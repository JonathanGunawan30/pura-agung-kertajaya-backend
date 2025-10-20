package usecase

import (
    "pura-agung-kertajaya-backend/internal/model"

    "github.com/stretchr/testify/mock"
)

type AboutUsecaseMock struct{ mock.Mock }

func (m *AboutUsecaseMock) GetAll() ([]model.AboutSectionResponse, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]model.AboutSectionResponse), args.Error(1)
}

func (m *AboutUsecaseMock) GetPublic() ([]model.AboutSectionResponse, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]model.AboutSectionResponse), args.Error(1)
}

func (m *AboutUsecaseMock) GetByID(id string) (*model.AboutSectionResponse, error) {
    args := m.Called(id)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*model.AboutSectionResponse), args.Error(1)
}

func (m *AboutUsecaseMock) Create(req model.AboutSectionRequest) (*model.AboutSectionResponse, error) {
    args := m.Called(req)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*model.AboutSectionResponse), args.Error(1)
}

func (m *AboutUsecaseMock) Update(id string, req model.AboutSectionRequest) (*model.AboutSectionResponse, error) {
    args := m.Called(id, req)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*model.AboutSectionResponse), args.Error(1)
}

func (m *AboutUsecaseMock) Delete(id string) error {
    args := m.Called(id)
    return args.Error(0)
}
