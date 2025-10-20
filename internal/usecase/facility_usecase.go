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

type FacilityUsecase interface {
	GetAll() ([]model.FacilityResponse, error)
	GetPublic() ([]model.FacilityResponse, error)
	GetByID(id string) (*model.FacilityResponse, error)
	Create(req model.FacilityRequest) (*model.FacilityResponse, error)
	Update(id string, req model.FacilityRequest) (*model.FacilityResponse, error)
	Delete(id string) error
}

type facilityUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Facility]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewFacilityUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) FacilityUsecase {
	return &facilityUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Facility]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *facilityUsecase) GetAll() ([]model.FacilityResponse, error) {
	var items []entity.Facility
	if err := u.db.Order("order_index ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	resp := make([]model.FacilityResponse, 0, len(items))
	for _, g := range items {
		resp = append(resp, converter.ToFacilityResponse(g))
	}
	return resp, nil
}

func (u *facilityUsecase) GetPublic() ([]model.FacilityResponse, error) {
	var items []entity.Facility
	if err := u.db.Where("is_active = ?", true).Order("order_index ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	resp := make([]model.FacilityResponse, 0, len(items))
	for _, g := range items {
		resp = append(resp, converter.ToFacilityResponse(g))
	}
	return resp, nil
}

func (u *facilityUsecase) GetByID(id string) (*model.FacilityResponse, error) {
	var g entity.Facility
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return nil, err
	}
	r := converter.ToFacilityResponse(g)
	return &r, nil
}

func (u *facilityUsecase) Create(req model.FacilityRequest) (*model.FacilityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	g := entity.Facility{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
	}
	if err := u.repo.Create(u.db, &g); err != nil {
		return nil, err
	}
	r := converter.ToFacilityResponse(g)
	return &r, nil
}

func (u *facilityUsecase) Update(id string, req model.FacilityRequest) (*model.FacilityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var g entity.Facility
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return nil, err
	}
	g.Name = req.Name
	g.Description = req.Description
	g.ImageURL = req.ImageURL
	g.OrderIndex = req.OrderIndex
	g.IsActive = req.IsActive
	if err := u.repo.Update(u.db, &g); err != nil {
		return nil, err
	}
	r := converter.ToFacilityResponse(g)
	return &r, nil
}

func (u *facilityUsecase) Delete(id string) error {
	var g entity.Facility
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &g)
}
