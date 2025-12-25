package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type ContactInfoUsecaseMock struct{ mock.Mock }

func (m *ContactInfoUsecaseMock) GetAll(entityType string) ([]model.ContactInfoResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ContactInfoResponse), args.Error(1)
}

func (m *ContactInfoUsecaseMock) GetByID(id string) (*model.ContactInfoResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContactInfoResponse), args.Error(1)
}

func (m *ContactInfoUsecaseMock) Create(req model.CreateContactInfoRequest) (*model.ContactInfoResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContactInfoResponse), args.Error(1)
}

func (m *ContactInfoUsecaseMock) Update(id string, req model.UpdateContactInfoRequest) (*model.ContactInfoResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContactInfoResponse), args.Error(1)
}

func (m *ContactInfoUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
