package usecase

import (
	"pura-agung-kertajaya-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

type ArticleUsecaseMock struct {
	mock.Mock
}

func (m *ArticleUsecaseMock) GetAll(filter string) ([]model.ArticleResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ArticleResponse), args.Error(1)
}

func (m *ArticleUsecaseMock) GetPublic(limit int) ([]model.ArticleResponse, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ArticleResponse), args.Error(1)
}

func (m *ArticleUsecaseMock) GetByID(id string) (*model.ArticleResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ArticleResponse), args.Error(1)
}

func (m *ArticleUsecaseMock) GetBySlug(slug string) (*model.ArticleResponse, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ArticleResponse), args.Error(1)
}

func (m *ArticleUsecaseMock) Create(req model.CreateArticleRequest) (*model.ArticleResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ArticleResponse), args.Error(1)
}

func (m *ArticleUsecaseMock) Update(id string, req model.UpdateArticleRequest) (*model.ArticleResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ArticleResponse), args.Error(1)
}

func (m *ArticleUsecaseMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
