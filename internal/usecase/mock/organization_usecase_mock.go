package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type OrganizationMemberUsecaseMock struct {
	mock.Mock
}

// GetAll mocks base method.
func (m *OrganizationMemberUsecaseMock) GetAll() ([]model.OrganizationResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.OrganizationResponse), args.Error(1)
}

// GetPublic mocks base method.
func (m *OrganizationMemberUsecaseMock) GetPublic() ([]model.OrganizationResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.OrganizationResponse), args.Error(1)
}

// GetByID mocks base method.
func (m *OrganizationMemberUsecaseMock) GetByID(id string) (*model.OrganizationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationResponse), args.Error(1)
}

// Create mocks base method.
func (m *OrganizationMemberUsecaseMock) Create(req model.OrganizationRequest) (*model.OrganizationResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationResponse), args.Error(1)
}

// Update mocks base method.
func (m *OrganizationMemberUsecaseMock) Update(id string, req model.OrganizationRequest) (*model.OrganizationResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationResponse), args.Error(1)
}

// Delete mocks base method.
func (m *OrganizationMemberUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
