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

func setupAboutController() (*fiber.App, *usecasemock.AboutUsecaseMock) {
	mockUC := &usecasemock.AboutUsecaseMock{}
	controller := httpdelivery.NewAboutController(mockUC, logrus.New())
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "An internal server error occurred. Please try again later."
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				code = fiber.StatusNotFound
				message = "The requested resource was not found." // Standard 404 message
			} else if _, ok := err.(validator.ValidationErrors); ok {
				code = fiber.StatusBadRequest
				message = "Validation failed" // Simplified message for test validation errors
			}
			return ctx.Status(code).JSON(model.WebResponse[any]{Errors: message})
		},
	})

	api := app.Group("/api")
	api.Get("/about", controller.GetAll)
	api.Get("/about/:id", controller.GetByID)
	api.Post("/about", controller.Create)
	api.Put("/about/:id", controller.Update)
	api.Delete("/about/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/about", controller.GetAllPublic)

	return app, mockUC
}

func TestAboutController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupAboutController()
	items := []model.AboutSectionResponse{{ID: "1", EntityType: "pura", Title: "A"}, {ID: "2", EntityType: "pura", Title: "B"}}
	mockUC.On("GetPublic", "pura").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/public/about?entity_type=pura", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupAboutController()
	mockUC.On("GetPublic", "").Return(([]model.AboutSectionResponse)(nil), errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/public/about", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestAboutController_GetAll_Success(t *testing.T) {
	app, mockUC := setupAboutController()
	items := []model.AboutSectionResponse{{ID: "1", EntityType: "pura", Title: "A"}, {ID: "2", EntityType: "pura", Title: "B"}}
	mockUC.On("GetAll", "").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/about", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_GetAll_Error(t *testing.T) {
	app, mockUC := setupAboutController()
	mockUC.On("GetAll", "").Return(nil, errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/about", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestAboutController_GetByID_Success(t *testing.T) {
	app, mockUC := setupAboutController()
	item := &model.AboutSectionResponse{ID: "x", Title: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)
	req := httptest.NewRequest("GET", "/api/about/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupAboutController()
	mockUC.On("GetByID", "missing").Return((*model.AboutSectionResponse)(nil), gorm.ErrRecordNotFound)
	req := httptest.NewRequest("GET", "/api/about/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "The requested resource was not found.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestAboutController_Create_Success(t *testing.T) {
	app, mockUC := setupAboutController()
	reqBody := model.AboutSectionRequest{EntityType: "pura", Title: "T", Description: "D", IsActive: true}
	resBody := &model.AboutSectionResponse{ID: "1", EntityType: "pura", Title: "T"}
	mockUC.On("Create", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/about", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_Create_BadBody(t *testing.T) {
	app, _ := setupAboutController()
	req := httptest.NewRequest("POST", "/api/about", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAboutController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupAboutController()
	reqBody := model.AboutSectionRequest{} // Invalid request
	validate := validator.New()
	err := validate.Struct(reqBody)
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)
	mockUC.On("Create", reqBody).Return((*model.AboutSectionResponse)(nil), validationErrs)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/about", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Validation failed", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestAboutController_Update_Success(t *testing.T) {
	app, mockUC := setupAboutController()
	reqBody := model.AboutSectionRequest{Title: "New", Description: "D", IsActive: true}
	resBody := &model.AboutSectionResponse{ID: "2", Title: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/about/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_Update_BadBody(t *testing.T) {
	app, _ := setupAboutController()
	req := httptest.NewRequest("PUT", "/api/about/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAboutController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupAboutController()
	reqBody := model.AboutSectionRequest{Title: "N", Description: "D", IsActive: true}
	mockUC.On("Update", "3", reqBody).Return((*model.AboutSectionResponse)(nil), errors.New("update failed"))
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/about/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestAboutController_Delete_Success(t *testing.T) {
	app, mockUC := setupAboutController()
	mockUC.On("Delete", "7").Return(nil)
	req := httptest.NewRequest("DELETE", "/api/about/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestAboutController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupAboutController()
	mockUC.On("Delete", "8").Return(errors.New("delete failed"))
	req := httptest.NewRequest("DELETE", "/api/about/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}
