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

func setupHeroSlideController(mockUC *usecasemock.HeroSlideUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewHeroSlideController(mockUC, logger)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
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

	return app
}

func TestHeroSlideController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	items := []model.HeroSlideResponse{
		{ID: "a", EntityType: "pura", Images: model.ImageVariants{Lg: "https://a"}},
		{ID: "b", EntityType: "pura", Images: model.ImageVariants{Lg: "https://b"}},
	}

	mockUC.On("GetPublic", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/hero-slides?entity_type=pura", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.HeroSlideResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/hero-slides", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestHeroSlideController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	items := []model.HeroSlideResponse{{ID: "a"}}
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/hero-slides", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	item := &model.HeroSlideResponse{ID: "x", Images: model.ImageVariants{Lg: "https://x"}}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/api/hero-slides/x", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	mockUC.On("GetByID", "missing").Return((*model.HeroSlideResponse)(nil), model.ErrNotFound("hero slide not found"))

	req := httptest.NewRequest("GET", "/api/hero-slides/missing", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "hero slide not found", response.Errors)
}

func TestHeroSlideController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://img.com/lg.jpg"}, OrderIndex: 1, IsActive: true}
	resBody := &model.HeroSlideResponse{ID: "1", EntityType: "pura"}

	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hero-slides", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestHeroSlideController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	reqBody := model.HeroSlideRequest{}

	validate := validator.New()
	type Dummy struct {
		Images map[string]string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.HeroSlideResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hero-slides", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHeroSlideController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	reqBody := model.HeroSlideRequest{EntityType: "pura", Images: map[string]string{"lg": "https://new"}, OrderIndex: 3, IsActive: false}
	resBody := &model.HeroSlideResponse{ID: "2", EntityType: "pura", Images: model.ImageVariants{Lg: "https://new"}}

	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/hero-slides/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	reqBody := model.HeroSlideRequest{Images: map[string]string{"lg": "https://img"}}
	mockUC.On("Update", "3", reqBody).Return((*model.HeroSlideResponse)(nil), model.ErrNotFound("hero slide not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/hero-slides/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestHeroSlideController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/hero-slides/7", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHeroSlideController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	mockUC.On("Delete", "8").Return(model.ErrNotFound("hero slide not found"))

	req := httptest.NewRequest("DELETE", "/api/hero-slides/8", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestHeroSlideController_Delete_InternalError(t *testing.T) {
	mockUC := &usecasemock.HeroSlideUsecaseMock{}
	app := setupHeroSlideController(mockUC)

	mockUC.On("Delete", "9").Return(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/hero-slides/9", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
