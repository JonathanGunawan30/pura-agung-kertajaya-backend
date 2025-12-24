package usecase

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RemarkUsecase interface {
	GetAll(entityType string) ([]model.RemarkResponse, error)
	GetPublic(entityType string) ([]model.RemarkResponse, error)
	GetByID(id string) (*model.RemarkResponse, error)
	Create(req model.CreateRemarkRequest) (*model.RemarkResponse, error)
	Update(id string, req model.UpdateRemarkRequest) (*model.RemarkResponse, error)
	Delete(id string) error
}

type remarkUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Remark]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewRemarkUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) RemarkUsecase {
	return &remarkUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Remark]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *remarkUsecase) GetAll(entityType string) ([]model.RemarkResponse, error) {
	var remarks []entity.Remark

	if entityType == "" {
		entityType = "pura"
	}

	query := u.db.Where("entity_type = ?", entityType).Order("order_index ASC")

	if err := u.repo.FindAll(query, &remarks); err != nil {
		return nil, err
	}

	return converter.ToRemarkResponses(remarks), nil
}

func (u *remarkUsecase) GetPublic(entityType string) ([]model.RemarkResponse, error) {
	var remarks []entity.Remark
	if entityType == "" {
		entityType = "pura"
	}

	query := u.db.Where("is_active = ? AND entity_type = ?", true, entityType).
		Order("order_index ASC")

	if err := u.repo.FindAll(query, &remarks); err != nil {
		return nil, err
	}
	return converter.ToRemarkResponses(remarks), nil
}

func (u *remarkUsecase) GetByID(id string) (*model.RemarkResponse, error) {
	var r entity.Remark
	if err := u.repo.FindById(u.db, &r, id); err != nil {
		return nil, err
	}

	response := converter.ToRemarkResponse(&r)
	return &response, nil
}

func (u *remarkUsecase) Create(req model.CreateRemarkRequest) (*model.RemarkResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	r := entity.Remark{
		ID:         uuid.New().String(),
		EntityType: req.EntityType,
		Name:       req.Name,
		Position:   req.Position,
		ImageURL:   req.ImageURL,
		Content:    req.Content,
		IsActive:   req.IsActive,
		OrderIndex: req.OrderIndex,
	}

	if r.OrderIndex == 0 {
		r.OrderIndex = 1
	}

	if err := u.repo.Create(u.db, &r); err != nil {
		return nil, err
	}

	response := converter.ToRemarkResponse(&r)
	return &response, nil
}

func (u *remarkUsecase) Update(id string, req model.UpdateRemarkRequest) (*model.RemarkResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var r entity.Remark
	if err := u.repo.FindById(u.db, &r, id); err != nil {
		return nil, err
	}

	r.Name = req.Name
	r.Position = req.Position
	r.ImageURL = req.ImageURL
	r.Content = req.Content
	r.OrderIndex = req.OrderIndex
	r.IsActive = req.IsActive

	if err := u.repo.Update(u.db, &r); err != nil {
		return nil, err
	}

	response := converter.ToRemarkResponse(&r)
	return &response, nil
}

func (u *remarkUsecase) Delete(id string) error {
	var r entity.Remark
	if err := u.repo.FindById(u.db, &r, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &r)
}
