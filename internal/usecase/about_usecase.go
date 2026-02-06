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

type AboutUsecase interface {
	GetAll(entityType string) ([]model.AboutSectionResponse, error)
	GetPublic(entityType string) ([]model.AboutSectionResponse, error)
	GetByID(id string) (*model.AboutSectionResponse, error)
	Create(req model.AboutSectionRequest) (*model.AboutSectionResponse, error)
	Update(id string, req model.AboutSectionRequest) (*model.AboutSectionResponse, error)
	Delete(id string) error
}

type aboutUsecase struct {
	db        *gorm.DB
	repoAbout *repository.Repository[entity.AboutSection]
	repoValue *repository.Repository[entity.AboutValue]
	validate  *validator.Validate
}

func NewAboutUsecase(db *gorm.DB, validate *validator.Validate) AboutUsecase {
	return &aboutUsecase{
		db:        db,
		repoAbout: &repository.Repository[entity.AboutSection]{DB: db},
		repoValue: &repository.Repository[entity.AboutValue]{DB: db},
		validate:  validate,
	}
}

func preloadValuesOrdered(tx *gorm.DB) *gorm.DB {
	return tx.Preload("Values", func(db *gorm.DB) *gorm.DB { return db.Order("order_index ASC") })
}

func (u *aboutUsecase) GetAll(entityType string) ([]model.AboutSectionResponse, error) {
	var list []entity.AboutSection
	query := preloadValuesOrdered(u.db).Order("created_at ASC")
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	resp := make([]model.AboutSectionResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, converter.ToAboutSectionResponse(a))
	}
	return resp, nil
}

func (u *aboutUsecase) GetPublic(entityType string) ([]model.AboutSectionResponse, error) {
	var list []entity.AboutSection
	query := preloadValuesOrdered(u.db).Where("is_active = ?", true).Order("created_at ASC")
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	resp := make([]model.AboutSectionResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, converter.ToAboutSectionResponse(a))
	}
	return resp, nil
}

func (u *aboutUsecase) GetByID(id string) (*model.AboutSectionResponse, error) {
	var a entity.AboutSection
	if err := preloadValuesOrdered(u.db).Where("id = ?", id).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("about section not found")
		}
		return nil, err
	}
	r := converter.ToAboutSectionResponse(a)
	return &r, nil
}

func (u *aboutUsecase) Create(req model.AboutSectionRequest) (*model.AboutSectionResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	id := uuid.New().String()
	a := entity.AboutSection{
		ID:          id,
		EntityType:  req.EntityType,
		Title:       req.Title,
		Description: req.Description,
		Images:      util.ImageMap(req.Images),
		IsActive:    req.IsActive,
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := u.repoAbout.Create(tx, &a); err != nil {
			return err
		}
		if len(req.Values) > 0 {
			values := make([]entity.AboutValue, 0, len(req.Values))
			for _, v := range req.Values {
				values = append(values, entity.AboutValue{
					ID:         uuid.New().String(),
					AboutID:    a.ID,
					Title:      v.Title,
					Value:      v.Value,
					OrderIndex: v.OrderIndex,
				})
			}
			if err := tx.Create(&values).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := preloadValuesOrdered(u.db).Where("id = ?", a.ID).Take(&a).Error; err != nil {
		return nil, err
	}
	r := converter.ToAboutSectionResponse(a)
	return &r, nil
}

func (u *aboutUsecase) Update(id string, req model.AboutSectionRequest) (*model.AboutSectionResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var a entity.AboutSection
	if err := u.repoAbout.FindById(u.db, &a, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("about section not found")
		}
		return nil, err
	}

	a.EntityType = req.EntityType
	a.Title = req.Title
	a.Description = req.Description
	a.Images = util.ImageMap(req.Images)
	a.IsActive = req.IsActive

	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := u.repoAbout.Update(tx, &a); err != nil {
			return err
		}
		if err := tx.Where("about_id = ?", a.ID).Delete(&entity.AboutValue{}).Error; err != nil {
			return err
		}
		if len(req.Values) > 0 {
			values := make([]entity.AboutValue, 0, len(req.Values))
			for _, v := range req.Values {
				values = append(values, entity.AboutValue{
					ID:         uuid.New().String(),
					AboutID:    a.ID,
					Title:      v.Title,
					Value:      v.Value,
					OrderIndex: v.OrderIndex,
				})
			}
			if err := tx.Create(&values).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := preloadValuesOrdered(u.db).Where("id = ?", a.ID).Take(&a).Error; err != nil {
		return nil, err
	}
	r := converter.ToAboutSectionResponse(a)
	return &r, nil
}

func (u *aboutUsecase) Delete(id string) error {
	exists, err := u.repoAbout.CountById(u.db, id)
	if err != nil {
		return err
	}
	if exists == 0 {
		return model.ErrNotFound("about section not found")
	}
	return u.repoAbout.Delete(u.db, &entity.AboutSection{ID: id})
}
