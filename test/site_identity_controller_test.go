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

func setupSiteIdentityController(mockUC *usecasemock.SiteIdentityUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewSiteIdentityController(mockUC, logger)

	publicApi := app.Group("/api/public")
	publicApi.Get("/site-identity", controller.GetPublic)

	api := app.Group("/api", func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})
	api.Get("/site-identity", controller.GetAll)
	api.Get("/site-identity/:id", controller.GetByID)
	api.Post("/site-identity", controller.Create)
	api.Put("/site-identity/:id", controller.Update)
	api.Delete("/site-identity/:id", controller.Delete)

	return app
}

func TestSiteIdentityController_GetPublic_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	item := &model.SiteIdentityResponse{ID: "x", EntityType: "pura", SiteName: "Pura"}
	mockUC.On("GetPublic", "pura").Return(item, nil)

	req := httptest.NewRequest("GET", "/api/public/site-identity?entity_type=pura", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.SiteIdentityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestSiteIdentityController_GetPublic_NotFound(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	mockUC.On("GetPublic", "").Return((*model.SiteIdentityResponse)(nil), model.ErrNotFound("site identity not found"))

	req := httptest.NewRequest("GET", "/api/public/site-identity", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSiteIdentityController_GetPublic_Error(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	mockUC.On("GetPublic", "").Return((*model.SiteIdentityResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/site-identity", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSiteIdentityController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	items := []model.SiteIdentityResponse{{ID: "1", EntityType: "pura", SiteName: "A"}}
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/site-identity", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestSiteIdentityController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	item := &model.SiteIdentityResponse{ID: "x", SiteName: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/api/site-identity/x", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSiteIdentityController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	mockUC.On("GetByID", "missing").Return((*model.SiteIdentityResponse)(nil), model.ErrNotFound("site identity not found"))

	req := httptest.NewRequest("GET", "/api/site-identity/missing", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSiteIdentityController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	reqBody := model.SiteIdentityRequest{SiteName: "Pura", LogoURL: "http://logo.com/logo.png"}
	resBody := &model.SiteIdentityResponse{ID: "1", SiteName: "Pura"}

	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/site-identity", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestSiteIdentityController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	reqBody := model.SiteIdentityRequest{}

	validate := validator.New()
	type Dummy struct {
		SiteName string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.SiteIdentityResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/site-identity", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSiteIdentityController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	reqBody := model.SiteIdentityRequest{SiteName: "New"}
	resBody := &model.SiteIdentityResponse{ID: "2", SiteName: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/site-identity/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSiteIdentityController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	reqBody := model.SiteIdentityRequest{SiteName: "X"}
	mockUC.On("Update", "3", reqBody).Return((*model.SiteIdentityResponse)(nil), model.ErrNotFound("site identity not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/site-identity/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSiteIdentityController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/site-identity/7", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSiteIdentityController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	app := setupSiteIdentityController(mockUC)

	mockUC.On("Delete", "8").Return(model.ErrNotFound("site identity not found"))

	req := httptest.NewRequest("DELETE", "/api/site-identity/8", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
