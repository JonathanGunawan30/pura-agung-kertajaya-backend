package test

import (
	"context"
	"errors"
	"io"
	mock2 "pura-agung-kertajaya-backend/internal/repository/mock"
	"strings"
	"testing"

	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStorageUsecase_UploadFile_Success(t *testing.T) {
	// Setup
	mockRepo := mock2.NewMockStorageRepository()
	log := logrus.New()
	validate := validator.New()
	u := usecase.NewStorageUsecase(mockRepo, log, validate)

	ctx := context.Background()
	filename := "test.jpg"
	file := strings.NewReader("test file content")
	contentType := "image/jpeg"
	fileSize := int64(len("test file content"))
	expectedURL := "https://example.com/uploads/test_1234567890.jpg"

	// Mock expectation
	mockRepo.On("Upload", ctx, mock.AnythingOfType("string"), file, contentType, fileSize).
		Return(expectedURL, nil)

	// Execute
	url, err := u.UploadFile(ctx, filename, file, contentType, fileSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_UploadFile_Error(t *testing.T) {
	// Setup
	mockRepo := mock2.NewMockStorageRepository()
	log := logrus.New()
	validate := validator.New()
	u := usecase.NewStorageUsecase(mockRepo, log, validate)

	ctx := context.Background()
	filename := "test.jpg"
	file := strings.NewReader("test file content")
	contentType := "image/jpeg"
	fileSize := int64(len("test file content"))
	expectedError := errors.New("upload failed")

	// Mock expectation
	mockRepo.On("Upload", ctx, mock.AnythingOfType("string"), file, contentType, fileSize).
		Return("", expectedError)

	// Execute
	url, err := u.UploadFile(ctx, filename, file, contentType, fileSize)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", url)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_DeleteFile_Success(t *testing.T) {
	// Setup
	mockRepo := mock2.NewMockStorageRepository()
	log := logrus.New()
	validate := validator.New()
	u := usecase.NewStorageUsecase(mockRepo, log, validate)

	ctx := context.Background()
	key := "uploads/test_1234567890.jpg"

	// Mock expectation
	mockRepo.On("Delete", ctx, key).Return(nil)

	// Execute
	err := u.DeleteFile(ctx, key)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_DeleteFile_Error(t *testing.T) {
	// Setup
	mockRepo := mock2.NewMockStorageRepository()
	log := logrus.New()
	validate := validator.New()
	u := usecase.NewStorageUsecase(mockRepo, log, validate)

	ctx := context.Background()
	key := "uploads/test_1234567890.jpg"
	expectedError := errors.New("delete failed")

	// Mock expectation
	mockRepo.On("Delete", ctx, key).Return(expectedError)

	// Execute
	err := u.DeleteFile(ctx, key)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_GetPresignedURL_Success(t *testing.T) {
	// Setup
	mockRepo := mock2.NewMockStorageRepository()
	log := logrus.New()
	validate := validator.New()
	u := usecase.NewStorageUsecase(mockRepo, log, validate)

	ctx := context.Background()
	key := "uploads/test_1234567890.jpg"
	expiration := 3600
	expectedURL := "https://presigned-url.com/test"

	// Mock expectation
	mockRepo.On("GetPresignedURL", ctx, key, expiration).Return(expectedURL, nil)

	// Execute
	url, err := u.GetPresignedURL(ctx, key, expiration)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_DownloadFile_Success(t *testing.T) {
	// Setup
	mockRepo := mock2.NewMockStorageRepository()
	log := logrus.New()
	validate := validator.New()
	u := usecase.NewStorageUsecase(mockRepo, log, validate)

	ctx := context.Background()
	key := "uploads/test_1234567890.jpg"
	mockReader := io.NopCloser(strings.NewReader("file content"))

	// Mock expectation
	mockRepo.On("Download", ctx, key).Return(mockReader, nil)

	// Execute
	reader, err := u.DownloadFile(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	mockRepo.AssertExpectations(t)
}
