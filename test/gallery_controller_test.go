package test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupGalleryController(t *testing.T) (*fiber.App, *usecasemock.GalleryUsecaseMock) {
	mockUC := new(usecasemock.GalleryUsecaseMock)
	controller := httpdelivery.NewGalleryController(mockUC, logrus.New())

	app := fiber.New()

	api := app.Group("/api")
	api.Get("/galleries", controller.GetAll)
	api.Get("/galleries/:id", controller.GetByID)
	api.Post("/galleries", controller.Create)
	api.Put("/galleries/:id", controller.Update)
	api.Delete("/galleries/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/galleries", controller.GetAllPublic)

	return app, mockUC
}

func TestGalleryController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupGalleryController(t)

	items := []model.GalleryResponse{{ID: "g1", Title: "Image 1"}}

	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/galleries", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_GetAllPublic_WithFilter(t *testing.T) {
	app, mockUC := setupGalleryController(t)
	items := []model.GalleryResponse{{ID: "g1", Title: "Pura Image"}}

	mockUC.On("GetPublic", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/galleries?entity_type=pura", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_GetAll_Success(t *testing.T) {
	app, mockUC := setupGalleryController(t)
	items := []model.GalleryResponse{{ID: "g1", Title: "Admin View"}}

	mockUC.On("GetAll", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/galleries", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_GetByID_Success(t *testing.T) {
	app, mockUC := setupGalleryController(t)
	targetID := "g1"
	item := &model.GalleryResponse{ID: targetID, Title: "Detail"}

	mockUC.On("GetByID", targetID).Return(item, nil)

	req := httptest.NewRequest("GET", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_Create_Success(t *testing.T) {
	app, mockUC := setupGalleryController(t)

	reqBody := model.CreateGalleryRequest{Title: "New", ImageURL: "http://img.com", EntityType: "pura"}
	resBody := &model.GalleryResponse{ID: "new-id", Title: "New"}

	mockUC.On("Create", mock.AnythingOfType("model.CreateGalleryRequest")).Return(resBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/galleries", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_Update_Success(t *testing.T) {
	app, mockUC := setupGalleryController(t)
	targetID := "g1"

	reqBody := model.UpdateGalleryRequest{Title: "Updated"}
	resBody := &model.GalleryResponse{ID: targetID, Title: "Updated"}

	mockUC.On("Update", targetID, mock.AnythingOfType("model.UpdateGalleryRequest")).Return(resBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/galleries/"+targetID, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestGalleryController_Delete_Success(t *testing.T) {
	app, mockUC := setupGalleryController(t)
	targetID := "g1"

	mockUC.On("Delete", targetID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/galleries/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}
