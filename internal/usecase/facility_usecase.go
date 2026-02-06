package usecase

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FacilityUsecase interface {
	GetAll(entityType string) ([]model.FacilityResponse, error)
	GetPublic(entityType string) ([]model.FacilityResponse, error)
	GetByID(id string) (*model.FacilityResponse, error)
	Create(entityType string, req model.CreateFacilityRequest) (*model.FacilityResponse, error)
	Update(id string, req model.UpdateFacilityRequest) (*model.FacilityResponse, error)
	Delete(id string) error
}

type facilityUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Facility]
	validate *validator.Validate
}

func NewFacilityUsecase(db *gorm.DB, validate *validator.Validate) FacilityUsecase {
	return &facilityUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Facility]{DB: db},
		validate: validate,
	}
}

func (u *facilityUsecase) GetAll(entityType string) ([]model.FacilityResponse, error) {
	var items []entity.Facility

	query := u.db.Where("entity_type = ?", entityType).Order("order_index ASC")
	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToFacilityResponses(items), nil
}

func (u *facilityUsecase) GetPublic(entityType string) ([]model.FacilityResponse, error) {
	var items []entity.Facility

	query := u.db.Where("entity_type = ?", entityType).Where("is_active = ?", true).Order("order_index ASC")
	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToFacilityResponses(items), nil
}

func (u *facilityUsecase) GetByID(id string) (*model.FacilityResponse, error) {
	var f entity.Facility
	if err := u.repo.FindById(u.db, &f, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("facility not found")
		}
		return nil, err
	}
	r := converter.ToFacilityResponse(f)
	return &r, nil
}

func (u *facilityUsecase) Create(entityType string, req model.CreateFacilityRequest) (*model.FacilityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	f := entity.Facility{
		ID:          uuid.New().String(),
		EntityType:  entityType,
		Name:        req.Name,
		Description: req.Description,
		Images:      util.ImageMap(req.Images),
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
	}
	if err := u.repo.Create(u.db, &f); err != nil {
		return nil, err
	}
	r := converter.ToFacilityResponse(f)
	return &r, nil
}

func (u *facilityUsecase) Update(id string, req model.UpdateFacilityRequest) (*model.FacilityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var f entity.Facility
	if err := u.repo.FindById(u.db, &f, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("facility not found")
		}
		return nil, err
	}

	f.Name = req.Name
	f.Description = req.Description
	f.Images = util.ImageMap(req.Images)
	f.OrderIndex = req.OrderIndex
	f.IsActive = req.IsActive

	if err := u.repo.Update(u.db, &f); err != nil {
		return nil, err
	}
	r := converter.ToFacilityResponse(f)
	return &r, nil
}

func (u *facilityUsecase) Delete(id string) error {
	var f entity.Facility
	if err := u.repo.FindById(u.db, &f, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound("facility not found")
		}
		return err
	}
	return u.repo.Delete(u.db, &f)
}
