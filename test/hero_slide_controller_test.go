package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupHeroSlideController() (*fiber.App, *usecasemock.HeroSlideUsecaseMock) {
	log := logrus.New()
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	controller := httpdelivery.NewHeroSlideController(mockUC, logrus.New())

	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	app.Use(middleware.ErrorHandlerMiddleware(log))

	api := app.Group("/api")
	api.Get("/hero-slides", controller.GetAll)
	api.Get("/hero-slides/:id", controller.GetByID)
	api.Post("/hero-slides", controller.Create)
	api.Put("/hero-slides/:id", controller.Update)
	api.Delete("/hero-slides/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/hero-slides", controller.GetAllPublic)

	return app, mockUC
}

func TestHeroSlideController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	items := []model.HeroSlideResponse{
		{ID: "a", EntityType: "pura", Images: map[string]string{"lg": "https://a"}},
		{ID: "b", EntityType: "pura", Images: map[string]string{"lg": "https://b"}},
	}
	mockUC.On("GetPublic", "pura").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/public/hero-slides?entity_type=pura", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	mockUC.On("GetPublic", "").Return(([]model.HeroSlideResponse)(nil), errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/public/hero-slides", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestHeroSlideController_GetAll_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	items := []model.HeroSlideResponse{
		{ID: "a", EntityType: "pura", Images: map[string]string{"lg": "https://a"}},
		{ID: "b", EntityType: "pura", Images: map[string]string{"lg": "https://b"}},
	}
	mockUC.On("GetAll", "").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/hero-slides", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_GetAll_Error(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	mockUC.On("GetAll", "").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/hero-slides", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetByID_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	item := &model.HeroSlideResponse{ID: "x", Images: map[string]string{"lg": "https://x"}}
	mockUC.On("GetByID", "x").Return(item, nil)
	req := httptest.NewRequest("GET", "/api/hero-slides/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	mockUC.On("GetByID", "missing").Return((*model.HeroSlideResponse)(nil), gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/api/hero-slides/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "The requested resource was not found.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Create_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://img"}, OrderIndex: 1, IsActive: true}
	resBody := &model.HeroSlideResponse{ID: "1", EntityType: "pura", Images: map[string]string{"lg": "https://img"}}

	mockUC.On("Create", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hero-slides", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestHeroSlideController_Create_BadBody(t *testing.T) {
	app, _ := setupHeroSlideController()
	req := httptest.NewRequest("POST", "/api/hero-slides", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHeroSlideController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	reqBody := model.HeroSlideRequest{}

	validate := validator.New()
	err := validate.Struct(reqBody)
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)

	mockUC.On("Create", reqBody).Return((*model.HeroSlideResponse)(nil), validationErrs)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hero-slides", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Contains(t, responseBody, "Images")
	assert.Contains(t, responseBody, "required")

	var response model.WebResponse[any]
	json.Unmarshal(bodyBytes, &response)
	assert.Contains(t, response.Errors, "Images")
	assert.Contains(t, response.Errors, "required")

	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Update_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://new"}, OrderIndex: 3, IsActive: false}
	resBody := &model.HeroSlideResponse{ID: "2", EntityType: "pura", Images: map[string]string{"lg": "https://new"}}

	mockUC.On("Update", "2", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/hero-slides/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_Update_BadBody(t *testing.T) {
	app, _ := setupHeroSlideController()
	req := httptest.NewRequest("PUT", "/api/hero-slides/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHeroSlideController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://img"}}
	mockUC.On("Update", "3", reqBody).Return((*model.HeroSlideResponse)(nil), errors.New("update failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/hero-slides/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Delete_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	mockUC.On("Delete", "7").Return(nil)
	req := httptest.NewRequest("DELETE", "/api/hero-slides/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	mockUC.On("Delete", "8").Return(errors.New("delete failed"))

	req := httptest.NewRequest("DELETE", "/api/hero-slides/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}
