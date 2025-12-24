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

func setupMockGalleryUsecase(t *testing.T) (usecase.GalleryUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewGalleryUsecase(gormDB, logrus.New(), validator.New())
	return u, mock
}

func TestGalleryUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockGalleryUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "entity_type", "title", "image_url", "is_active", "order_index"}).
		AddRow("g3", "pura", "C", "https://img3", true, 1).
		AddRow("g1", "pura", "B", "https://img1", true, 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `galleries` WHERE entity_type = ? AND is_active = ? ORDER BY order_index ASC")).
		WithArgs("pura", true).
		WillReturnRows(rows)

	list, err := u.GetPublic("pura")

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "g3", list[0].ID)
	assert.Equal(t, "g1", list[1].ID)
}

func TestGalleryUsecase_Create(t *testing.T) {
	u, mock := setupMockGalleryUsecase(t)

	req := model.CreateGalleryRequest{
		EntityType: "pura",
		Title:      "Title",
		ImageURL:   "https://img",
		IsActive:   true,
		OrderIndex: 3,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `galleries`")).
		WithArgs(
			sqlmock.AnyArg(),
			"pura",
			"Title",
			"",
			"https://img",
			3,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "entity_type", "title", "description", "image_url", "order_index", "is_active"}).
		AddRow("uuid-1", "pura", "Title", "", "https://img", 3, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `galleries`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	created, err := u.Create(req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, created)
	if created == nil {
		return
	}
	assert.NotEmpty(t, created.ID)
}

func TestGalleryUsecase_Update(t *testing.T) {
	u, mock := setupMockGalleryUsecase(t)
	targetID := "g-1"

	req := model.UpdateGalleryRequest{
		Title:      "New Title",
		ImageURL:   "https://img2",
		IsActive:   false,
		OrderIndex: 5,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `galleries` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "title"}).AddRow(targetID, "pura", "Old Title"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `galleries`")).
		WithArgs("pura", "New Title", "", "https://img2", 5, false, sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, updated)
	if updated == nil {
		return
	}
	assert.Equal(t, "New Title", updated.Title)
	assert.Equal(t, false, updated.IsActive)
	assert.Equal(t, 5, updated.OrderIndex)
}

func TestGalleryUsecase_Delete(t *testing.T) {
	u, mock := setupMockGalleryUsecase(t)
	targetID := "g-1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `galleries` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "entity_type"}).AddRow(targetID, "Title", "pura"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `galleries` WHERE `galleries`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}
