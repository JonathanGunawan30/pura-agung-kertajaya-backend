package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
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
	app := fiber.New(fiber.Config{
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
	api.Get("/testimonials", controller.GetAll)
	api.Get("/testimonials/:id", controller.GetByID)
	api.Post("/testimonials", controller.Create)
	api.Put("/testimonials/:id", controller.Update)
	api.Delete("/testimonials/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/testimonials", controller.GetAllPublic)

	return app, mockUC
}

func TestTestimonialController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	items := []model.TestimonialResponse{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	mockUC.On("GetPublic").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/public/testimonials", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("GetPublic").Return(([]model.TestimonialResponse)(nil), errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/public/testimonials", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetAll_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	items := []model.TestimonialResponse{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	mockUC.On("GetAll").Return(items, nil)
	req := httptest.NewRequest("GET", "/api/testimonials", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetAll_Error(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("GetAll").Return(nil, errors.New("db error"))
	req := httptest.NewRequest("GET", "/api/testimonials", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetByID_Valid_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	item := &model.TestimonialResponse{ID: 10, Name: "X"}
	mockUC.On("GetByID", 10).Return(item, nil)
	req := httptest.NewRequest("GET", "/api/testimonials/10", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var response model.WebResponse[*model.TestimonialResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 10, response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_GetByID_InvalidID(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("GET", "/api/testimonials/abc", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupTestimonialController()
	mockUC.On("GetByID", 999).Return((*model.TestimonialResponse)(nil), gorm.ErrRecordNotFound) // Mock terima string
	req := httptest.NewRequest("GET", "/api/testimonials/999", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "The requested resource was not found.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Create_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	reqBody := model.TestimonialRequest{Name: "John", Rating: 5, Comment: "Good", IsActive: true, OrderIndex: 1}
	resBody := &model.TestimonialResponse{ID: 1, Name: "John"}
	mockUC.On("Create", reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/testimonials", bytes.NewReader(b))
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
	req := httptest.NewRequest("POST", "/api/testimonials", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupTestimonialController()
	reqBody := model.TestimonialRequest{Name: "", Rating: 0}
	validate := validator.New()
	err := validate.Struct(reqBody)
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)
	mockUC.On("Create", reqBody).Return((*model.TestimonialResponse)(nil), validationErrs)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/testimonials", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Update_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	idToUpdate := 2
	reqBody := model.TestimonialRequest{Name: "Jane", Rating: 4, Comment: "Nice", IsActive: true, OrderIndex: 1}
	resBody := &model.TestimonialResponse{ID: 2, Name: "Jane"}
	mockUC.On("Update", idToUpdate, reqBody).Return(resBody, nil)
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/testimonials/2", bytes.NewReader(b))
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
	req := httptest.NewRequest("PUT", "/api/testimonials/abc", bytes.NewBufferString("{}")) // ID non-numerik
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Update_BadBody(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("PUT", "/api/testimonials/1", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupTestimonialController()
	idToUpdate := 3
	reqBody := model.TestimonialRequest{Name: "Jane", Rating: 6}
	mockUC.On("Update", idToUpdate, reqBody).Return((*model.TestimonialResponse)(nil), errors.New("update failed"))
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/testimonials/3", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Delete_Success(t *testing.T) {
	app, mockUC := setupTestimonialController()
	idToDelete := 7
	mockUC.On("Delete", idToDelete).Return(nil)
	req := httptest.NewRequest("DELETE", "/api/testimonials/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Testimonial deleted successfully", response.Data)
	mockUC.AssertExpectations(t)
}

func TestTestimonialController_Delete_InvalidID(t *testing.T) {
	app, _ := setupTestimonialController()
	req := httptest.NewRequest("DELETE", "/api/testimonials/abc", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTestimonialController_Delete_UsecaseError(t *testing.T) {
	app, mockUC := setupTestimonialController()
	idToDelete := 8
	mockUC.On("Delete", idToDelete).Return(errors.New("delete failed"))
	req := httptest.NewRequest("DELETE", "/api/testimonials/8", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)
	mockUC.AssertExpectations(t)
}
