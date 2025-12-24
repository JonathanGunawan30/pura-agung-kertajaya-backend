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

func setupRemarkController() (*fiber.App, *usecasemock.RemarkUsecaseMock) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	controller := httpdelivery.NewRemarkController(mockUC, logrus.New())

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "An internal server error occurred."

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
	api.Get("/remarks", controller.GetAll)
	api.Get("/remarks/:id", controller.GetByID)
	api.Post("/remarks", controller.Create)
	api.Put("/remarks/:id", controller.Update)
	api.Delete("/remarks/:id", controller.Delete)

	publicApi := app.Group("/api/public")
	publicApi.Get("/remarks", controller.GetAllPublic)

	return app, mockUC
}

func TestRemarkController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupRemarkController()

	items := []model.RemarkResponse{
		{ID: "uuid-1", Name: "Pak Ketua", Position: "Ketua", EntityType: "pura"},
		{ID: "uuid-2", Name: "Pak Wakil", Position: "Wakil", EntityType: "pura"},
	}

	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/remarks", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetAllPublic_WithFilter_Success(t *testing.T) {
	app, mockUC := setupRemarkController()
	items := []model.RemarkResponse{}

	mockUC.On("GetPublic", "yayasan").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/remarks?entity_type=yayasan", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupRemarkController()
	mockUC.On("GetPublic", "").Return(([]model.RemarkResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/remarks", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetAll_Success(t *testing.T) {
	app, mockUC := setupRemarkController()
	items := []model.RemarkResponse{{ID: "uuid-1", Name: "A"}}

	mockUC.On("GetAll", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/remarks", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetByID_Success(t *testing.T) {
	app, mockUC := setupRemarkController()
	targetID := "uuid-10"
	item := &model.RemarkResponse{ID: targetID, Name: "Pak Bos"}

	mockUC.On("GetByID", targetID).Return(item, nil)

	req := httptest.NewRequest("GET", "/api/remarks/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.RemarkResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, targetID, response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupRemarkController()
	targetID := "uuid-unknown"
	mockUC.On("GetByID", targetID).Return((*model.RemarkResponse)(nil), gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/api/remarks/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_Create_Success(t *testing.T) {
	app, mockUC := setupRemarkController()

	reqBody := model.CreateRemarkRequest{
		EntityType: "pura",
		Name:       "New Person", Position: "Staff", Content: "Hello", OrderIndex: 1,
	}
	resBody := &model.RemarkResponse{ID: "new-uuid-1", Name: "New Person", EntityType: "pura"}

	mockUC.On("Create", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/remarks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	var response model.WebResponse[*model.RemarkResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "new-uuid-1", response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_Update_Success(t *testing.T) {
	app, mockUC := setupRemarkController()
	idToUpdate := "uuid-5"

	reqBodyRaw := model.UpdateRemarkRequest{
		Name:       "Updated Name",
		Position:   "Updated Pos",
		Content:    "Updated Content",
		IsActive:   true,
		OrderIndex: 2,
	}

	resBody := &model.RemarkResponse{
		ID:         idToUpdate,
		Name:       "Updated Name",
		EntityType: "pura",
	}

	mockUC.On("Update", idToUpdate, reqBodyRaw).Return(resBody, nil)

	b, _ := json.Marshal(reqBodyRaw)
	req := httptest.NewRequest("PUT", "/api/remarks/"+idToUpdate, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.RemarkResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Updated Name", response.Data.Name)

	mockUC.AssertExpectations(t)
}
func TestRemarkController_Delete_Success(t *testing.T) {
	app, mockUC := setupRemarkController()
	idToDelete := "uuid-7"

	mockUC.On("Delete", idToDelete).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/remarks/"+idToDelete, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}
