package test

import (
	"bytes"
	"context"
	"errors"
	"io"
	mock2 "pura-agung-kertajaya-backend/internal/repository/mock"
	"testing"

	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var minimalPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
	0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

func TestStorageUsecase_UploadFile_Success(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	filename := "test.png"
	file := bytes.NewReader(minimalPNG)
	contentType := "image/png"
	fileSize := int64(len(minimalPNG))

	mockRepo.On("Upload", ctx, mock.AnythingOfType("string"), mock.Anything, "image/webp", mock.AnythingOfType("int64")).
		Return("uploads/test_variant.webp", nil)

	result, err := u.UploadFile(ctx, filename, file, contentType, fileSize)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_UploadFile_InvalidImage(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	filename := "test.png"
	file := bytes.NewReader([]byte("not an image"))
	contentType := "image/png"
	fileSize := int64(12)

	result, err := u.UploadFile(ctx, filename, file, contentType, fileSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "image processing failed")

	mockRepo.AssertNotCalled(t, "Upload")
}

func TestStorageUsecase_UploadFile_UploadError_RollbackTriggered(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	filename := "test.png"
	file := bytes.NewReader(minimalPNG)
	contentType := "image/png"
	fileSize := int64(len(minimalPNG))

	mockRepo.On("Upload", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("s3 connection error"))

	mockRepo.On("Delete", mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := u.UploadFile(ctx, filename, file, contentType, fileSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "image processing failed")

	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_DeleteFile_Success(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	key := "uploads/test.jpg"

	mockRepo.On("Delete", ctx, key).Return(nil)

	err := u.DeleteFile(ctx, key)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_DeleteFile_Error(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	key := "uploads/test.jpg"
	expectedError := errors.New("delete failed")

	mockRepo.On("Delete", ctx, key).Return(expectedError)

	err := u.DeleteFile(ctx, key)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_GetPresignedURL_Success(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	key := "uploads/test.jpg"
	expiration := 3600
	expectedURL := "https://presigned-url.com/test"

	mockRepo.On("GetPresignedURL", ctx, key, expiration).Return(expectedURL, nil)

	url, err := u.GetPresignedURL(ctx, key, expiration)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
	mockRepo.AssertExpectations(t)
}

func TestStorageUsecase_DownloadFile_Success(t *testing.T) {
	mockRepo := mock2.NewMockStorageRepository()
	u := usecase.NewStorageUsecase(mockRepo)

	ctx := context.Background()
	key := "uploads/test.jpg"
	mockReader := io.NopCloser(bytes.NewReader([]byte("file content")))

	mockRepo.On("Download", ctx, key).Return(mockReader, nil)

	reader, err := u.DownloadFile(ctx, key)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	mockRepo.AssertExpectations(t)
}
