package usecase

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ActivityUsecase interface {
	GetAll(entityType string) ([]model.ActivityResponse, error)
	GetPublic(entityType string) ([]model.ActivityResponse, error)
	GetByID(id string) (*model.ActivityResponse, error)
	Create(req model.CreateActivityRequest) (*model.ActivityResponse, error)
	Update(id string, req model.UpdateActivityRequest) (*model.ActivityResponse, error)
	Delete(id string) error
}

type activityUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Activity]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewActivityUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) ActivityUsecase {
	return &activityUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Activity]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *activityUsecase) GetAll(entityType string) ([]model.ActivityResponse, error) {
	var items []entity.Activity
	if entityType == "" {
		entityType = "pura"
	}
	query := u.db.Where("entity_type = ?", entityType).Order("event_date DESC").Order("order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToActivityResponses(items), nil
}

func (u *activityUsecase) GetPublic(entityType string) ([]model.ActivityResponse, error) {
	var items []entity.Activity
	if entityType == "" {
		entityType = "pura"
	}
	query := u.db.Where("entity_type = ?", entityType).Where("is_active = ?", true).Order("event_date DESC").Order("order_index ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToActivityResponses(items), nil
}

func (u *activityUsecase) GetByID(id string) (*model.ActivityResponse, error) {
	var a entity.Activity
	if err := u.repo.FindById(u.db, &a, id); err != nil {
		return nil, err
	}
	r := converter.ToActivityResponse(&a)
	return &r, nil
}

func (u *activityUsecase) Create(req model.CreateActivityRequest) (*model.ActivityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, errors.New("format event_date is invalid")
	}

	a := entity.Activity{
		ID:          uuid.New().String(),
		EntityType:  req.EntityType,
		Title:       req.Title,
		Description: req.Description,
		TimeInfo:    req.TimeInfo,
		Location:    req.Location,
		EventDate:   eventDate,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
	}
	if err := u.repo.Create(u.db, &a); err != nil {
		return nil, err
	}
	r := converter.ToActivityResponse(&a)
	return &r, nil
}

func (u *activityUsecase) Update(id string, req model.UpdateActivityRequest) (*model.ActivityResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}
	var a entity.Activity
	if err := u.repo.FindById(u.db, &a, id); err != nil {
		return nil, err
	}

	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, errors.New("format event_date is invalid")
	}

	a.Title = req.Title
	a.Description = req.Description
	a.TimeInfo = req.TimeInfo
	a.Location = req.Location
	a.EventDate = eventDate
	a.OrderIndex = req.OrderIndex
	a.IsActive = req.IsActive

	if err := u.repo.Update(u.db, &a); err != nil {
		return nil, err
	}
	r := converter.ToActivityResponse(&a)
	return &r, nil
}

func (u *activityUsecase) Delete(id string) error {
	var a entity.Activity
	if err := u.repo.FindById(u.db, &a, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &a)
}
