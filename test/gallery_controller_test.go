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
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupGalleryController(mockUC *usecasemock.GalleryUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewGalleryController(mockUC, logger)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})

	api := app.Group("/api")
	api.Get("/galleries", controller.GetAll)
	api.Get("/galleries/:id", controller.GetByID)
	api.Post("/galleries", controller.Create)
	api.Put("/galleries/:id", controller.Update)
	api.Delete("/galleries/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/galleries", controller.GetAllPublic)

	return app
}

func TestGalleryController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	items := []model.GalleryResponse{{ID: "g1", Title: "Image 1"}}
	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/galleries", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_GetAllPublic_WithFilter(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	items := []model.GalleryResponse{{ID: "g1", Title: "Pura Image"}}
	mockUC.On("GetPublic", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/galleries?entity_type=pura", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.GalleryResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/galleries", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestGalleryController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	items := []model.GalleryResponse{{ID: "g1", Title: "Admin View"}}
	// "pura" injected by middleware mock
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/galleries", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "g1"
	item := &model.GalleryResponse{ID: targetID, Title: "Detail"}
	mockUC.On("GetByID", targetID).Return(item, nil)

	req := httptest.NewRequest("GET", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGalleryController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "missing"
	mockUC.On("GetByID", targetID).Return((*model.GalleryResponse)(nil), model.ErrNotFound("gallery not found"))

	req := httptest.NewRequest("GET", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGalleryController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	reqBody := model.CreateGalleryRequest{
		Title:      "New",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
		EntityType: "pura",
	}
	resBody := &model.GalleryResponse{ID: "1", Title: "New"}

	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/galleries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestGalleryController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	reqBody := model.CreateGalleryRequest{}

	validate := validator.New()
	type Dummy struct {
		Title string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.GalleryResponse)(nil), realValErr)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/galleries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGalleryController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "g1"
	reqBody := model.UpdateGalleryRequest{Title: "Updated"}
	resBody := &model.GalleryResponse{ID: targetID, Title: "Updated"}

	mockUC.On("Update", targetID, reqBody).Return(resBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/galleries/"+targetID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGalleryController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "missing"
	reqBody := model.UpdateGalleryRequest{Title: "Updated"}
	mockUC.On("Update", targetID, reqBody).Return((*model.GalleryResponse)(nil), model.ErrNotFound("gallery not found"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/galleries/"+targetID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGalleryController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "g1"
	mockUC.On("Delete", targetID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGalleryController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "missing"
	mockUC.On("Delete", targetID).Return(model.ErrNotFound("gallery not found"))

	req := httptest.NewRequest("DELETE", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGalleryController_Delete_InternalError(t *testing.T) {
	mockUC := &usecasemock.GalleryUsecaseMock{}
	app := setupGalleryController(mockUC)

	targetID := "error"
	mockUC.On("Delete", targetID).Return(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
