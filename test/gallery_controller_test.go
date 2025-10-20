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

func TestGalleryController_GetAllPublic_Success(t *testing.T) {
    mockUC := &usecasemock.GalleryUsecaseMock{}
    controller := httpdelivery.NewGalleryController(mockUC, logrus.New())
    app := fiber.New()
    app.Get("/public/galleries", controller.GetAllPublic)

    items := []model.GalleryResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
    mockUC.On("GetPublic").Return(items, nil)

    req := httptest.NewRequest("GET", "/public/galleries", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[[]model.GalleryResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Len(t, response.Data, 2)
    mockUC.AssertExpectations(t)
}

func TestGalleryController_GetAllPublic_Error(t *testing.T) {
    mockUC := &usecasemock.GalleryUsecaseMock{}
    controller := httpdelivery.NewGalleryController(mockUC, logrus.New())
    app := fiber.New()
    app.Get("/public/galleries", controller.GetAllPublic)

    mockUC.On("GetPublic").Return(([]model.GalleryResponse)(nil), errors.New("db error"))

    req := httptest.NewRequest("GET", "/public/galleries", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
