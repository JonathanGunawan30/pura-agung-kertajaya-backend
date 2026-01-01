package usecase

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GalleryUsecase interface {
	GetAll(entityType string) ([]model.GalleryResponse, error)
	GetPublic(entityType string) ([]model.GalleryResponse, error)
	GetByID(id string) (*model.GalleryResponse, error)
	Create(req model.CreateGalleryRequest) (*model.GalleryResponse, error)
	Update(id string, req model.UpdateGalleryRequest) (*model.GalleryResponse, error)
	Delete(id string) error
}

type galleryUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Gallery]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewGalleryUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) GalleryUsecase {
	return &galleryUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Gallery]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *galleryUsecase) GetAll(entityType string) ([]model.GalleryResponse, error) {
	var items []entity.Gallery

	if entityType == "" {
		entityType = "pura"
	}

	query := u.db.Where("entity_type = ?", entityType).Order("order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToGalleryResponses(items), nil
}

func (u *galleryUsecase) GetPublic(entityType string) ([]model.GalleryResponse, error) {
	var items []entity.Gallery

	if entityType == "" {
		entityType = "pura"
	}

	query := u.db.Where("entity_type = ?", entityType).Where("is_active = ?", true).Order("order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToGalleryResponses(items), nil
}

func (u *galleryUsecase) GetByID(id string) (*model.GalleryResponse, error) {
	var g entity.Gallery
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return nil, err
	}
	r := converter.ToGalleryResponse(&g)
	return &r, nil
}

func (u *galleryUsecase) Create(req model.CreateGalleryRequest) (*model.GalleryResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	g := entity.Gallery{
		ID:          uuid.New().String(),
		EntityType:  req.EntityType,
		Title:       req.Title,
		Description: req.Description,
		Images:      util.ImageMap(req.Images),
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
	}
	if err := u.repo.Create(u.db, &g); err != nil {
		return nil, err
	}
	r := converter.ToGalleryResponse(&g)
	return &r, nil
}

func (u *galleryUsecase) Update(id string, req model.UpdateGalleryRequest) (*model.GalleryResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var g entity.Gallery
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return nil, err
	}
	g.Title = req.Title
	g.Description = req.Description
	g.Images = util.ImageMap(req.Images)
	g.OrderIndex = req.OrderIndex
	g.IsActive = req.IsActive
	if err := u.repo.Update(u.db, &g); err != nil {
		return nil, err
	}
	r := converter.ToGalleryResponse(&g)
	return &r, nil
}

func (u *galleryUsecase) Delete(id string) error {
	var g entity.Gallery
	if err := u.repo.FindById(u.db, &g, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &g)
}
