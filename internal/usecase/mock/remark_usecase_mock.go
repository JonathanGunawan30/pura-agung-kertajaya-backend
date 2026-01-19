package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type RemarkUsecaseMock struct {
	mock.Mock
}

func (m *RemarkUsecaseMock) GetAll(entityType string) ([]model.RemarkResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RemarkResponse), args.Error(1)
}

func (m *RemarkUsecaseMock) GetPublic(entityType string) ([]model.RemarkResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RemarkResponse), args.Error(1)
}

func (m *RemarkUsecaseMock) GetByID(id string) (*model.RemarkResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RemarkResponse), args.Error(1)
}

func (m *RemarkUsecaseMock) Create(entityType string, req model.CreateRemarkRequest) (*model.RemarkResponse, error) {
	args := m.Called(entityType, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RemarkResponse), args.Error(1)
}

func (m *RemarkUsecaseMock) Update(id string, req model.UpdateRemarkRequest) (*model.RemarkResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RemarkResponse), args.Error(1)
}

func (m *RemarkUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
