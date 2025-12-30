package usecase

import (
	"fmt"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUsecase interface {
	GetAll() ([]model.CategoryResponse, error)
	GetByID(id string) (*model.CategoryResponse, error)
	Create(req model.CreateCategoryRequest) (*model.CategoryResponse, error)
	Update(id string, req model.UpdateCategoryRequest) (*model.CategoryResponse, error)
	Delete(id string) error
}

type categoryUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Category]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewCategoryUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) CategoryUsecase {
	return &categoryUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Category]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *categoryUsecase) GetAll() ([]model.CategoryResponse, error) {
	var items []entity.Category

	query := u.db.Order("name ASC")

	if err := u.repo.FindAll(query, &items); err != nil {
		return nil, err
	}

	return converter.ToCategoryResponses(items), nil
}

func (u *categoryUsecase) GetByID(id string) (*model.CategoryResponse, error) {
	var c entity.Category
	if err := u.repo.FindById(u.db, &c, id); err != nil {
		return nil, err
	}
	r := converter.ToCategoryResponse(&c)
	return &r, nil
}

func (u *categoryUsecase) Create(req model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	baseSlug := slug.Make(req.Name)
	finalSlug := baseSlug
	counter := 1

	for {
		count, err := u.repo.CountBySlug(u.db, finalSlug)
		if err != nil {
			return nil, err
		}

		if count == 0 {
			break
		}
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	c := entity.Category{
		Name: req.Name,
		Slug: finalSlug,
	}

	if err := u.repo.Create(u.db, &c); err != nil {
		return nil, err
	}

	r := converter.ToCategoryResponse(&c)
	return &r, nil
}

func (u *categoryUsecase) Update(id string, req model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var c entity.Category
	if err := u.repo.FindById(u.db, &c, id); err != nil {
		return nil, err
	}

	if c.Name != req.Name {
		baseSlug := slug.Make(req.Name)
		finalSlug := baseSlug
		counter := 1

		for {
			count, err := u.repo.CountBySlugIgnoringID(u.db, finalSlug, id)

			if err != nil {
				return nil, err
			}

			if count == 0 {
				break
			}
			finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
			counter++
		}
		c.Slug = finalSlug
	}

	c.Name = req.Name

	if err := u.repo.Update(u.db, &c); err != nil {
		return nil, err
	}

	r := converter.ToCategoryResponse(&c)
	return &r, nil
}

func (u *categoryUsecase) Delete(id string) error {
	var c entity.Category
	if err := u.repo.FindById(u.db, &c, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &c)
}
