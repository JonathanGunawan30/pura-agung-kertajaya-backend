package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupRemarkController(mockUC *usecasemock.RemarkUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewRemarkController(mockUC, logger)

	publicApi := app.Group("/api/public")
	publicApi.Get("/remarks", controller.GetAllPublic)

	api := app.Group("/api", func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})
	api.Get("/remarks", controller.GetAll)
	api.Get("/remarks/:id", controller.GetByID)
	api.Post("/remarks", controller.Create)
	api.Put("/remarks/:id", controller.Update)
	api.Delete("/remarks/:id", controller.Delete)

	return app
}

func TestRemarkController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	items := []model.RemarkResponse{
		{ID: "uuid-1", Name: "Pak Ketua", Position: "Ketua", EntityType: "pura"},
		{ID: "uuid-2", Name: "Pak Wakil", Position: "Wakil", EntityType: "pura"},
	}

	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/remarks", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetAllPublic_WithFilter_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	items := []model.RemarkResponse{}
	mockUC.On("GetPublic", "yayasan").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/remarks?entity_type=yayasan", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.RemarkResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/remarks", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestRemarkController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	items := []model.RemarkResponse{{ID: "uuid-1", Name: "A"}}
	// "pura" injected by middleware
	mockUC.On("GetAll", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/remarks", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	targetID := "uuid-10"
	item := &model.RemarkResponse{ID: targetID, Name: "Pak Bos"}

	mockUC.On("GetByID", targetID).Return(item, nil)

	req := httptest.NewRequest("GET", "/api/remarks/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.RemarkResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, targetID, response.Data.ID)
	mockUC.AssertExpectations(t)
}

func TestRemarkController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	targetID := "uuid-unknown"
	mockUC.On("GetByID", targetID).Return((*model.RemarkResponse)(nil), model.ErrNotFound("remark not found"))

	req := httptest.NewRequest("GET", "/api/remarks/"+targetID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestRemarkController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	reqBody := model.CreateRemarkRequest{
		Name:       "Test",
		Position:   "Pos",
		Content:    "Content",
		OrderIndex: 1,
	}
	resBody := &model.RemarkResponse{ID: "new-uuid-1", Name: "Test"}

	mockUC.On("Create", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/remarks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	var response model.WebResponse[*model.RemarkResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "new-uuid-1", response.Data.ID)
}

func TestRemarkController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	reqBody := model.CreateRemarkRequest{}

	validate := validator.New()
	type Dummy struct {
		Name string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.RemarkResponse)(nil), realValErr)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/remarks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRemarkController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	idToUpdate := "uuid-5"
	reqBody := model.UpdateRemarkRequest{
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

	mockUC.On("Update", idToUpdate, reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/remarks/"+idToUpdate, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var response model.WebResponse[*model.RemarkResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Updated Name", response.Data.Name)
}

func TestRemarkController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	idToUpdate := "uuid-missing"
	reqBody := model.UpdateRemarkRequest{Name: "Update"}

	mockUC.On("Update", idToUpdate, reqBody).Return((*model.RemarkResponse)(nil), model.ErrNotFound("remark not found"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/remarks/"+idToUpdate, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestRemarkController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	idToDelete := "uuid-7"
	mockUC.On("Delete", idToDelete).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/remarks/"+idToDelete, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRemarkController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.RemarkUsecaseMock{}
	app := setupRemarkController(mockUC)

	idToDelete := "uuid-missing"
	mockUC.On("Delete", idToDelete).Return(model.ErrNotFound("remark not found"))

	req := httptest.NewRequest("DELETE", "/api/remarks/"+idToDelete, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
