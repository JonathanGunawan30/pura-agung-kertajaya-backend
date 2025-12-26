package usecase

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationDetailUsecase interface {
	GetByEntityType(entityType string) (*model.OrganizationDetailResponse, error)
	Update(entityType string, req model.UpdateOrganizationDetailRequest) (*model.OrganizationDetailResponse, error)
}

type organizationDetailUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.OrganizationDetail]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewOrganizationDetailUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) OrganizationDetailUsecase {
	return &organizationDetailUsecase{
		db:       db,
		repo:     &repository.Repository[entity.OrganizationDetail]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *organizationDetailUsecase) GetByEntityType(entityType string) (*model.OrganizationDetailResponse, error) {
	var detail entity.OrganizationDetail

	if entityType == "" {
		entityType = "pura"
	}

	query := u.db.Where("entity_type = ?", entityType)

	if err := query.First(&detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	response := converter.ToOrganizationDetailResponse(&detail)
	return &response, nil
}

func (u *organizationDetailUsecase) Update(entityType string, req model.UpdateOrganizationDetailRequest) (*model.OrganizationDetailResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	if entityType == "" {
		entityType = "pura"
	}

	var detail entity.OrganizationDetail

	err := u.db.Where("entity_type = ?", entityType).First(&detail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		detail = entity.OrganizationDetail{
			ID:          uuid.New().String(),
			EntityType:  entityType,
			Vision:      req.Vision,
			Mission:     req.Mission,
			Rules:       req.Rules,
			WorkProgram: req.WorkProgram,
			ImageURL:    req.ImageURL,
		}

		if err := u.repo.Create(u.db, &detail); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		detail.Vision = req.Vision
		detail.Mission = req.Mission
		detail.Rules = req.Rules
		detail.WorkProgram = req.WorkProgram
		detail.ImageURL = req.ImageURL

		if err := u.repo.Update(u.db, &detail); err != nil {
			return nil, err
		}
	}

	response := converter.ToOrganizationDetailResponse(&detail)
	return &response, nil
}
