package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupOrganizationDetailController(mockUC *usecasemock.OrganizationDetailUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewOrganizationDetailController(mockUC, logger)

	publicApi := app.Group("/api/public")
	publicApi.Get("/organization-details", controller.GetPublic)

	api := app.Group("/api", func(c *fiber.Ctx) error {
		c.Locals(middleware.CtxEntityType, "pura")
		return c.Next()
	})
	api.Get("/organization-details", controller.GetAdmin)
	api.Put("/organization-details", controller.Update)

	return app
}

func TestOrganizationDetailController_GetPublic_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	app := setupOrganizationDetailController(mockUC)

	expectedResp := &model.OrganizationDetailResponse{
		ID:         "uuid-1",
		EntityType: "pura",
		Vision:     "Visi Pura",
	}

	mockUC.On("GetByEntityType", "pura").Return(expectedResp, nil)

	req := httptest.NewRequest("GET", "/api/public/organization-details?entity_type=pura", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationDetailResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Visi Pura", response.Data.Vision)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_GetPublic_Error(t *testing.T) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	app := setupOrganizationDetailController(mockUC)

	mockUC.On("GetByEntityType", "pura").Return((*model.OrganizationDetailResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/organization-details?entity_type=pura", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestOrganizationDetailController_GetAdmin_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	app := setupOrganizationDetailController(mockUC)

	expectedResp := &model.OrganizationDetailResponse{
		ID:         "uuid-admin-1",
		EntityType: "pura",
		Mission:    "Misi Yayasan",
	}

	mockUC.On("GetByEntityType", "pura").Return(expectedResp, nil)

	req := httptest.NewRequest("GET", "/api/organization-details", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationDetailResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Misi Yayasan", response.Data.Mission)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_Update_Success(t *testing.T) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	app := setupOrganizationDetailController(mockUC)

	reqBody := model.UpdateOrganizationDetailRequest{
		Vision:                "Visi Baru",
		Mission:               "Misi Baru",
		VisionMissionImageURL: "new.jpg",
	}

	resBody := &model.OrganizationDetailResponse{
		ID:                    "uuid-upserted",
		EntityType:            "pura",
		Vision:                "Visi Baru",
		VisionMissionImageURL: "new.jpg",
	}

	mockUC.On("Update", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-details", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationDetailResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Visi Baru", response.Data.Vision)
	assert.Equal(t, "new.jpg", response.Data.VisionMissionImageURL)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_Update_InvalidBody(t *testing.T) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	app := setupOrganizationDetailController(mockUC)

	req := httptest.NewRequest("PUT", "/api/organization-details", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUC.AssertNotCalled(t, "Update")
}

func TestOrganizationDetailController_Update_UsecaseError(t *testing.T) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	app := setupOrganizationDetailController(mockUC)

	reqBody := model.UpdateOrganizationDetailRequest{Vision: "Test"}

	mockUC.On("Update", "pura", reqBody).Return((*model.OrganizationDetailResponse)(nil), errors.New("db connection failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-details", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUC.AssertExpectations(t)
}
