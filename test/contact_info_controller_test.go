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
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupContactInfoController(mockUC *usecasemock.ContactInfoUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewContactInfoController(mockUC, logger)

	app.Get("/contact-info", controller.GetAll)
	app.Get("/contact-info/:id", controller.GetByID)
	app.Post("/contact-info", controller.Create)
	app.Put("/contact-info/:id", controller.Update)
	app.Delete("/contact-info/:id", controller.Delete)

	return app
}

func TestContactInfoController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	items := []model.ContactInfoResponse{{ID: "1", Address: "A"}, {ID: "2", Address: "B"}}
	mockUC.On("GetAll", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/contact-info", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.ContactInfoResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestContactInfoController_GetAll_Error(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	mockUC.On("GetAll", "").Return(([]model.ContactInfoResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/contact-info", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}

func TestContactInfoController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	item := &model.ContactInfoResponse{ID: "x", Address: "Addr"}
	mockUC.On("GetByID", "x").Return(item, nil)

	req := httptest.NewRequest("GET", "/contact-info/x", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.ContactInfoResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "x", response.Data.ID)
}

func TestContactInfoController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	mockUC.On("GetByID", "missing").Return((*model.ContactInfoResponse)(nil), model.ErrNotFound("contact info not found"))

	req := httptest.NewRequest("GET", "/contact-info/missing", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "contact info not found", response.Errors)
}

func TestContactInfoController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	reqBody := model.CreateContactInfoRequest{EntityType: "pura", Address: "A", Email: "e@x.com"}
	resBody := &model.ContactInfoResponse{ID: "1", Address: "A"}
	mockUC.On("Create", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/contact-info", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestContactInfoController_Create_BadBody(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	req := httptest.NewRequest("POST", "/contact-info", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestContactInfoController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	reqBody := model.CreateContactInfoRequest{}

	validate := validator.New()
	type Dummy struct {
		Address string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", reqBody).Return((*model.ContactInfoResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/contact-info", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestContactInfoController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	reqBody := model.UpdateContactInfoRequest{Address: "New"}
	resBody := &model.ContactInfoResponse{ID: "2", Address: "New"}
	mockUC.On("Update", "2", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/contact-info/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestContactInfoController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	reqBody := model.UpdateContactInfoRequest{Address: "A"}
	mockUC.On("Update", "3", reqBody).Return((*model.ContactInfoResponse)(nil), model.ErrNotFound("contact info not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/contact-info/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestContactInfoController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	mockUC.On("Delete", "7").Return(nil)

	req := httptest.NewRequest("DELETE", "/contact-info/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestContactInfoController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	mockUC.On("Delete", "8").Return(model.ErrNotFound("contact info not found"))

	req := httptest.NewRequest("DELETE", "/contact-info/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestContactInfoController_Delete_InternalError(t *testing.T) {
	mockUC := &usecasemock.ContactInfoUsecaseMock{}
	app := setupContactInfoController(mockUC)

	mockUC.On("Delete", "9").Return(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/contact-info/9", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal Server Error", response.Errors)
}
