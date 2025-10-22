package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupContactInfoController() (*fiber.App, *usecasemock.ContactInfoUsecaseMock) {
	log = logrus.New()
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	controller := httpdelivery.NewContactInfoController(mockUC, logrus.New())
	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	app.Use(middleware.ErrorHandlerMiddleware(log))

	app.Get("/contact-info", controller.GetAll)
	app.Get("/contact-info/:id", controller.GetByID)
	app.Post("/contact-info", controller.Create)
	app.Put("/contact-info/:id", controller.Update)
	app.Delete("/contact-info/:id", controller.Delete)

	return app, mockUC
}

func TestContactInfoController_GetAll_Success(t *testing.T) {
	app, mockUC := setupContactInfoController()
	items := []model.ContactInfoResponse{{ID: "1", Address: "A"}, {ID: "2", Address: "B"}}
	mockUC.On("GetAll").Return(items, nil)

	req := httptest.NewRequest("GET", "/contact-info", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.ContactInfoResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestContactInfoController_GetAll_Error(t *testing.T) {
	app, mockUC := setupContactInfoController()
	mockUC.On("GetAll").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/contact-info", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestContactInfoController_GetByID_Success(t *testing.T) {
	app, mockUC := setupContactInfoController()
	item := &model.ContactInfoResponse{ID: "x", Address: "Addr"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/contact-info/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.ContactInfoResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestContactInfoController_GetByID_Invalid(t *testing.T) {
	app, _ := setupContactInfoController()
	req := httptest.NewRequest("GET", "/contact-info/", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode) // no route matches
}

func TestContactInfoController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupContactInfoController()
	mockUC.On("GetByID", "missing").Return((*model.ContactInfoResponse)(nil), errors.New("not found"))

	req := httptest.NewRequest("GET", "/contact-info/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestContactInfoController_Create_Success(t *testing.T) {
	app, mockUC := setupContactInfoController()
	reqBody := model.ContactInfoRequest{Address: "A", Email: "e@x.com"}
	resBody := &model.ContactInfoResponse{ID: "1", Address: "A"}
	mockUC.On("Create", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/contact-info", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response model.WebResponse[*model.ContactInfoResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "1", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestContactInfoController_Create_BadBody(t *testing.T) {
	app, _ := setupContactInfoController()
	req := httptest.NewRequest("POST", "/contact-info", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestContactInfoController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupContactInfoController()
	reqBody := model.ContactInfoRequest{}
	mockUC.On("Create", reqBody).Return((*model.ContactInfoResponse)(nil), errors.New("validation failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/contact-info", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestContactInfoController_Update_Success(t *testing.T) {
	app, mockUC := setupContactInfoController()
	reqBody := model.ContactInfoRequest{Address: "New"}
	resBody := &model.ContactInfoResponse{ID: "2", Address: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/contact-info/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.ContactInfoResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "2", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestContactInfoController_Update_BadBody(t *testing.T) {
	app, _ := setupContactInfoController()
	req := httptest.NewRequest("PUT", "/contact-info/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestContactInfoController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupContactInfoController()
	reqBody := model.ContactInfoRequest{Address: "A"}
	mockUC.On("Update", "3", reqBody).Return((*model.ContactInfoResponse)(nil), errors.New("update failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/contact-info/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestContactInfoController_Delete_Success(t *testing.T) {
	app, mockUC := setupContactInfoController()
	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/contact-info/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Contact info deleted successfully", response.Data)
	mockUC.AssertExpectations(t)
}

func TestContactInfoController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupContactInfoController()
	mockUC.On("Delete", "8").Return(errors.New("delete failed"))

	req := httptest.NewRequest("DELETE", "/contact-info/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
