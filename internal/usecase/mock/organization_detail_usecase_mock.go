package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type OrganizationDetailUsecaseMock struct {
	mock.Mock
}

func (m *OrganizationDetailUsecaseMock) GetByEntityType(entityType string) (*model.OrganizationDetailResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationDetailResponse), args.Error(1)
}

func (m *OrganizationDetailUsecaseMock) Update(entityType string, req model.UpdateOrganizationDetailRequest) (*model.OrganizationDetailResponse, error) {
	args := m.Called(entityType, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrganizationDetailResponse), args.Error(1)
}
