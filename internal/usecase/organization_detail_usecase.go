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

type OrganizationDetailUsecase interface {
	GetByEntityType(entityType string) (*model.OrganizationDetailResponse, error)
	Update(entityType string, req model.UpdateOrganizationDetailRequest) (*model.OrganizationDetailResponse, error)
}

type organizationDetailUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.OrganizationDetail]
	validate *validator.Validate
}

func NewOrganizationDetailUsecase(db *gorm.DB, validate *validator.Validate) OrganizationDetailUsecase {
	return &organizationDetailUsecase{
		db:       db,
		repo:     &repository.Repository[entity.OrganizationDetail]{DB: db},
		validate: validate,
	}
}

func (u *organizationDetailUsecase) GetByEntityType(entityType string) (*model.OrganizationDetailResponse, error) {
	var detail entity.OrganizationDetail

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
		detail.Vision = req.Vision
		detail.Mission = req.Mission
		detail.Rules = req.Rules
		detail.WorkProgram = req.WorkProgram
		detail.VisionMissionImageURL = req.VisionMissionImageURL
		detail.WorkProgramImageURL = req.WorkProgramImageURL
		detail.RulesImageURL = req.RulesImageURL
		detail.StructureImageURL = req.StructureImageURL

		if err := u.repo.Update(u.db, &detail); err != nil {
			return nil, err
		}
	}

	response := converter.ToOrganizationDetailResponse(&detail)
	return &response, nil
}
