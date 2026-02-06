package test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

func setupMockTestimonialUsecase(t *testing.T) (usecase.TestimonialUsecase, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub db: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	u := usecase.NewTestimonialUsecase(gormDB, validator.New())
	return u, mock
}

func TestTestimonialUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)

	req := model.TestimonialRequest{
		Name:       "John Doe",
		AvatarURL:  "https://example.com/avatar.jpg",
		Rating:     5,
		Comment:    "Excellent!",
		IsActive:   true,
		OrderIndex: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `testimonials`")).
		WithArgs(
			"John Doe",
			"https://example.com/avatar.jpg",
			5,
			"Excellent!",
			true,
			1,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Rating, res.Rating)
}

func TestTestimonialUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockTestimonialUsecase(t)

	req := model.TestimonialRequest{
		Name:    "",
		Rating:  0,
		Comment: "",
	}

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestTestimonialUsecase_GetAll(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "rating", "order_index"}).
		AddRow(2, "B", 4, 1).
		AddRow(1, "A", 5, 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` ORDER BY order_index ASC")).
		WillReturnRows(rows)

	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestTestimonialUsecase_GetPublic(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "is_active"}).
		AddRow(1, "Public User", true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE is_active = ? ORDER BY order_index ASC")).
		WithArgs(true).
		WillReturnRows(rows)

	list, err := u.GetPublic()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestTestimonialUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 10

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Found Me")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(targetID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, targetID, res.ID)
}

func TestTestimonialUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ? LIMIT ?")).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "testimonial not found", e.Message)
	}
}

func TestTestimonialUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 1
	now := time.Now()

	req := model.TestimonialRequest{
		Name:       "New Name",
		AvatarURL:  "https://example.com/new.jpg",
		Rating:     4,
		Comment:    "Updated",
		IsActive:   false,
		OrderIndex: 3,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(targetID, "Old", now))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `testimonials`")).
		WithArgs(
			"New Name",
			"https://example.com/new.jpg",
			4,
			"Updated",
			false,
			3,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
}

func TestTestimonialUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 99

	req := model.TestimonialRequest{
		Name:    "Valid Name",
		Rating:  5,
		Comment: "Valid Comment",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "testimonial not found", e.Message)
	}
}

func TestTestimonialUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 1

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Delete Me"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `testimonials` WHERE `testimonials`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestTestimonialUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 99

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "testimonial not found", e.Message)
	}
}
