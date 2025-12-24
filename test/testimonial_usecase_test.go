package test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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

	u := usecase.NewTestimonialUsecase(gormDB, logrus.New(), validator.New())
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
	// FIX: Argumen disesuaikan menjadi 8 (tanpa ID auto-increment)
	// name, avatar_url, rating, comment, is_active, order_index, created_at, updated_at
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `testimonials`")).
		WithArgs("John Doe", "https://example.com/avatar.jpg", 5, "Excellent!", true, 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "name", "rating"}).
		AddRow(1, "John Doe", 5)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	res, err := u.Create(req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, res)
	if res == nil {
		return
	}
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Rating, res.Rating)
}

func TestTestimonialUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockTestimonialUsecase(t)

	req := model.TestimonialRequest{
		Name:     "",
		Rating:   0,
		Comment:  "",
		IsActive: true,
	}

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestTestimonialUsecase_GetAll_OrderedByIndex(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "rating", "order_index"}).
		AddRow(2, "B", 4, 1).
		AddRow(1, "A", 5, 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` ORDER BY order_index ASC")).
		WillReturnRows(rows)

	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "B", list[0].Name)
	assert.Equal(t, "A", list[1].Name)
}

func TestTestimonialUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ?")).
		WithArgs(999999, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(999999)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestTestimonialUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 1

	req := model.TestimonialRequest{
		Name:       "New",
		AvatarURL:  "https://example.com/new.jpg",
		Rating:     4,
		Comment:    "Updated",
		IsActive:   false,
		OrderIndex: 3,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Old"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `testimonials`")).
		WithArgs("New", "https://example.com/new.jpg", 4, "Updated", false, 3, sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, res)
	if res == nil {
		return
	}
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Comment, res.Comment)
}

func TestTestimonialUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockTestimonialUsecase(t)
	targetID := 1

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `testimonials` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Delete Me"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `testimonials` WHERE")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}
