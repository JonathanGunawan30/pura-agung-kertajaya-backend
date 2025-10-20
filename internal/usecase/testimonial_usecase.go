package usecase

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TestimonialUsecase interface {
	GetAll() ([]model.TestimonialResponse, error)
	GetByID(id int) (*model.TestimonialResponse, error)
	Create(req model.TestimonialRequest) (*model.TestimonialResponse, error)
	Update(id int, req model.TestimonialRequest) (*model.TestimonialResponse, error)
	Delete(id int) error
}

type testimonialUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Testimonial]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewTestimonialUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) TestimonialUsecase {
	return &testimonialUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Testimonial]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *testimonialUsecase) GetAll() ([]model.TestimonialResponse, error) {
	var testimonials []entity.Testimonial
	if err := u.db.Order("order_index ASC").Find(&testimonials).Error; err != nil {
		return nil, err
	}

	responses := make([]model.TestimonialResponse, 0, len(testimonials))
	for _, t := range testimonials {
		responses = append(responses, model.TestimonialResponse{
			ID:         t.ID,
			Name:       t.Name,
			AvatarURL:  t.AvatarURL,
			Rating:     t.Rating,
			Comment:    t.Comment,
			IsActive:   t.IsActive,
			OrderIndex: t.OrderIndex,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		})
	}
	return responses, nil
}

func (u *testimonialUsecase) GetByID(id int) (*model.TestimonialResponse, error) {
	var t entity.Testimonial
	if err := u.repo.FindById(u.db, &t, id); err != nil {
		return nil, err
	}

	return &model.TestimonialResponse{
		ID:         t.ID,
		Name:       t.Name,
		AvatarURL:  t.AvatarURL,
		Rating:     t.Rating,
		Comment:    t.Comment,
		IsActive:   t.IsActive,
		OrderIndex: t.OrderIndex,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}, nil
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

	return &model.TestimonialResponse{
		ID:         t.ID,
		Name:       t.Name,
		AvatarURL:  t.AvatarURL,
		Rating:     t.Rating,
		Comment:    t.Comment,
		IsActive:   t.IsActive,
		OrderIndex: t.OrderIndex,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}, nil
}

func (u *testimonialUsecase) Update(id int, req model.TestimonialRequest) (*model.TestimonialResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var t entity.Testimonial
	if err := u.repo.FindById(u.db, &t, id); err != nil {
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

	return &model.TestimonialResponse{
		ID:         t.ID,
		Name:       t.Name,
		AvatarURL:  t.AvatarURL,
		Rating:     t.Rating,
		Comment:    t.Comment,
		IsActive:   t.IsActive,
		OrderIndex: t.OrderIndex,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}, nil
}

func (u *testimonialUsecase) Delete(id int) error {
	var t entity.Testimonial
	if err := u.repo.FindById(u.db, &t, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &t)
}
