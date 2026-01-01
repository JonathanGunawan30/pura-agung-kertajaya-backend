package usecase

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type MockStorageUsecase struct {
	mock.Mock
}

func NewMockStorageUsecase() *MockStorageUsecase {
	return &MockStorageUsecase{}
}

func (m *MockStorageUsecase) UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (map[string]string, error) {
	args := m.Called(ctx, filename, file, contentType, fileSize)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockStorageUsecase) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageUsecase) DeleteFile(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorageUsecase) GetPresignedURL(ctx context.Context, key string, expiration int) (string, error) {
	args := m.Called(ctx, key, expiration)
	return args.String(0), args.Error(1)
}
