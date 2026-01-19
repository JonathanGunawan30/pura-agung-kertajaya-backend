package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type GalleryUsecaseMock struct {
	mock.Mock
}

func (m *GalleryUsecaseMock) GetAll(entityType string) ([]model.GalleryResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.GalleryResponse), args.Error(1)
}

func (m *GalleryUsecaseMock) GetPublic(entityType string) ([]model.GalleryResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.GalleryResponse), args.Error(1)
}

func (m *GalleryUsecaseMock) GetByID(id string) (*model.GalleryResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GalleryResponse), args.Error(1)
}

func (m *GalleryUsecaseMock) Create(entityType string, req model.CreateGalleryRequest) (*model.GalleryResponse, error) {
	args := m.Called(entityType, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GalleryResponse), args.Error(1)
}

func (m *GalleryUsecaseMock) Update(id string, req model.UpdateGalleryRequest) (*model.GalleryResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GalleryResponse), args.Error(1)
}

func (m *GalleryUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
