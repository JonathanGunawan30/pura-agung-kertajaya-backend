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

type SiteIdentityUsecase interface {
	GetAll(entityType string) ([]model.SiteIdentityResponse, error)
	GetPublic(entityType string) (*model.SiteIdentityResponse, error)
	GetByID(id string) (*model.SiteIdentityResponse, error)
	Create(entityType string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error)
	Update(id string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error)
	Delete(id string) error
}

type siteIdentityUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.SiteIdentity]
	validate *validator.Validate
}

func NewSiteIdentityUsecase(db *gorm.DB, validate *validator.Validate) SiteIdentityUsecase {
	return &siteIdentityUsecase{
		db:       db,
		repo:     &repository.Repository[entity.SiteIdentity]{DB: db},
		validate: validate,
	}
}

func (u *siteIdentityUsecase) GetAll(entityType string) ([]model.SiteIdentityResponse, error) {
	var items []entity.SiteIdentity
	query := u.db.Order("created_at ASC")
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}
	resp := make([]model.SiteIdentityResponse, 0, len(items))
	for _, e := range items {
		resp = append(resp, converter.ToSiteIdentityResponse(e))
	}
	return resp, nil
}

func (u *siteIdentityUsecase) GetPublic(entityType string) (*model.SiteIdentityResponse, error) {
	var e entity.SiteIdentity
	query := u.db.Order("created_at DESC")

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	if err := query.Limit(1).Take(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("site identity not found")
		}
		return nil, err
	}

	r := converter.ToSiteIdentityResponse(e)
	return &r, nil
}

func (u *siteIdentityUsecase) GetByID(id string) (*model.SiteIdentityResponse, error) {
	var e entity.SiteIdentity
	if err := u.repo.FindById(u.db, &e, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("site identity not found")
		}
		return nil, err
	}
	r := converter.ToSiteIdentityResponse(e)
	return &r, nil
}

func (u *siteIdentityUsecase) Create(entityType string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	e := entity.SiteIdentity{
		ID:                  uuid.New().String(),
		EntityType:          entityType,
		SiteName:            req.SiteName,
		LogoURL:             req.LogoURL,
		Tagline:             req.Tagline,
		PrimaryButtonText:   req.PrimaryButtonText,
		PrimaryButtonLink:   req.PrimaryButtonLink,
		SecondaryButtonText: req.SecondaryButtonText,
		SecondaryButtonLink: req.SecondaryButtonLink,
	}
	if err := u.repo.Create(u.db, &e); err != nil {
		return nil, err
	}
	r := converter.ToSiteIdentityResponse(e)
	return &r, nil
}

func (u *siteIdentityUsecase) Update(id string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var e entity.SiteIdentity
	if err := u.repo.FindById(u.db, &e, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("site identity not found")
		}
		return nil, err
	}

	e.EntityType = req.EntityType
	e.SiteName = req.SiteName
	e.LogoURL = req.LogoURL
	e.Tagline = req.Tagline
	e.PrimaryButtonText = req.PrimaryButtonText
	e.PrimaryButtonLink = req.PrimaryButtonLink
	e.SecondaryButtonText = req.SecondaryButtonText
	e.SecondaryButtonLink = req.SecondaryButtonLink

	if err := u.repo.Update(u.db, &e); err != nil {
		return nil, err
	}
	r := converter.ToSiteIdentityResponse(e)
	return &r, nil
}

func (u *siteIdentityUsecase) Delete(id string) error {
	var e entity.SiteIdentity
	if err := u.repo.FindById(u.db, &e, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound("site identity not found")
		}
		return err
	}
	return u.repo.Delete(u.db, &e)
}
