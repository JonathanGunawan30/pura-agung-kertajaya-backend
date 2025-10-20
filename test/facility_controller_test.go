package test

import (
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

func TestFacilityController_GetAllPublic_Success(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	controller := httpdelivery.NewFacilityController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/facilities", controller.GetAllPublic)

	items := []model.FacilityResponse{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	mockUC.On("GetPublic").Return(items, nil)

	req := httptest.NewRequest("GET", "/public/facilities", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.FacilityResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 2)
	mockUC.AssertExpectations(t)
}

func TestFacilityController_GetAllPublic_Error(t *testing.T) {
	mockUC := &usecasemock.FacilityUsecaseMock{}
	controller := httpdelivery.NewFacilityController(mockUC, logrus.New())
	app := fiber.New()
	app.Get("/public/facilities", controller.GetAllPublic)

	mockUC.On("GetPublic").Return(([]model.FacilityResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/public/facilities", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
