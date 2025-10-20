package usecase

import (
    "pura-agung-kertajaya-backend/internal/model"

    "github.com/stretchr/testify/mock"
)

type ActivityUsecaseMock struct{ mock.Mock }

func (m *ActivityUsecaseMock) GetAll() ([]model.ActivityResponse, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]model.ActivityResponse), args.Error(1)
}

func (m *ActivityUsecaseMock) GetPublic() ([]model.ActivityResponse, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]model.ActivityResponse), args.Error(1)
}

func (m *ActivityUsecaseMock) GetByID(id string) (*model.ActivityResponse, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.ActivityResponse), args.Error(1)
}

func (m *ActivityUsecaseMock) Create(req model.ActivityRequest) (*model.ActivityResponse, error) {
    args := m.Called(req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.ActivityResponse), args.Error(1)
}

func (m *ActivityUsecaseMock) Update(id string, req model.ActivityRequest) (*model.ActivityResponse, error) {
    args := m.Called(id, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.ActivityResponse), args.Error(1)
}

func (m *ActivityUsecaseMock) Delete(id string) error {
    args := m.Called(id)
    return args.Error(0)
}
