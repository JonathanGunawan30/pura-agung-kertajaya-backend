package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type FacilityUsecaseMock struct{ mock.Mock }

func (m *FacilityUsecaseMock) GetAll(entityType string) ([]model.FacilityResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.FacilityResponse), args.Error(1)
}

func (m *FacilityUsecaseMock) GetPublic(entityType string) ([]model.FacilityResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.FacilityResponse), args.Error(1)
}

func (m *FacilityUsecaseMock) GetByID(id string) (*model.FacilityResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FacilityResponse), args.Error(1)
}

func (m *FacilityUsecaseMock) Create(req model.CreateFacilityRequest) (*model.FacilityResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FacilityResponse), args.Error(1)
}

func (m *FacilityUsecaseMock) Update(id string, req model.UpdateFacilityRequest) (*model.FacilityResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FacilityResponse), args.Error(1)
}

func (m *FacilityUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
