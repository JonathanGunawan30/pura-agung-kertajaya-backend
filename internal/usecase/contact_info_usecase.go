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

type ContactInfoUsecase interface {
	GetAll(entityType string) ([]model.ContactInfoResponse, error)
	GetByID(id string) (*model.ContactInfoResponse, error)
	Create(req model.CreateContactInfoRequest) (*model.ContactInfoResponse, error)
	Update(id string, req model.UpdateContactInfoRequest) (*model.ContactInfoResponse, error)
	Delete(id string) error
}

type contactInfoUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.ContactInfo]
	validate *validator.Validate
}

func NewContactInfoUsecase(db *gorm.DB, validate *validator.Validate) ContactInfoUsecase {
	return &contactInfoUsecase{
		db:       db,
		repo:     &repository.Repository[entity.ContactInfo]{DB: db},
		validate: validate,
	}
}

func (u *contactInfoUsecase) GetAll(entityType string) ([]model.ContactInfoResponse, error) {
	var items []entity.ContactInfo
	if entityType == "" {
		entityType = "pura"
	}
	if err := u.db.Where("entity_type = ?", entityType).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return converter.ToContactInfoResponses(items), nil
}

func (u *contactInfoUsecase) GetByID(id string) (*model.ContactInfoResponse, error) {
	var e entity.ContactInfo
	if err := u.repo.FindById(u.db, &e, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("contact info not found")
		}
		return nil, err
	}
	r := converter.ToContactInfoResponse(e)
	return &r, nil
}

func (u *contactInfoUsecase) Create(req model.CreateContactInfoRequest) (*model.ContactInfoResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	e := entity.ContactInfo{
		ID:            uuid.New().String(),
		EntityType:    req.EntityType,
		Address:       req.Address,
		Phone:         req.Phone,
		Email:         req.Email,
		VisitingHours: req.VisitingHours,
		MapEmbedURL:   req.MapEmbedURL,
	}

	if err := u.repo.Create(u.db, &e); err != nil {
		return nil, err
	}

	r := converter.ToContactInfoResponse(e)
	return &r, nil
}

func (u *contactInfoUsecase) Update(id string, req model.UpdateContactInfoRequest) (*model.ContactInfoResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var e entity.ContactInfo
	if err := u.repo.FindById(u.db, &e, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("contact info not found")
		}
		return nil, err
	}

	e.Address = req.Address
	e.Phone = req.Phone
	e.Email = req.Email
	e.VisitingHours = req.VisitingHours
	e.MapEmbedURL = req.MapEmbedURL

	if err := u.repo.Update(u.db, &e); err != nil {
		return nil, err
	}

	r := converter.ToContactInfoResponse(e)
	return &r, nil
}

func (u *contactInfoUsecase) Delete(id string) error {
	var e entity.ContactInfo
	if err := u.repo.FindById(u.db, &e, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound("contact info not found")
		}
		return err
	}
	return u.repo.Delete(u.db, &e)
}
