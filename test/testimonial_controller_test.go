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

func setupTestimonialController() (*fiber.App, *usecasemock.TestimonialUsecaseMock) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	controller := httpdelivery.NewTestimonialController(mockUC, logrus.New())
	app := fiber.New()

	app.Get("/testimonials", controller.GetAll)
	app.Get("/testimonials/:id", controller.GetByID)
	app.Post("/testimonials", controller.Create)
	app.Put("/testimonials/:id", controller.Update)
	app.Delete("/testimonials/:id", controller.Delete)

	return app, mockUC
}

func TestTestimonialController_GetAll_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()

	items := []model.TestimonialResponse{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	mockUC.On("GetAll").Return(items, nil)

	req := httptest.NewRequest("GET", "/testimonials", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.TestimonialResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetAll_Error(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("GetAll").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/testimonials", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "db error", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetByID_Valid_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	item := &model.TestimonialResponse{ID: 10, Name: "X"}
	mockUC.On("GetByID", 10).Return(item, nil)

	req := httptest.NewRequest("GET", "/testimonials/10", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.TestimonialResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 10, response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetByID_InvalidID(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("GET", "/testimonials/abc", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("GetByID", 999).Return((*model.TestimonialResponse)(nil), errors.New("not found"))

	req := httptest.NewRequest("GET", "/testimonials/999", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "not found", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Create_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	reqBody := model.TestimonialRequest{Name: "John", Rating: 5, Comment: "Good"}
	resBody := &model.TestimonialResponse{ID: 1, Name: "John"}
	mockUC.On("Create", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/testimonials", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response model.WebResponse[*model.TestimonialResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 1, response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Create_BadBody(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("POST", "/testimonials", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupTestimonialController()
	reqBody := model.TestimonialRequest{Name: "", Rating: 0}
	mockUC.On("Create", reqBody).Return((*model.TestimonialResponse)(nil), errors.New("validation failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/testimonials", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "validation failed", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Update_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	reqBody := model.TestimonialRequest{Name: "Jane", Rating: 4, Comment: "Nice"}
	resBody := &model.TestimonialResponse{ID: 2, Name: "Jane"}
	mockUC.On("Update", 2, reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/testimonials/2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.TestimonialResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 2, response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Update_InvalidID(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("PUT", "/testimonials/zero", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Update_BadBody(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("PUT", "/testimonials/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupTestimonialController()
	reqBody := model.TestimonialRequest{Name: "Jane", Rating: 6}
	mockUC.On("Update", 3, reqBody).Return((*model.TestimonialResponse)(nil), errors.New("update failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/testimonials/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "update failed", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Delete_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("Delete", 7).Return(nil)

	req := httptest.NewRequest("DELETE", "/testimonials/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Testimonial deleted successfully", response.Data)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Delete_InvalidID(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("DELETE", "/testimonials/zero", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("Delete", 8).Return(errors.New("delete failed"))

	req := httptest.NewRequest("DELETE", "/testimonials/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "delete failed", response.Errors)
	mockUC.AssertExpectations(t)
}
