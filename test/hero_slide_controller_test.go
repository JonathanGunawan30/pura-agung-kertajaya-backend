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

func setupHeroSlideController() (*fiber.App, *usecasemock.HeroSlideUsecaseMock) {
    mockUC := &usecasemock.HeroSlideUsecaseMock{}
    controller := httpdelivery.NewHeroSlideController(mockUC, logrus.New())
    app := fiber.New()

    app.Get("/hero-slides", controller.GetAll)
    app.Get("/hero-slides/:id", controller.GetByID)
    app.Post("/hero-slides", controller.Create)
    app.Put("/hero-slides/:id", controller.Update)
    app.Delete("/hero-slides/:id", controller.Delete)

    return app, mockUC
}

func TestHeroSlideController_GetAllPublic_Success(t *testing.T) {
    mockUC := &usecasemock.HeroSlideUsecaseMock{}
    controller := httpdelivery.NewHeroSlideController(mockUC, logrus.New())
    app := fiber.New()
    app.Get("/public/hero-slides", controller.GetAllPublic)

    items := []model.HeroSlideResponse{{ID: "a", ImageURL: "https://a"}, {ID: "b", ImageURL: "https://b"}}
    mockUC.On("GetPublic").Return(items, nil)

    req := httptest.NewRequest("GET", "/public/hero-slides", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[[]model.HeroSlideResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Len(t, response.Data, 2)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetAllPublic_Error(t *testing.T) {
    mockUC := &usecasemock.HeroSlideUsecaseMock{}
    controller := httpdelivery.NewHeroSlideController(mockUC, logrus.New())
    app := fiber.New()
    app.Get("/public/hero-slides", controller.GetAllPublic)

    mockUC.On("GetPublic").Return(([]model.HeroSlideResponse)(nil), errors.New("db error"))

    req := httptest.NewRequest("GET", "/public/hero-slides", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestHeroSlideController_GetAll_Success(t *testing.T) {
    app, mockUC := setupHeroSlideController()

    items := []model.HeroSlideResponse{{ID: "a", ImageURL: "https://a"}, {ID: "b", ImageURL: "https://b"}}
    mockUC.On("GetAll").Return(items, nil)

    req := httptest.NewRequest("GET", "/hero-slides", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[[]model.HeroSlideResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Len(t, response.Data, 2)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetAll_Error(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    mockUC.On("GetAll").Return(nil, errors.New("db error"))

    req := httptest.NewRequest("GET", "/hero-slides", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

    var response model.WebResponse[any]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "db error", response.Errors)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetByID_Success(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    item := &model.HeroSlideResponse{ID: "x", ImageURL: "https://x"}
    mockUC.On("GetByID", "x").Return(item, nil)

    req := httptest.NewRequest("GET", "/hero-slides/x", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[*model.HeroSlideResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "x", response.Data.ID)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_GetByID_NotFound(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    mockUC.On("GetByID", "missing").Return((*model.HeroSlideResponse)(nil), errors.New("not found"))

    req := httptest.NewRequest("GET", "/hero-slides/missing", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

    var response model.WebResponse[any]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "not found", response.Errors)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Create_Success(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    reqBody := model.HeroSlideRequest{ImageURL: "https://img", OrderIndex: 1, IsActive: true}
    resBody := &model.HeroSlideResponse{ID: "1", ImageURL: "https://img"}
    mockUC.On("Create", reqBody).Return(resBody, nil)

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/hero-slides", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

    var response model.WebResponse[*model.HeroSlideResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "1", response.Data.ID)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Create_BadBody(t *testing.T) {
    app, _ := setupHeroSlideController()
    req := httptest.NewRequest("POST", "/hero-slides", bytes.NewBufferString("{bad json}"))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHeroSlideController_Create_UsecaseError(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    reqBody := model.HeroSlideRequest{}
    mockUC.On("Create", reqBody).Return((*model.HeroSlideResponse)(nil), errors.New("validation failed"))

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/hero-slides", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

    var response model.WebResponse[any]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "validation failed", response.Errors)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Update_Success(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    reqBody := model.HeroSlideRequest{ImageURL: "https://new", OrderIndex: 3, IsActive: false}
    resBody := &model.HeroSlideResponse{ID: "2", ImageURL: "https://new"}
    mockUC.On("Update", "2", reqBody).Return(resBody, nil)

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("PUT", "/hero-slides/2", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[*model.HeroSlideResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "2", response.Data.ID)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Update_BadBody(t *testing.T) {
    app, _ := setupHeroSlideController()
    req := httptest.NewRequest("PUT", "/hero-slides/1", bytes.NewBufferString("{bad json}"))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHeroSlideController_Update_UsecaseError(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    reqBody := model.HeroSlideRequest{ImageURL: "https://img"}
    mockUC.On("Update", "3", reqBody).Return((*model.HeroSlideResponse)(nil), errors.New("update failed"))

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("PUT", "/hero-slides/3", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

    var response model.WebResponse[any]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "update failed", response.Errors)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Delete_Success(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    mockUC.On("Delete", "7").Return(nil)

    req := httptest.NewRequest("DELETE", "/hero-slides/7", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[string]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "Hero slide deleted successfully", response.Data)
    mockUC.AssertExpectations(t)
}

func TestHeroSlideController_Delete_UsecaseError(t *testing.T) {
    app, mockUC := setupHeroSlideController()
    mockUC.On("Delete", "8").Return(errors.New("delete failed"))

    req := httptest.NewRequest("DELETE", "/hero-slides/8", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

    var response model.WebResponse[any]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "delete failed", response.Errors)
    mockUC.AssertExpectations(t)
}
