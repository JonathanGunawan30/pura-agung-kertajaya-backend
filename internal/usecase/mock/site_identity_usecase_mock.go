package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type SiteIdentityUsecaseMock struct{ mock.Mock }

func (m *SiteIdentityUsecaseMock) GetAll(entityType string) ([]model.SiteIdentityResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SiteIdentityResponse), args.Error(1)
}

func (m *SiteIdentityUsecaseMock) GetPublic(entityType string) (*model.SiteIdentityResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SiteIdentityResponse), args.Error(1)
}

func (m *SiteIdentityUsecaseMock) GetByID(id string) (*model.SiteIdentityResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SiteIdentityResponse), args.Error(1)
}

func (m *SiteIdentityUsecaseMock) Create(entityType string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error) {
	args := m.Called(entityType, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SiteIdentityResponse), args.Error(1)
}

func (m *SiteIdentityUsecaseMock) Update(id string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SiteIdentityResponse), args.Error(1)
}

func (m *SiteIdentityUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
