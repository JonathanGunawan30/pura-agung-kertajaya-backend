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

func setupActivityController() (*fiber.App, *usecasemock.ActivityUsecaseMock) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	controller := httpdelivery.NewActivityController(mockUC, logrus.New())
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

	api := app.Group("/api")
	api.Get("/activities", controller.GetAll)
	api.Get("/activities/:id", controller.GetByID)
	api.Post("/activities", controller.Create)
	api.Put("/activities/:id", controller.Update)
	api.Delete("/activities/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/activities", controller.GetAllPublic)

	return app, mockUC
}

func TestActivityController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	items := []model.ActivityResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	mockUC.On("GetPublic", "").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/public/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("GetPublic", "").Return(([]model.ActivityResponse)(nil), errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/public/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAll_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	items := []model.ActivityResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	mockUC.On("GetAll", "").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAll_Error(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("GetAll", "").Return(nil, errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetByID_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	item := &model.ActivityResponse{ID: "x", Title: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)
	req := httptest.NewRequest("GET", "/api/activities/x", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetByID_Invalid(t *testing.T) {
	app, _ := setupActivityController()
	req := httptest.NewRequest("GET", "/api/activities/", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestActivityController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("GetByID", "missing").Return((*model.ActivityResponse)(nil), gorm.ErrRecordNotFound)
	req := httptest.NewRequest("GET", "/api/activities/missing", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "The requested resource was not found.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Create_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.CreateActivityRequest{EntityType: "pura", Title: "T", Description: "D", OrderIndex: 1, IsActive: true}
	resBody := &model.ActivityResponse{ID: "1", Title: "T"}
	mockUC.On("Create", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "1", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Create_BadBody(t *testing.T) {
	app, _ := setupActivityController()
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.CreateActivityRequest{}
	validate := validator.New()
	err := validate.Struct(reqBody)
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)
	mockUC.On("Create", reqBody).Return((*model.ActivityResponse)(nil), validationErrs)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Validation failed", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Update_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.UpdateActivityRequest{Title: "New", Description: "D", IsActive: true, OrderIndex: 1}
	resBody := &model.ActivityResponse{ID: "2", Title: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/activities/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "2", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Update_BadBody(t *testing.T) {
	app, _ := setupActivityController()
	req := httptest.NewRequest("PUT", "/api/activities/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.UpdateActivityRequest{Title: "N", Description: "D", IsActive: true, OrderIndex: 1}
	mockUC.On("Update", "3", reqBody).Return((*model.ActivityResponse)(nil), errors.New("update failed"))
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/activities/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Delete_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("Delete", "7").Return(nil)
	req := httptest.NewRequest("DELETE", "/api/activities/7", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Activity deleted successfully", response.Data)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("Delete", "8").Return(errors.New("delete failed"))
	req := httptest.NewRequest("DELETE", "/api/activities/8", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}
