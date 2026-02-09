package usecase

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TestimonialUsecase interface {
	GetAll() ([]model.TestimonialResponse, error)
	GetPublic() ([]model.TestimonialResponse, error)
	GetByID(id string) (*model.TestimonialResponse, error)
	Create(req model.TestimonialRequest) (*model.TestimonialResponse, error)
	Update(id string, req model.TestimonialRequest) (*model.TestimonialResponse, error)
	Delete(id string) error
}

type testimonialUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Testimonial]
	validate *validator.Validate
}

func NewTestimonialUsecase(db *gorm.DB, validate *validator.Validate) TestimonialUsecase {
	return &testimonialUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Testimonial]{DB: db},
		validate: validate,
	}
}

func (u *testimonialUsecase) GetAll() ([]model.TestimonialResponse, error) {
	var testimonials []entity.Testimonial

	query := u.db.Order("order_index ASC")

	if err := u.repo.FindAll(query, &testimonials); err != nil {
		return nil, err
	}

	return converter.ToTestimonialResponses(testimonials), nil
}

func (u *testimonialUsecase) GetPublic() ([]model.TestimonialResponse, error) {
	var testimonials []entity.Testimonial

	query := u.db.Where("is_active = ?", true).Order("order_index ASC")

	if err := u.repo.FindAll(query, &testimonials); err != nil {
		return nil, err
	}

	return converter.ToTestimonialResponses(testimonials), nil
}

func (u *testimonialUsecase) GetByID(id string) (*model.TestimonialResponse, error) {
	var t entity.Testimonial
	if err := u.repo.FindById(u.db, &t, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("testimonial not found")
		}
		return nil, err
	}

	response := converter.ToTestimonialResponse(&t)
	return &response, nil
}

func (u *testimonialUsecase) Create(req model.TestimonialRequest) (*model.TestimonialResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	t := entity.Testimonial{
		Name:       req.Name,
		AvatarURL:  req.AvatarURL,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsActive:   req.IsActive,
		OrderIndex: req.OrderIndex,
	}

	if err := u.repo.Create(u.db, &t); err != nil {
		return nil, err
	}

	response := converter.ToTestimonialResponse(&t)
	return &response, nil
}

func (u *testimonialUsecase) Update(id string, req model.TestimonialRequest) (*model.TestimonialResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var t entity.Testimonial
	if err := u.repo.FindById(u.db, &t, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("testimonial not found")
		}
		return nil, err
	}

	t.Name = req.Name
	t.AvatarURL = req.AvatarURL
	t.Rating = req.Rating
	t.Comment = req.Comment
	t.IsActive = req.IsActive
	t.OrderIndex = req.OrderIndex

	if err := u.repo.Update(u.db, &t); err != nil {
		return nil, err
	}

	response := converter.ToTestimonialResponse(&t)
	return &response, nil
}

func (u *testimonialUsecase) Delete(id string) error {
	var t entity.Testimonial
	if err := u.repo.FindById(u.db, &t, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound("testimonial not found")
		}
		return err
	}
	return u.repo.Delete(u.db, &t)
}
