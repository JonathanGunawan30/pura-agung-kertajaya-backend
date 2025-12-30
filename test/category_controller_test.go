package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func TestCategoryController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	controller := httpdelivery.NewCategoryController(mockUC, logrus.New())
	app := fiber.New()

	app.Get("/public/categories", controller.GetAllPublic)

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
	controller := httpdelivery.NewCategoryController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/categories", controller.GetAllPublic)

	mockUC.On("GetAll").Return(([]model.CategoryResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/public/categories", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestCategoryController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	controller := httpdelivery.NewCategoryController(mockUC, logrus.New())
	app := fiber.New()
	app.Post("/categories", controller.Create)

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

func TestCategoryController_Create_BadRequest(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	controller := httpdelivery.NewCategoryController(mockUC, logrus.New())
	app := fiber.New()
	app.Post("/categories", controller.Create)

	req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCategoryController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.CategoryUsecaseMock{}
	controller := httpdelivery.NewCategoryController(mockUC, logrus.New())
	app := fiber.New()
	app.Delete("/categories/:id", controller.Delete)

	targetID := "uuid-123"
	mockUC.On("Delete", targetID).Return(nil)

	req := httptest.NewRequest("DELETE", "/categories/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}
