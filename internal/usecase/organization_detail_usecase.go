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
			return &model.OrganizationDetailResponse{
				EntityType:            entityType,
				Vision:                "",
				Mission:               "",
				Rules:                 "",
				WorkProgram:           "",
				VisionMissionImageURL: "",
				WorkProgramImageURL:   "",
				RulesImageURL:         "",
				StructureImageURL:     "",
			}, nil
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
			ID:                    uuid.New().String(),
			EntityType:            entityType,
			Vision:                req.Vision,
			Mission:               req.Mission,
			Rules:                 req.Rules,
			WorkProgram:           req.WorkProgram,
			VisionMissionImageURL: req.VisionMissionImageURL,
			WorkProgramImageURL:   req.WorkProgramImageURL,
			RulesImageURL:         req.RulesImageURL,
			StructureImageURL:     req.StructureImageURL,
		}

		if err := u.repo.Create(u.db, &detail); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		if req.Vision != "" {
			detail.Vision = req.Vision
		}
		if req.Mission != "" {
			detail.Mission = req.Mission
		}
		if req.Rules != "" {
			detail.Rules = req.Rules
		}
		if req.WorkProgram != "" {
			detail.WorkProgram = req.WorkProgram
		}
		if req.VisionMissionImageURL != "" {
			detail.VisionMissionImageURL = req.VisionMissionImageURL
		}
		if req.WorkProgramImageURL != "" {
			detail.WorkProgramImageURL = req.WorkProgramImageURL
		}
		if req.RulesImageURL != "" {
			detail.RulesImageURL = req.RulesImageURL
		}

		if req.StructureImageURL != "" {
			detail.StructureImageURL = req.StructureImageURL
		}

		if err := u.repo.Update(u.db, &detail); err != nil {
			return nil, err
		}
	}

	response := converter.ToOrganizationDetailResponse(&detail)
	return &response, nil
}
