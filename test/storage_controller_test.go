package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
)

type StorageUsecaseMock struct {
	mock.Mock
}

func (m *StorageUsecaseMock) UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (map[string]string, error) {
	args := m.Called(ctx, filename, file, contentType, fileSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *StorageUsecaseMock) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *StorageUsecaseMock) DeleteFile(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *StorageUsecaseMock) GetPresignedURL(ctx context.Context, key string, expiration int) (string, error) {
	args := m.Called(ctx, key, expiration)
	return args.String(0), args.Error(1)
}

func setupStorageController(mockUsecase *StorageUsecaseMock) *fiber.App {
	app, logger, _ := NewTestApp()
	controller := httpdelivery.NewStorageController(mockUsecase, logger)

	app.Post("/api/storage/upload", controller.Upload)
	app.Delete("/api/storage/delete", controller.Delete)
	app.Get("/api/storage/url", controller.GetPresignedURL)

	return app
}

func TestStorageController_Upload_Success(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "test.jpg"))
	h.Set("Content-Type", "image/jpeg")
	part, _ := writer.CreatePart(h)
	part.Write([]byte("dummy image content"))
	writer.Close()

	expectedVariants := map[string]string{
		"original": "uploads/test_123.jpg",
		"lg":       "uploads/test_123_lg.jpg",
	}

	mockUsecase.On("UploadFile", mock.Anything, "test.jpg", mock.Anything, "image/jpeg", mock.Anything).
		Return(expectedVariants, nil)

	req := httptest.NewRequest("POST", "/api/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[map[string]interface{}]
	json.NewDecoder(resp.Body).Decode(&response)

	variants := response.Data["variants"].(map[string]interface{})
	assert.Equal(t, expectedVariants["original"], variants["original"])
	assert.Equal(t, expectedVariants["lg"], variants["lg"])
	assert.Equal(t, "test.jpg", response.Data["filename"])

	mockUsecase.AssertExpectations(t)
}

func TestStorageController_Upload_NoFile(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	req := httptest.NewRequest("POST", "/api/storage/upload", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "No file uploaded", response.Errors)
}

func TestStorageController_Upload_InvalidFileType(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "doc.pdf"))
	h.Set("Content-Type", "application/pdf")
	part, _ := writer.CreatePart(h)
	part.Write([]byte("pdf content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response.Errors, "Only image files are allowed")
}

func TestStorageController_Upload_UsecaseError(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "test.jpg"))
	h.Set("Content-Type", "image/jpeg")
	part, _ := writer.CreatePart(h)
	part.Write([]byte("test content"))
	writer.Close()

	expectedError := errors.New("s3 upload failed")
	mockUsecase.On("UploadFile", mock.Anything, "test.jpg", mock.Anything, "image/jpeg", mock.Anything).
		Return((map[string]string)(nil), expectedError)

	req := httptest.NewRequest("POST", "/api/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Internal server error during file upload", response.Errors)
}

func TestStorageController_Upload_UserError(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "corrupt.jpg"))
	h.Set("Content-Type", "image/jpeg")
	part, _ := writer.CreatePart(h)
	part.Write([]byte("corrupt data"))
	writer.Close()

	expectedError := errors.New("invalid image format")
	mockUsecase.On("UploadFile", mock.Anything, "corrupt.jpg", mock.Anything, "image/jpeg", mock.Anything).
		Return((map[string]string)(nil), expectedError)

	req := httptest.NewRequest("POST", "/api/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Invalid image format or corrupted file", response.Errors)
}

func TestStorageController_Delete_Success(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	key := "uploads/test.jpg"
	mockUsecase.On("DeleteFile", mock.Anything, key).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/storage/delete?key="+key, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "File deleted successfully", response.Data)
}

func TestStorageController_Delete_NoKey(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	req := httptest.NewRequest("DELETE", "/api/storage/delete", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Key parameter is required", response.Errors)
}

func TestStorageController_Delete_Error(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	key := "uploads/error.jpg"
	expectedError := errors.New("delete failed")
	mockUsecase.On("DeleteFile", mock.Anything, key).Return(expectedError)

	req := httptest.NewRequest("DELETE", "/api/storage/delete?key="+key, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "delete failed", response.Errors)
}

func TestStorageController_GetPresignedURL_Success(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	key := "uploads/test.jpg"
	expiration := 3600
	expectedURL := "https://presigned-url.com/test"

	mockUsecase.On("GetPresignedURL", mock.Anything, key, expiration).Return(expectedURL, nil)

	req := httptest.NewRequest("GET", "/api/storage/url?key="+key, nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[fiber.Map]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, expectedURL, response.Data["url"])
}

func TestStorageController_GetPresignedURL_NoKey(t *testing.T) {
	mockUsecase := new(StorageUsecaseMock)
	app := setupStorageController(mockUsecase)

	req := httptest.NewRequest("GET", "/api/storage/url", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Key parameter is required", response.Errors)
}
