package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupArticleController(mockUC *usecasemock.ArticleUsecaseMock) *fiber.App {
	controller := httpdelivery.NewArticleController(mockUC, logrus.New())
	app := fiber.New()

	app.Get("/public/articles", controller.GetPublic)
	app.Get("/public/articles/:slug", controller.GetBySlug)

	app.Get("/articles", controller.GetAll)
	app.Post("/articles", controller.Create)
	app.Put("/articles/:id", controller.Update)
	app.Delete("/articles/:id", controller.Delete)

	return app
}

func TestArticleController_GetPublic_Success(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	now := time.Now()
	mockData := []model.ArticleResponse{
		{ID: "1", Title: "Berita A", Slug: "berita-a", PublishedAt: &now},
		{ID: "2", Title: "Berita B", Slug: "berita-b", PublishedAt: &now},
	}

	mockUC.On("GetPublic", 5).Return(mockData, nil)

	req := httptest.NewRequest("GET", "/public/articles?limit=5", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.ArticleResponse]
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Len(t, response.Data, 2)
	assert.Equal(t, "Berita A", response.Data[0].Title)
	mockUC.AssertExpectations(t)
}

func TestArticleController_GetBySlug_Success(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	slug := "upacara-besar"
	mockData := &model.ArticleResponse{
		ID: "1", Title: "Upacara Besar", Slug: slug,
	}

	mockUC.On("GetBySlug", slug).Return(mockData, nil)

	req := httptest.NewRequest("GET", "/public/articles/"+slug, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[model.ArticleResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, slug, response.Data.Slug)
}

func TestArticleController_GetBySlug_NotFound(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	slug := "tidak-ada"
	mockUC.On("GetBySlug", slug).Return((*model.ArticleResponse)(nil), errors.New("record not found"))

	req := httptest.NewRequest("GET", "/public/articles/"+slug, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestArticleController_GetAll_CMS(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	mockUC.On("GetAll", "").Return([]model.ArticleResponse{}, nil)

	req := httptest.NewRequest("GET", "/articles", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestArticleController_Create_Success(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	reqBody := model.CreateArticleRequest{
		Title:      "Judul Baru",
		AuthorName: "Admin",
		Content:    "Konten Panjang...",
		Status:     "DRAFT",
	}

	mockResp := &model.ArticleResponse{
		ID: "new-id", Title: "Judul Baru", Status: "DRAFT",
	}

	mockUC.On("Create", mock.Anything).Return(mockResp, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response model.WebResponse[model.ArticleResponse]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "new-id", response.Data.ID)
}

func TestArticleController_Create_ValidationError(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	reqBody := model.CreateArticleRequest{Title: "Judul"}

	mockUC.On("Create", mock.Anything).Return((*model.ArticleResponse)(nil), errors.New("validation failed"))

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestArticleController_Delete_Success(t *testing.T) {
	mockUC := &usecasemock.ArticleUsecaseMock{}
	app := setupArticleController(mockUC)

	targetID := "uuid-123"
	mockUC.On("Delete", targetID).Return(nil)

	req := httptest.NewRequest("DELETE", "/articles/"+targetID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
