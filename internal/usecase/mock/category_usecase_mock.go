package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type CategoryUsecaseMock struct {
	mock.Mock
}

func (m *CategoryUsecaseMock) GetAll() ([]model.CategoryResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CategoryResponse), args.Error(1)
}

func (m *CategoryUsecaseMock) GetByID(id string) (*model.CategoryResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CategoryResponse), args.Error(1)
}

func (m *CategoryUsecaseMock) Create(req model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CategoryResponse), args.Error(1)
}

func (m *CategoryUsecaseMock) Update(id string, req model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CategoryResponse), args.Error(1)
}

func (m *CategoryUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
