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

func setupActivityController(mockUC *usecasemock.ActivityUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewActivityController(mockUC, logger)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})

	api := app.Group("/api")
	api.Get("/activities", controller.GetAll)
	api.Get("/activities/:id", controller.GetByID)
	api.Post("/activities", controller.Create)
	api.Put("/activities/:id", controller.Update)
	api.Delete("/activities/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/activities", controller.GetAllPublic)

	return app
}

func TestActivityController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	items := []model.ActivityResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.ActivityResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}

func TestActivityController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	items := []model.ActivityResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAll_Error(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	mockUC.On("GetAll", "pura").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/activities", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}

func TestActivityController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	item := &model.ActivityResponse{ID: "x", Title: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/api/activities/x", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
}

func TestActivityController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	mockUC.On("GetByID", "missing").Return((*model.ActivityResponse)(nil), model.ErrNotFound("activity not found"))

	req := httptest.NewRequest("GET", "/api/activities/missing", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "activity not found", response.Errors)
}

func TestActivityController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	reqBody := model.CreateActivityRequest{EntityType: "pura", Title: "T", Description: "D", OrderIndex: 1, IsActive: true}
	resBody := &model.ActivityResponse{ID: "1", Title: "T"}

	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestActivityController_Create_BadBody(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Create_BadRequest_Date(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	reqBody := model.CreateActivityRequest{Title: "Bad Date"}
	mockUC.On("Create", "pura", reqBody).Return((*model.ActivityResponse)(nil), model.ErrBadRequest("invalid date"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "invalid date", response.Errors)
}

func TestActivityController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	reqBody := model.CreateActivityRequest{}
	validate := validator.New()
	type Dummy struct {
		Title string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.ActivityResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	reqBody := model.UpdateActivityRequest{Title: "New", Description: "D", IsActive: true, OrderIndex: 1}
	resBody := &model.ActivityResponse{ID: "2", Title: "New"}

	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/activities/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestActivityController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	reqBody := model.UpdateActivityRequest{Title: "New"}
	mockUC.On("Update", "3", reqBody).Return((*model.ActivityResponse)(nil), model.ErrNotFound("activity not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/activities/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestActivityController_Update_BadRequest_Date(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	reqBody := model.UpdateActivityRequest{Title: "Bad Date"}
	mockUC.On("Update", "4", reqBody).Return((*model.ActivityResponse)(nil), model.ErrBadRequest("invalid date"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/activities/4", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/activities/7", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestActivityController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	mockUC.On("Delete", "8").Return(model.ErrNotFound("activity not found"))

	req := httptest.NewRequest("DELETE", "/api/activities/8", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestActivityController_Delete_InternalError(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	app := setupActivityController(mockUC)

	mockUC.On("Delete", "9").Return(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/activities/9", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
