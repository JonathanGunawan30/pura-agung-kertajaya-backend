package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type HeroSlideUsecaseMock struct {
	mock.Mock
}

func (m *HeroSlideUsecaseMock) GetAll(entityType string) ([]model.HeroSlideResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.HeroSlideResponse), args.Error(1)
}

func (m *HeroSlideUsecaseMock) GetPublic(entityType string) ([]model.HeroSlideResponse, error) {
	args := m.Called(entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.HeroSlideResponse), args.Error(1)
}

func (m *HeroSlideUsecaseMock) GetByID(id string) (*model.HeroSlideResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.HeroSlideResponse), args.Error(1)
}

func (m *HeroSlideUsecaseMock) Create(req model.HeroSlideRequest) (*model.HeroSlideResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.HeroSlideResponse), args.Error(1)
}

func (m *HeroSlideUsecaseMock) Update(id string, req model.HeroSlideRequest) (*model.HeroSlideResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.HeroSlideResponse), args.Error(1)
}

func (m *HeroSlideUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
