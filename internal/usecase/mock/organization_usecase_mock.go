package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type OrganizationMemberUsecaseMock struct {
	mock.Mock
}

func (m *OrganizationMemberUsecaseMock) GetAll(entityType string) ([]model.OrganizationResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.OrganizationResponse), args.Error(1)
}

func (m *OrganizationMemberUsecaseMock) GetPublic(entityType string) ([]model.OrganizationResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.OrganizationResponse), args.Error(1)
}

func (m *OrganizationMemberUsecaseMock) GetByID(id string) (*model.OrganizationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationResponse), args.Error(1)
}

func (m *OrganizationMemberUsecaseMock) Create(entityType string, req model.CreateOrganizationRequest) (*model.OrganizationResponse, error) {
	args := m.Called(entityType, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationResponse), args.Error(1)
}

func (m *OrganizationMemberUsecaseMock) Update(id string, req model.UpdateOrganizationRequest) (*model.OrganizationResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationResponse), args.Error(1)
}

func (m *OrganizationMemberUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
