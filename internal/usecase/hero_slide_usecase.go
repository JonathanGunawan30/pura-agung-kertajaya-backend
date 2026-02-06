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

type HeroSlideUsecase interface {
	GetAll(entityType string) ([]model.HeroSlideResponse, error)
	GetPublic(entityType string) ([]model.HeroSlideResponse, error)
	GetByID(id string) (*model.HeroSlideResponse, error)
	Create(entityType string, req model.HeroSlideRequest) (*model.HeroSlideResponse, error)
	Update(id string, req model.HeroSlideRequest) (*model.HeroSlideResponse, error)
	Delete(id string) error
}

type heroSlideUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.HeroSlide]
	validate *validator.Validate
}

func NewHeroSlideUsecase(db *gorm.DB, validate *validator.Validate) HeroSlideUsecase {
	return &heroSlideUsecase{
		db:       db,
		repo:     &repository.Repository[entity.HeroSlide]{DB: db},
		validate: validate,
	}
}

func (u *heroSlideUsecase) GetAll(entityType string) ([]model.HeroSlideResponse, error) {
	var slides []entity.HeroSlide
	query := u.db.Order("order_index ASC")
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if err := u.repo.FindAll(query, &slides); err != nil {
		return nil, err
	}
	responses := make([]model.HeroSlideResponse, 0, len(slides))
	for _, s := range slides {
		responses = append(responses, converter.ToHeroSlideResponse(s))
	}
	return responses, nil
}

func (u *heroSlideUsecase) GetPublic(entityType string) ([]model.HeroSlideResponse, error) {
	var slides []entity.HeroSlide
	query := u.db.Where("is_active = ?", true).Order("order_index ASC")
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if err := u.repo.FindAll(query, &slides); err != nil {
		return nil, err
	}
	responses := make([]model.HeroSlideResponse, 0, len(slides))
	for _, s := range slides {
		responses = append(responses, converter.ToHeroSlideResponse(s))
	}
	return responses, nil
}

func (u *heroSlideUsecase) GetByID(id string) (*model.HeroSlideResponse, error) {
	var s entity.HeroSlide
	if err := u.repo.FindById(u.db, &s, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("hero slide not found")
		}
		return nil, err
	}
	resp := converter.ToHeroSlideResponse(s)
	return &resp, nil
}

func (u *heroSlideUsecase) Create(entityType string, req model.HeroSlideRequest) (*model.HeroSlideResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	s := entity.HeroSlide{
		ID:         uuid.New().String(),
		EntityType: entityType,
		Images:     util.ImageMap(req.Images),
		OrderIndex: req.OrderIndex,
		IsActive:   req.IsActive,
	}

	if err := u.repo.Create(u.db, &s); err != nil {
		return nil, err
	}

	resp := converter.ToHeroSlideResponse(s)
	return &resp, nil
}

func (u *heroSlideUsecase) Update(id string, req model.HeroSlideRequest) (*model.HeroSlideResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var s entity.HeroSlide
	if err := u.repo.FindById(u.db, &s, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("hero slide not found")
		}
		return nil, err
	}

	s.EntityType = req.EntityType
	s.Images = util.ImageMap(req.Images)
	s.OrderIndex = req.OrderIndex
	s.IsActive = req.IsActive

	if err := u.repo.Update(u.db, &s); err != nil {
		return nil, err
	}

	resp := converter.ToHeroSlideResponse(s)
	return &resp, nil
}

func (u *heroSlideUsecase) Delete(id string) error {
	var s entity.HeroSlide
	if err := u.repo.FindById(u.db, &s, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound("hero slide not found")
		}
		return err
	}
	return u.repo.Delete(u.db, &s)
}
