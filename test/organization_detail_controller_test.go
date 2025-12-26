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

func setupOrganizationDetailController() (*fiber.App, *usecasemock.OrganizationDetailUsecaseMock) {
	mockUC := &usecasemock.OrganizationDetailUsecaseMock{}
	controller := httpdelivery.NewOrganizationDetailController(mockUC, logrus.New())

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
	api.Get("/organization-details", controller.GetAdmin)
	api.Put("/organization-details", controller.Update)

	publicApi := app.Group("/api/public")
	publicApi.Get("/organization-details", controller.GetPublic)

	return app, mockUC
}

func TestOrganizationDetailController_GetPublic_Success(t *testing.T) {
	app, mockUC := setupOrganizationDetailController()

	expectedResp := &model.OrganizationDetailResponse{
		ID: "uuid-1", EntityType: "pura", Vision: "Visi Pura",
	}

	mockUC.On("GetByEntityType", "pura").Return(expectedResp, nil)

	req := httptest.NewRequest("GET", "/api/public/organization-details?entity_type=pura", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationDetailResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Visi Pura", response.Data.Vision)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_GetPublic_NotFound(t *testing.T) {
	app, mockUC := setupOrganizationDetailController()

	mockUC.On("GetByEntityType", "pasraman").Return((*model.OrganizationDetailResponse)(nil), gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/api/public/organization-details?entity_type=pasraman", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_GetAdmin_Success(t *testing.T) {
	app, mockUC := setupOrganizationDetailController()

	expectedResp := &model.OrganizationDetailResponse{
		ID: "uuid-admin-1", EntityType: "yayasan", Mission: "Misi Yayasan",
	}

	mockUC.On("GetByEntityType", "yayasan").Return(expectedResp, nil)

	req := httptest.NewRequest("GET", "/api/organization-details?entity_type=yayasan", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationDetailResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Misi Yayasan", response.Data.Mission)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_Update_Success(t *testing.T) {
	app, mockUC := setupOrganizationDetailController()

	reqBody := model.UpdateOrganizationDetailRequest{
		Vision: "Visi Baru", Mission: "Misi Baru", Rules: "Rules Baru",
	}

	resBody := &model.OrganizationDetailResponse{
		ID: "uuid-upserted", EntityType: "pura", Vision: "Visi Baru",
	}

	mockUC.On("Update", "pura", reqBody).Return(resBody, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-details?entity_type=pura", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationDetailResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Visi Baru", response.Data.Vision)

	mockUC.AssertExpectations(t)
}

func TestOrganizationDetailController_Update_InvalidBody(t *testing.T) {
	app, mockUC := setupOrganizationDetailController()

	req := httptest.NewRequest("PUT", "/api/organization-details?entity_type=pura", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUC.AssertNotCalled(t, "Update")
}

func TestOrganizationDetailController_Update_UsecaseError(t *testing.T) {
	app, mockUC := setupOrganizationDetailController()

	reqBody := model.UpdateOrganizationDetailRequest{Vision: "Test"}

	mockUC.On("Update", "yayasan", reqBody).Return((*model.OrganizationDetailResponse)(nil), errors.New("db connection failed"))

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-details?entity_type=yayasan", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUC.AssertExpectations(t)
}
