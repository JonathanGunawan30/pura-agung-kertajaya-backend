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

func setupAboutController() (*fiber.App, *usecasemock.AboutUsecaseMock) {
    mockUC := &usecasemock.AboutUsecaseMock{}
    controller := httpdelivery.NewAboutController(mockUC, logrus.New())
    app := fiber.New(fiber.Config{StrictRouting: true})

    app.Get("/about", controller.GetAll)
    app.Get("/about/:id", controller.GetByID)
    app.Post("/about", controller.Create)
    app.Put("/about/:id", controller.Update)
    app.Delete("/about/:id", controller.Delete)

    return app, mockUC
}

func TestAboutController_GetAllPublic_Success(t *testing.T) {
    mockUC := &usecasemock.AboutUsecaseMock{}
    controller := httpdelivery.NewAboutController(mockUC, logrus.New())
    app := fiber.New()
    app.Get("/public/about", controller.GetAllPublic)

    items := []model.AboutSectionResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
    mockUC.On("GetPublic").Return(items, nil)

    req := httptest.NewRequest("GET", "/public/about", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[[]model.AboutSectionResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Len(t, response.Data, 2)
    mockUC.AssertExpectations(t)
}

func TestAboutController_GetAllPublic_Error(t *testing.T) {
    mockUC := &usecasemock.AboutUsecaseMock{}
    controller := httpdelivery.NewAboutController(mockUC, logrus.New())
    app := fiber.New()
    app.Get("/public/about", controller.GetAllPublic)

    mockUC.On("GetPublic").Return(([]model.AboutSectionResponse)(nil), errors.New("db error"))

    req := httptest.NewRequest("GET", "/public/about", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestAboutController_GetAll_Success(t *testing.T) {
    app, mockUC := setupAboutController()
    items := []model.AboutSectionResponse{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
    mockUC.On("GetAll").Return(items, nil)

    req := httptest.NewRequest("GET", "/about", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[[]model.AboutSectionResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Len(t, response.Data, 2)
    mockUC.AssertExpectations(t)
}

func TestAboutController_GetAll_Error(t *testing.T) {
    app, mockUC := setupAboutController()
    mockUC.On("GetAll").Return(nil, errors.New("db error"))

    req := httptest.NewRequest("GET", "/about", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestAboutController_GetByID_Success(t *testing.T) {
    app, mockUC := setupAboutController()
    item := &model.AboutSectionResponse{ID: "x", Title: "X"}
    mockUC.On("GetByID", "x").Return(item, nil)

    req := httptest.NewRequest("GET", "/about/x", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[*model.AboutSectionResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "x", response.Data.ID)
    mockUC.AssertExpectations(t)
}

func TestAboutController_GetByID_NotFound(t *testing.T) {
    app, mockUC := setupAboutController()
    mockUC.On("GetByID", "missing").Return((*model.AboutSectionResponse)(nil), errors.New("not found"))

    req := httptest.NewRequest("GET", "/about/missing", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestAboutController_Create_Success(t *testing.T) {
    app, mockUC := setupAboutController()
    reqBody := model.AboutSectionRequest{Title: "T", Description: "D"}
    resBody := &model.AboutSectionResponse{ID: "1", Title: "T"}
    mockUC.On("Create", reqBody).Return(resBody, nil)

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/about", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

    var response model.WebResponse[*model.AboutSectionResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "1", response.Data.ID)
    mockUC.AssertExpectations(t)
}

func TestAboutController_Create_BadBody(t *testing.T) {
    app, _ := setupAboutController()
    req := httptest.NewRequest("POST", "/about", bytes.NewBufferString("{bad json}"))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAboutController_Create_UsecaseError(t *testing.T) {
    app, mockUC := setupAboutController()
    reqBody := model.AboutSectionRequest{}
    mockUC.On("Create", reqBody).Return((*model.AboutSectionResponse)(nil), errors.New("validation failed"))

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/about", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAboutController_Update_Success(t *testing.T) {
    app, mockUC := setupAboutController()
    reqBody := model.AboutSectionRequest{Title: "New", Description: "D"}
    resBody := &model.AboutSectionResponse{ID: "2", Title: "New"}
    mockUC.On("Update", "2", reqBody).Return(resBody, nil)

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("PUT", "/about/2", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[*model.AboutSectionResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "2", response.Data.ID)
    mockUC.AssertExpectations(t)
}

func TestAboutController_Update_BadBody(t *testing.T) {
    app, _ := setupAboutController()
    req := httptest.NewRequest("PUT", "/about/1", bytes.NewBufferString("{bad json}"))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAboutController_Update_UsecaseError(t *testing.T) {
    app, mockUC := setupAboutController()
    reqBody := model.AboutSectionRequest{Title: "N"}
    mockUC.On("Update", "3", reqBody).Return((*model.AboutSectionResponse)(nil), errors.New("update failed"))

    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("PUT", "/about/3", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestAboutController_Delete_Success(t *testing.T) {
    app, mockUC := setupAboutController()
    mockUC.On("Delete", "7").Return(nil)

    req := httptest.NewRequest("DELETE", "/about/7", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response model.WebResponse[string]
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "About section deleted successfully", response.Data)
    mockUC.AssertExpectations(t)
}

func TestAboutController_Delete_UsecaseError(t *testing.T) {
    app, mockUC := setupAboutController()
    mockUC.On("Delete", "8").Return(errors.New("delete failed"))

    req := httptest.NewRequest("DELETE", "/about/8", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
