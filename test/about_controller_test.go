package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupAboutController(mockUC *usecasemock.AboutUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()

	controller := httpdelivery.NewAboutController(mockUC, logger)

	api := app.Group("/api")
	api.Get("/about", controller.GetAll)
	api.Get("/about/:id", controller.GetByID)
	api.Post("/about", controller.Create)
	api.Put("/about/:id", controller.Update)
	api.Delete("/about/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/about", controller.GetAllPublic)

	return app
}

func TestAboutController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	items := []model.AboutSectionResponse{{ID: "1", EntityType: "pura", Title: "A"}, {ID: "2", EntityType: "pura", Title: "B"}}
	mockUC.On("GetPublic", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/about?entity_type=pura", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.AboutSectionResponse)(nil), errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/public/about", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}

func TestAboutController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	items := []model.AboutSectionResponse{{ID: "1", EntityType: "pura", Title: "A"}, {ID: "2", EntityType: "pura", Title: "B"}}
	mockUC.On("GetAll", "").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/about", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAboutController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	item := &model.AboutSectionResponse{ID: "x", Title: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)
	req := httptest.NewRequest("GET", "/api/about/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAboutController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	expectedErr := model.ErrNotFound("about section not found")
	mockUC.On("GetByID", "missing").Return((*model.AboutSectionResponse)(nil), expectedErr)

	req := httptest.NewRequest("GET", "/api/about/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "about section not found", response.Errors)
}

func TestAboutController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	reqBody := model.AboutSectionRequest{EntityType: "pura", Title: "T", Description: "D", IsActive: true}
	resBody := &model.AboutSectionResponse{ID: "1", EntityType: "pura", Title: "T"}
	mockUC.On("Create", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/about", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestAboutController_Create_Validation_Error(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	reqBody := model.AboutSectionRequest{}

	validate := validator.New()
	type Dummy struct {
		Title string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", reqBody).Return((*model.AboutSectionResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/about", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAboutController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	reqBody := model.AboutSectionRequest{Title: "New", Description: "D", IsActive: true}
	resBody := &model.AboutSectionResponse{ID: "2", Title: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/about/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAboutController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	reqBody := model.AboutSectionRequest{Title: "N", Description: "D", IsActive: true}
	mockUC.On("Update", "3", reqBody).Return((*model.AboutSectionResponse)(nil), model.ErrNotFound("about section not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/about/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestAboutController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	mockUC.On("Delete", "7").Return(nil)
	req := httptest.NewRequest("DELETE", "/api/about/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAboutController_Delete_UsecaseError(t *testing.T) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	app := setupAboutController(mockUC)

	mockUC.On("Delete", "8").Return(errors.New("delete failed"))
	req := httptest.NewRequest("DELETE", "/api/about/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}
