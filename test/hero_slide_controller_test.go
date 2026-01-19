package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
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
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	controller := httpdelivery.NewHeroSlideController(mockUC, logrus.New())

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "An internal server error occurred. Please try again later."
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				code = fiber.StatusNotFound
				message = "The requested resource was not found."
			} else if _, ok := err.(validator.ValidationErrors); ok {
				code = fiber.StatusBadRequest
				message = "Validation failed"
			}
			return ctx.Status(code).JSON(model.WebResponse[any]{Errors: message})
		},
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("entity_type", "pura")
		return c.Next()
	})
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
		{ID: "a", EntityType: "pura", Images: model.ImageVariants{Lg: "https://a"}},
		{ID: "b", EntityType: "pura", Images: model.ImageVariants{Lg: "https://b"}},
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
		{ID: "a", EntityType: "pura", Images: model.ImageVariants{Lg: "https://a"}},
		{ID: "b", EntityType: "pura", Images: model.ImageVariants{Lg: "https://b"}},
	}
	mockUC.On("GetAll", "pura").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/hero-slides", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_GetAll_Error(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	mockUC.On("GetAll", "pura").Return(nil, errors.New("db error"))

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
	item := &model.HeroSlideResponse{ID: "x", Images: model.ImageVariants{Lg: "https://x"}}
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
	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://img.com/lg.jpg"}, OrderIndex: 1, IsActive: true}
	resBody := &model.HeroSlideResponse{ID: "1", EntityType: "pura"}
	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)
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
		// entity_type is set to "pura" by test middleware
		validate := validator.New()
		err := validate.Struct(reqBody)
		var validationErrs validator.ValidationErrors
		errors.As(err, &validationErrs)
		mockUC.On("Create", "pura", reqBody).Return((*model.HeroSlideResponse)(nil), validationErrs)
	

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hero-slides", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Validation failed", response.Errors)

	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Update_Success(t *testing.T) {
	app, mockUC := setupHeroSlideController()
	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://new"}, OrderIndex: 3, IsActive: false}
	resBody := &model.HeroSlideResponse{ID: "2", EntityType: "pura", Images: model.ImageVariants{Lg: "https://new"}}

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
