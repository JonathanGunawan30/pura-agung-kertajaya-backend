package test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

func newTestimonialUsecase() usecase.TestimonialUsecase {
	return usecase.NewTestimonialUsecase(db, logrus.New(), validator.New())
}

func prepareTestimonial(name string, rating int, order int) entity.Testimonial {
	return entity.Testimonial{
		Name:       name,
		AvatarURL:  "https://example.com/avatar.jpg",
		Rating:     rating,
		Comment:    "Great service",
		IsActive:   true,
		OrderIndex: order,
	}
}

func ClearTestimonials() {
	err := db.Where("id IS NOT NULL").Delete(&entity.Testimonial{}).Error
	if err != nil {
		log.Fatalf("Failed clear testimonials: %+v", err)
	}
}

func TestTestimonialUsecase_Create_Success(t *testing.T) {
	ClearTestimonials()
	u := newTestimonialUsecase()

	req := model.TestimonialRequest{
		Name:       "John Doe",
		AvatarURL:  "https://example.com/avatar.jpg",
		Rating:     5,
		Comment:    "Excellent!",
		IsActive:   true,
		OrderIndex: 1,
	}

	res, err := u.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Rating, res.Rating)
}

func TestTestimonialUsecase_Create_ValidationError(t *testing.T) {
	ClearTestimonials()
	u := newTestimonialUsecase()

	req := model.TestimonialRequest{
		Name:      "",
		Rating:    0,
		Comment:   "",
		IsActive:  true,
	}

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestTestimonialUsecase_GetAll_OrderedByIndex(t *testing.T) {
	ClearTestimonials()

	// seed records with different order_index
	t1 := prepareTestimonial("A", 5, 2)
	t2 := prepareTestimonial("B", 4, 1)
	assert.NoError(t, db.Create(&t1).Error)
	assert.NoError(t, db.Create(&t2).Error)

	u := newTestimonialUsecase()
	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "B", list[0].Name) // order_index 1 comes first
	assert.Equal(t, "A", list[1].Name)
}

func TestTestimonialUsecase_GetByID_NotFound(t *testing.T) {
	ClearTestimonials()
	u := newTestimonialUsecase()

	res, err := u.GetByID(999999)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestTestimonialUsecase_Update_Success(t *testing.T) {
	ClearTestimonials()

	seed := prepareTestimonial("Old", 3, 1)
	assert.NoError(t, db.Create(&seed).Error)

	u := newTestimonialUsecase()

	req := model.TestimonialRequest{
		Name:       "New",
		AvatarURL:  "https://example.com/new.jpg",
		Rating:     4,
		Comment:    "Updated",
		IsActive:   false,
		OrderIndex: 3,
	}

	res, err := u.Update(seed.ID, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.AvatarURL, res.AvatarURL)
	assert.Equal(t, req.Rating, res.Rating)
	assert.Equal(t, req.Comment, res.Comment)
	assert.Equal(t, req.IsActive, res.IsActive)
	assert.Equal(t, req.OrderIndex, res.OrderIndex)
}

func TestTestimonialUsecase_Delete_Success(t *testing.T) {
	ClearTestimonials()
	seed := prepareTestimonial("Delete Me", 5, 1)
	assert.NoError(t, db.Create(&seed).Error)

	u := newTestimonialUsecase()
	err := u.Delete(seed.ID)
	assert.NoError(t, err)

	var count int64
	assert.NoError(t, db.Model(&entity.Testimonial{}).Where("id = ?", seed.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}
