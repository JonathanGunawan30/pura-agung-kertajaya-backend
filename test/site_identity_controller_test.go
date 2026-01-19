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

func setupSiteIdentityController() (*fiber.App, *usecasemock.SiteIdentityUsecaseMock) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	controller := httpdelivery.NewSiteIdentityController(mockUC, logrus.New())
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

	app.Get("/site-identity", controller.GetAll)
	app.Get("/site-identity/:id", controller.GetByID)
	app.Post("/site-identity", controller.Create)
	app.Put("/site-identity/:id", controller.Update)
	app.Delete("/site-identity/:id", controller.Delete)

	return app, mockUC
}

func TestSiteIdentityController_GetPublic_Success(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	controller := httpdelivery.NewSiteIdentityController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/site-identity", controller.GetPublic)

	item := &model.SiteIdentityResponse{ID: "x", EntityType: "pura", SiteName: "Pura"}
	mockUC.On("GetPublic", "pura").Return(item, nil)

	req := httptest.NewRequest("GET", "/public/site-identity?entity_type=pura", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.SiteIdentityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestSiteIdentityController_GetPublic_Error(t *testing.T) {
	mockUC := &usecasemock.SiteIdentityUsecaseMock{}
	controller := httpdelivery.NewSiteIdentityController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/site-identity", controller.GetPublic)

	mockUC.On("GetPublic", "").Return((*model.SiteIdentityResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/public/site-identity", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSiteIdentityController_GetAll_Success(t *testing.T) {
	app, mockUC := setupSiteIdentityController()

	items := []model.SiteIdentityResponse{{ID: "1", EntityType: "pura", SiteName: "A"}, {ID: "2", EntityType: "pura", SiteName: "B"}}
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/site-identity?entity_type=pura", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.SiteIdentityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestSiteIdentityController_GetAll_Error(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	mockUC.On("GetAll", "pura").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/site-identity", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSiteIdentityController_GetByID_Success(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	item := &model.SiteIdentityResponse{ID: "x", SiteName: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/site-identity/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.SiteIdentityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
}

func TestSiteIdentityController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	mockUC.On("GetByID", "missing").Return((*model.SiteIdentityResponse)(nil), gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/site-identity/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSiteIdentityController_Create_Success(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	reqBody := model.SiteIdentityRequest{EntityType: "pura", SiteName: "Pura", LogoURL: "http://logo.com/logo.png", Tagline: "T"}
	resBody := &model.SiteIdentityResponse{ID: "1", SiteName: "Pura"}
	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/site-identity", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response model.WebResponse[*model.SiteIdentityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "1", response.Data.ID)
}

func TestSiteIdentityController_Create_BadBody(t *testing.T) {
	app, _ := setupSiteIdentityController()
	req := httptest.NewRequest("POST", "/site-identity", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSiteIdentityController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	reqBody := model.SiteIdentityRequest{}
	// entity_type is set to "pura" by test middleware
	validate := validator.New()
	err := validate.Struct(reqBody)
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)
	mockUC.On("Create", "pura", reqBody).Return((*model.SiteIdentityResponse)(nil), validationErrs)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/site-identity", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSiteIdentityController_Update_Success(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	reqBody := model.SiteIdentityRequest{SiteName: "New"}
	resBody := &model.SiteIdentityResponse{ID: "2", SiteName: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/site-identity/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.SiteIdentityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "2", response.Data.ID)
}

func TestSiteIdentityController_Update_BadBody(t *testing.T) {
	app, _ := setupSiteIdentityController()
	req := httptest.NewRequest("PUT", "/site-identity/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSiteIdentityController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	reqBody := model.SiteIdentityRequest{SiteName: "X"}
	mockUC.On("Update", "3", reqBody).Return((*model.SiteIdentityResponse)(nil), errors.New("update failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/site-identity/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSiteIdentityController_Delete_Success(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/site-identity/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Site identity deleted successfully", response.Data)
}

func TestSiteIdentityController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupSiteIdentityController()
	mockUC.On("Delete", "8").Return(errors.New("delete failed"))

	req := httptest.NewRequest("DELETE", "/site-identity/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
