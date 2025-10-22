package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestStorageController_Upload_Success(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()

	// multipart body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte("test file content"))
	writer.Close()

	expectedURL := "https://example.com/uploads/test_1234567890.jpg"
	mockUsecase.On("UploadFile", mock.Anything, "test.jpg", mock.Anything, mock.Anything, mock.Anything).
		Return(expectedURL, nil)

	app.Post("/upload", controller.Upload)
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[fiber.Map]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, expectedURL, response.Data["url"])
	assert.Equal(t, "test.jpg", response.Data["filename"])
	mockUsecase.AssertExpectations(t)
}

func TestStorageController_Upload_NoFile(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()
	app.Post("/upload", controller.Upload)

	req := httptest.NewRequest("POST", "/upload", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "No file uploaded", response.Errors)
}

func TestStorageController_Upload_UsecaseError(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte("test file content"))
	writer.Close()

	expectedError := errors.New("upload failed")
	mockUsecase.On("UploadFile", mock.Anything, "test.jpg", mock.Anything, mock.Anything, mock.Anything).
		Return("", expectedError)

	app.Post("/upload", controller.Upload)
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "upload failed", response.Errors)
	mockUsecase.AssertExpectations(t)
}

func TestStorageController_Delete_Success(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()

	key := "uploads/test_1234567890.jpg"
	mockUsecase.On("DeleteFile", mock.Anything, key).Return(nil)

	app.Delete("/delete", controller.Delete)
	req := httptest.NewRequest("DELETE", "/delete?key="+key, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[string]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "File deleted successfully", response.Data)
	mockUsecase.AssertExpectations(t)
}

func TestStorageController_Delete_NoKey(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()

	app.Delete("/delete", controller.Delete)
	req := httptest.NewRequest("DELETE", "/delete", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Key parameter is required", response.Errors)
}

func TestStorageController_Delete_Error(t *testing.T) {
	mockUsecase := &usecasemock.MockStorageUsecase{}
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "An internal server error occurred. Please try again later."
			if errors.Is(err, gorm.ErrRecordNotFound) {
				code = fiber.StatusNotFound
				message = "The requested resource was not found."
			}
			return ctx.Status(code).JSON(model.WebResponse[any]{Errors: message})
		},
	})

	app.Delete("/api/storage/delete", controller.Delete)

	key := "uploads/test_1234567890.jpg"
	expectedError := errors.New("delete failed")
	mockUsecase.On("DeleteFile", mock.Anything, key).Return(expectedError)

	req := httptest.NewRequest("DELETE", "/api/storage/delete?key="+key, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response model.WebResponse[any]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "An internal server error occurred. Please try again later.", response.Errors)

	mockUsecase.AssertExpectations(t)
}

func TestStorageController_GetPresignedURL_Success(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()

	key := "uploads/test_1234567890.jpg"
	expiration := 3600
	expectedURL := "https://presigned-url.com/test"

	mockUsecase.On("GetPresignedURL", mock.Anything, key, expiration).Return(expectedURL, nil)

	app.Get("/presigned-url", controller.GetPresignedURL)
	req := httptest.NewRequest("GET", "/presigned-url?key="+key, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[fiber.Map]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, expectedURL, response.Data["url"])
	mockUsecase.AssertExpectations(t)
}

func TestStorageController_GetPresignedURL_NoKey(t *testing.T) {
	mockUsecase := usecasemock.NewMockStorageUsecase()
	log := logrus.New()
	controller := http.NewStorageController(mockUsecase, log)
	app := fiber.New()

	app.Get("/presigned-url", controller.GetPresignedURL)
	req := httptest.NewRequest("GET", "/presigned-url", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response model.WebResponse[any]
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Key parameter is required", response.Errors)
}
