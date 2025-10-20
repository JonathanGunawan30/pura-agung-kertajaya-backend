package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type TestimonialUsecaseMock struct {
	mock.Mock
}

func (m *TestimonialUsecaseMock) GetAll() ([]model.TestimonialResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TestimonialResponse), args.Error(1)
}

func (m *TestimonialUsecaseMock) GetByID(id int) (*model.TestimonialResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TestimonialResponse), args.Error(1)
}

func (m *TestimonialUsecaseMock) Create(req model.TestimonialRequest) (*model.TestimonialResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TestimonialResponse), args.Error(1)
}

func (m *TestimonialUsecaseMock) Update(id int, req model.TestimonialRequest) (*model.TestimonialResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TestimonialResponse), args.Error(1)
}

func (m *TestimonialUsecaseMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
