package usecase

import (
	"context"
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (m *UserUsecaseMock) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*model.UserResponse), args.String(1), args.Error(2)
}

func (m *UserUsecaseMock) Current(ctx context.Context, userID string) (*model.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *UserUsecaseMock) UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *UserUsecaseMock) Logout(ctx context.Context, tokenString string) error {
	args := m.Called(ctx, tokenString)
	return args.Error(0)
}
