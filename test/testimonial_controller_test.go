package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestimonialController(mockUC *usecasemock.TestimonialUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewTestimonialController(mockUC, logger)

	publicApi := app.Group("/api/public")
	publicApi.Get("/testimonials", controller.GetAllPublic)

	api := app.Group("/api")
	api.Get("/testimonials", controller.GetAll)
	api.Get("/testimonials/:id", controller.GetByID)
	api.Post("/testimonials", controller.Create)
	api.Put("/testimonials/:id", controller.Update)
	api.Delete("/testimonials/:id", controller.Delete)

	return app
}

func TestTestimonialController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	items := []model.TestimonialResponse{
		{ID: "uuid-1", Name: "A"},
		{ID: "uuid-2", Name: "B"},
	}
	mockUC.On("GetPublic").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/testimonials", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	mockUC.On("GetPublic").Return(([]model.TestimonialResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/testimonials", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestTestimonialController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	items := []model.TestimonialResponse{
		{ID: "uuid-1", Name: "A"},
		{ID: "uuid-2", Name: "B"},
	}
	mockUC.On("GetAll").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/testimonials", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestTestimonialController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	targetID := "uuid-10"
	item := &model.TestimonialResponse{ID: targetID, Name: "X"}

	mockUC.On("GetByID", targetID).Return(item, nil)

	req := httptest.NewRequest("GET", "/api/testimonials/uuid-10", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.TestimonialResponse]
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, targetID, response.Data.ID)
}

func TestTestimonialController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	targetID := "non-existent-uuid"
	mockUC.On("GetByID", targetID).Return((*model.TestimonialResponse)(nil), model.ErrNotFound("testimonial not found"))

	req := httptest.NewRequest("GET", "/api/testimonials/non-existent-uuid", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestTestimonialController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	reqBody := model.TestimonialRequest{Name: "John", Rating: 5, Comment: "Good", IsActive: true, OrderIndex: 1}
	resBody := &model.TestimonialResponse{ID: "new-uuid-1", Name: "John"}

	mockUC.On("Create", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/testimonials", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestTestimonialController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	reqBody := model.TestimonialRequest{}

	validate := validator.New()
	type Dummy struct {
		Name string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", reqBody).Return((*model.TestimonialResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/testimonials", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestTestimonialController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	idToUpdate := "uuid-2"
	reqBody := model.TestimonialRequest{Name: "Jane", Rating: 4, Comment: "Nice", IsActive: true, OrderIndex: 1}
	resBody := &model.TestimonialResponse{ID: idToUpdate, Name: "Jane"}

	mockUC.On("Update", idToUpdate, reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/testimonials/uuid-2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestTestimonialController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	idToUpdate := "non-existent-uuid"
	reqBody := model.TestimonialRequest{Name: "Jane"}

	mockUC.On("Update", idToUpdate, reqBody).Return((*model.TestimonialResponse)(nil), model.ErrNotFound("testimonial not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/testimonials/non-existent-uuid", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestTestimonialController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	idToDelete := "uuid-7"
	mockUC.On("Delete", idToDelete).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/testimonials/uuid-7", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestTestimonialController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.TestimonialUsecaseMock{}
	app := setupTestimonialController(mockUC)

	idToDelete := "uuid-8"
	mockUC.On("Delete", idToDelete).Return(model.ErrNotFound("testimonial not found"))

	req := httptest.NewRequest("DELETE", "/api/testimonials/uuid-8", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
