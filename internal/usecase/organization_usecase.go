package usecase

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationUsecase interface {
	GetAll(entityType string) ([]model.OrganizationResponse, error)
	GetPublic(entityType string) ([]model.OrganizationResponse, error)
	GetByID(id string) (*model.OrganizationResponse, error)
	Create(entityType string, req model.CreateOrganizationRequest) (*model.OrganizationResponse, error)
	Update(id string, req model.UpdateOrganizationRequest) (*model.OrganizationResponse, error)
	Delete(id string) error
}

type organizationUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.OrganizationMember]
	validate *validator.Validate
}

func NewOrganizationUsecase(db *gorm.DB, validate *validator.Validate) OrganizationUsecase {
	return &organizationUsecase{
		db:       db,
		repo:     &repository.Repository[entity.OrganizationMember]{DB: db},
		validate: validate,
	}
}

func (u *organizationUsecase) GetAll(entityType string) ([]model.OrganizationResponse, error) {
	var items []entity.OrganizationMember

	query := u.db.Where("entity_type = ?", entityType).Order("position_order ASC, order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToOrganizationResponses(items), nil
}

func (u *organizationUsecase) GetPublic(entityType string) ([]model.OrganizationResponse, error) {
	var items []entity.OrganizationMember
	query := u.db.Where("entity_type = ?", entityType).
		Where("is_active = ?", true).
		Order("position_order ASC, order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToOrganizationResponses(items), nil
}

func (u *organizationUsecase) GetByID(id string) (*model.OrganizationResponse, error) {
	var m entity.OrganizationMember
	if err := u.repo.FindById(u.db, &m, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("organization member not found")
		}
		return nil, err
	}
	r := converter.ToOrganizationResponse(&m)
	return &r, nil
}

func (u *organizationUsecase) Create(entityType string, req model.CreateOrganizationRequest) (*model.OrganizationResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	m := entity.OrganizationMember{
		ID:            uuid.New().String(),
		EntityType:    entityType,
		Name:          req.Name,
		Position:      req.Position,
		PositionOrder: req.PositionOrder,
		OrderIndex:    req.OrderIndex,
		IsActive:      req.IsActive,
	}
	if err := u.repo.Create(u.db, &m); err != nil {
		return nil, err
	}
	r := converter.ToOrganizationResponse(&m)
	return &r, nil
}

func (u *organizationUsecase) Update(id string, req model.UpdateOrganizationRequest) (*model.OrganizationResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var m entity.OrganizationMember
	if err := u.repo.FindById(u.db, &m, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("organization member not found")
		}
		return nil, err
	}

	m.Name = req.Name
	m.Position = req.Position
	m.PositionOrder = req.PositionOrder
	m.OrderIndex = req.OrderIndex
	m.IsActive = req.IsActive

	if err := u.repo.Update(u.db, &m); err != nil {
		return nil, err
	}
	r := converter.ToOrganizationResponse(&m)
	return &r, nil
}

func (u *organizationUsecase) Delete(id string) error {
	var m entity.OrganizationMember
	if err := u.repo.FindById(u.db, &m, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound("organization member not found")
		}
		return err
	}
	return u.repo.Delete(u.db, &m)
}
