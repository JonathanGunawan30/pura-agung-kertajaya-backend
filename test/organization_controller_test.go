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

func setupOrganizationController(mockUC *usecasemock.OrganizationMemberUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewOrganizationController(mockUC, logger)

	publicApi := app.Group("/api/public")
	publicApi.Get("/organization-members", controller.GetAllPublic)

	api := app.Group("/api", func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})
	api.Get("/organization-members", controller.GetAll)
	api.Get("/organization-members/:id", controller.GetByID)
	api.Post("/organization-members", controller.Create)
	api.Put("/organization-members/:id", controller.Update)
	api.Delete("/organization-members/:id", controller.Delete)

	return app
}

func TestOrganizationController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	items := []model.OrganizationResponse{
		{ID: "1", Name: "Member A", Position: "Ketua", PositionOrder: 1, IsActive: true},
		{ID: "2", Name: "Member B", Position: "Sekretaris", PositionOrder: 2, IsActive: true},
	}
	mockUC.On("GetPublic", "pura").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/organization-members?entity_type=pura", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.OrganizationResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, "Member A", response.Data[0].Name)

	mockUC.AssertExpectations(t)
}

func TestOrganizationController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	mockUC.On("GetPublic", "").Return(([]model.OrganizationResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/organization-members", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestOrganizationController_GetAll_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	mockResponse := []model.OrganizationResponse{
		{ID: "1", Name: "Member A"},
	}
	mockUC.On("GetAll", "pura").Return(mockResponse, nil)

	req := httptest.NewRequest("GET", "/api/organization-members", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestOrganizationController_GetByID_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "member123"
	mockResponse := &model.OrganizationResponse{ID: memberID, Name: "Member Found"}
	mockUC.On("GetByID", memberID).Return(mockResponse, nil)

	req := httptest.NewRequest("GET", "/api/organization-members/"+memberID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUC.AssertExpectations(t)
}

func TestOrganizationController_GetByID_NotFound(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "notfound"
	mockUC.On("GetByID", memberID).Return((*model.OrganizationResponse)(nil), model.ErrNotFound("member not found"))

	req := httptest.NewRequest("GET", "/api/organization-members/"+memberID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestOrganizationController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	reqBody := model.CreateOrganizationRequest{
		Name:          "New Member",
		Position:      "Anggota",
		PositionOrder: 5,
		IsActive:      true,
	}
	mockResponse := &model.OrganizationResponse{
		ID:   "newID",
		Name: "New Member",
	}

	mockUC.On("Create", "pura", reqBody).Return(mockResponse, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organization-members", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestOrganizationController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	reqBody := model.CreateOrganizationRequest{}

	// Mock validator error
	validate := validator.New()
	type Dummy struct {
		Name string `validate:"required"`
	}
	realValErr := validate.Struct(Dummy{})

	mockUC.On("Create", "pura", reqBody).Return((*model.OrganizationResponse)(nil), realValErr)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organization-members", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestOrganizationController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "updateID"
	reqBody := model.UpdateOrganizationRequest{
		Name:          "Updated Name",
		Position:      "Ketua",
		PositionOrder: 1,
		IsActive:      true,
	}
	mockResponse := &model.OrganizationResponse{
		ID:   memberID,
		Name: "Updated Name",
	}

	mockUC.On("Update", memberID, reqBody).Return(mockResponse, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationController_Update_NotFound(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "notfound"
	reqBody := model.UpdateOrganizationRequest{Name: "Update"}

	mockUC.On("Update", memberID, reqBody).Return((*model.OrganizationResponse)(nil), model.ErrNotFound("member not found"))

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestOrganizationController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "deleteID"
	mockUC.On("Delete", memberID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationController_Delete_NotFound(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "notfound"
	mockUC.On("Delete", memberID).Return(model.ErrNotFound("member not found"))

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestOrganizationController_Delete_InternalError(t *testing.T) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	app := setupOrganizationController(mockUC)

	memberID := "errorID"
	mockUC.On("Delete", memberID).Return(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
