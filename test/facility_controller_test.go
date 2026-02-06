package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupFacilityController(mockUC *usecasemock.FacilityUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewFacilityController(mockUC, logger)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})

	api := app.Group("/api")
	api.Get("/facilities", controller.GetAll)
	api.Get("/facilities/:id", controller.GetByID)
	api.Post("/facilities", controller.Create)
	api.Put("/facilities/:id", controller.Update)
	api.Delete("/facilities/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/facilities", controller.GetAllPublic)

	return app
}

func TestFacilityController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	items := []model.FacilityResponse{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/facilities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.FacilityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestFacilityController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.FacilityResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/facilities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}

func TestFacilityController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	items := []model.FacilityResponse{{ID: "1", Name: "A"}}
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/facilities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestFacilityController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	item := &model.FacilityResponse{ID: "x", Name: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/api/facilities/x", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.FacilityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
}

func TestFacilityController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	mockUC.On("GetByID", "missing").Return((*model.FacilityResponse)(nil), model.ErrNotFound("facility not found"))

	req := httptest.NewRequest("GET", "/api/facilities/missing", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "facility not found", response.Errors)
}

func TestFacilityController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	reqBody := model.CreateFacilityRequest{EntityType: "pura", Name: "New Facility"}
	resBody := &model.FacilityResponse{ID: "1", Name: "New Facility"}

	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/facilities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestFacilityController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	reqBody := model.CreateFacilityRequest{}

	validate := validator.New()
	type Dummy struct {
		Name string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.FacilityResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/facilities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFacilityController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	reqBody := model.UpdateFacilityRequest{Name: "Updated"}
	resBody := &model.FacilityResponse{ID: "2", Name: "Updated"}

	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/facilities/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestFacilityController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	reqBody := model.UpdateFacilityRequest{Name: "Updated"}
	mockUC.On("Update", "3", reqBody).Return((*model.FacilityResponse)(nil), model.ErrNotFound("facility not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/facilities/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFacilityController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/facilities/7", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestFacilityController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	mockUC.On("Delete", "8").Return(model.ErrNotFound("facility not found"))

	req := httptest.NewRequest("DELETE", "/api/facilities/8", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFacilityController_Delete_InternalError(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	app := setupFacilityController(mockUC)

	mockUC.On("Delete", "9").Return(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/facilities/9", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
