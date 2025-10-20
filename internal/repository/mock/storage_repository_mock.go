package mock

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type MockStorageRepository struct {
	mock.Mock
}

func NewMockStorageRepository() *MockStorageRepository {
	return &MockStorageRepository{}
}

func (m *MockStorageRepository) Upload(ctx context.Context, key string, file io.Reader, contentType string, fileSize int64) (string, error) {
	args := m.Called(ctx, key, file, contentType, fileSize)
	return args.String(0), args.Error(1)
}

func (m *MockStorageRepository) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorageRepository) GetPresignedURL(ctx context.Context, key string, expiration int) (string, error) {
	args := m.Called(ctx, key, expiration)
	return args.String(0), args.Error(1)
}
