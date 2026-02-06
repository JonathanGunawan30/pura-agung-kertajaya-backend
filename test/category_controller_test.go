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

func setupCategoryController(mockUC *usecasemock.CategoryUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewCategoryController(mockUC, logger)

	app.Get("/public/categories", controller.GetAllPublic)
	app.Get("/categories", controller.GetAll)
	app.Get("/categories/:id", controller.GetByID)
	app.Post("/categories", controller.Create)
	app.Put("/categories/:id", controller.Update)
	app.Delete("/categories/:id", controller.Delete)

	return app
}

func TestCategoryController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	items := []model.CategoryResponse{
		{ID: "1", Name: "Adat", Slug: "adat"},
		{ID: "2", Name: "Upacara", Slug: "upacara"},
	}

	mockUC.On("GetAll").Return(items, nil)

	req := httptest.NewRequest("GET", "/public/categories", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.CategoryResponse]
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Len(t, response.Data, 2)
	assert.Equal(t, "Adat", response.Data[0].Name)
	mockUC.AssertExpectations(t)
}

func TestCategoryController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	mockUC.On("GetAll").Return(([]model.CategoryResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/public/categories", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}

func TestCategoryController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	item := &model.CategoryResponse{ID: "1", Name: "Adat", Slug: "adat"}
	mockUC.On("GetByID", "1").Return(item, nil)

	req := httptest.NewRequest("GET", "/categories/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestCategoryController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	mockUC.On("GetByID", "99").Return((*model.CategoryResponse)(nil), model.ErrNotFound("category not found"))

	req := httptest.NewRequest("GET", "/categories/99", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "category not found", response.Errors)
}

func TestCategoryController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	payload := model.CreateCategoryRequest{Name: "Baru"}
	mockResp := &model.CategoryResponse{ID: "100", Name: "Baru", Slug: "baru"}

	mockUC.On("Create", payload).Return(mockResp, nil)

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	var response model.WebResponse[model.CategoryResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "baru", response.Data.Slug)
}

func TestCategoryController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	payload := model.CreateCategoryRequest{}
	validate := validator.New()
	type Dummy struct {
		Name string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", payload).Return((*model.CategoryResponse)(nil), realValErr)

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCategoryController_Create_BodyParserError(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCategoryController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	payload := model.UpdateCategoryRequest{Name: "Updated"}
	mockResp := &model.CategoryResponse{ID: "1", Name: "Updated", Slug: "updated"}

	mockUC.On("Update", "1", payload).Return(mockResp, nil)

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/categories/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestCategoryController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	payload := model.UpdateCategoryRequest{Name: "Updated"}
	mockUC.On("Update", "99", payload).Return((*model.CategoryResponse)(nil), model.ErrNotFound("category not found"))

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/categories/99", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestCategoryController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	targetID := "uuid-123"
	mockUC.On("Delete", targetID).Return(nil)

	req := httptest.NewRequest("DELETE", "/categories/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestCategoryController_Delete_Conflict(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	targetID := "uuid-used"
	mockUC.On("Delete", targetID).Return(model.ErrConflict("category is currently in use"))

	req := httptest.NewRequest("DELETE", "/categories/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestCategoryController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	app := setupCategoryController(mockUC)

	targetID := "uuid-missing"
	mockUC.On("Delete", targetID).Return(model.ErrNotFound("category not found"))

	req := httptest.NewRequest("DELETE", "/categories/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
