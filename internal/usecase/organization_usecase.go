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
	log      *logrus.Logger
	validate *validator.Validate
}

func NewOrganizationRequest(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) OrganizationUsecase {
	return &organizationUsecase{
		db:       db,
		repo:     &repository.Repository[entity.OrganizationMember]{DB: db},
		log:      log,
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

	query := u.db.Where("entity_type = ?", entityType).Where("is_active = ?", true).Order("order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToOrganizationResponses(items), nil
}

func (u *organizationUsecase) GetByID(id string) (*model.OrganizationResponse, error) {
	var g entity.OrganizationMember
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return nil, err
	}
	r := converter.ToOrganizationResponse(&g)
	return &r, nil
}

func (u *organizationUsecase) Create(entityType string, req model.CreateOrganizationRequest) (*model.OrganizationResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	g := entity.OrganizationMember{
		ID:            uuid.New().String(),
		EntityType:    entityType,
		Name:          req.Name,
		Position:      req.Position,
		PositionOrder: req.PositionOrder,
		OrderIndex:    req.OrderIndex,
		IsActive:      req.IsActive,
	}
	if err := u.repo.Create(u.db, &g); err != nil {
		return nil, err
	}
	r := converter.ToOrganizationResponse(&g)
	return &r, nil
}

func (u *organizationUsecase) Update(id string, req model.UpdateOrganizationRequest) (*model.OrganizationResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var g entity.OrganizationMember
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return nil, err
	}
	g.Name = req.Name
	g.Position = req.Position
	g.PositionOrder = req.PositionOrder
	g.OrderIndex = req.OrderIndex
	g.IsActive = req.IsActive
	if err := u.repo.Update(u.db, &g); err != nil {
		return nil, err
	}
	r := converter.ToOrganizationResponse(&g)
	return &r, nil
}

func (u *organizationUsecase) Delete(id string) error {
	var g entity.OrganizationMember
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &g)
}
