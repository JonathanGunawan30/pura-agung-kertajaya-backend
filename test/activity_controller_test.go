package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupActivityController() (*fiber.App, *usecasemock.ActivityUsecaseMock) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	controller := httpdelivery.NewActivityController(mockUC, logrus.New())
	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	app.Get("/activities", controller.GetAll)
	app.Get("/activities/:id", controller.GetByID)
	app.Post("/activities", controller.Create)
	app.Put("/activities/:id", controller.Update)
	app.Delete("/activities/:id", controller.Delete)

	return app, mockUC
}

func TestActivityController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	controller := httpdelivery.NewActivityController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/activities", controller.GetAllPublic)

	items := []model.ActivityResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	mockUC.On("GetPublic").Return(items, nil)

	req := httptest.NewRequest("GET", "/public/activities", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.ActivityUsecaseMock{}
	controller := httpdelivery.NewActivityController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/activities", controller.GetAllPublic)

	mockUC.On("GetPublic").Return(([]model.ActivityResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/public/activities", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestActivityController_GetAll_Success(t *testing.T) {
	app, mockUC := setupActivityController()

	items := []model.ActivityResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	mockUC.On("GetAll").Return(items, nil)

	req := httptest.NewRequest("GET", "/activities", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetAll_Error(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("GetAll").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/activities", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestActivityController_GetByID_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	item := &model.ActivityResponse{ID: "x", Title: "X"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/activities/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestActivityController_GetByID_Invalid(t *testing.T) {
	app, _ := setupActivityController()
	req := httptest.NewRequest("GET", "/activities/", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode) // route not matched -> 404
}

func TestActivityController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("GetByID", "missing").Return((*model.ActivityResponse)(nil), errors.New("not found"))

	req := httptest.NewRequest("GET", "/activities/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestActivityController_Create_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.ActivityRequest{Title: "T", Description: "D", OrderIndex: 1}
	resBody := &model.ActivityResponse{ID: "1", Title: "T"}
	mockUC.On("Create", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "1", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Create_BadBody(t *testing.T) {
	app, _ := setupActivityController()
	req := httptest.NewRequest("POST", "/activities", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.ActivityRequest{}
	mockUC.On("Create", reqBody).Return((*model.ActivityResponse)(nil), errors.New("validation failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/activities", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Update_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.ActivityRequest{Title: "New", Description: "D"}
	resBody := &model.ActivityResponse{ID: "2", Title: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/activities/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.ActivityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "2", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Update_BadBody(t *testing.T) {
	app, _ := setupActivityController()
	req := httptest.NewRequest("PUT", "/activities/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestActivityController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupActivityController()
	reqBody := model.ActivityRequest{Title: "N"}
	mockUC.On("Update", "3", reqBody).Return((*model.ActivityResponse)(nil), errors.New("update failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/activities/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestActivityController_Delete_Success(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/activities/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Activity deleted successfully", response.Data)
	mockUC.AssertExpectations(t)
}

func TestActivityController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupActivityController()
	mockUC.On("Delete", "8").Return(errors.New("delete failed"))

	req := httptest.NewRequest("DELETE", "/activities/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
