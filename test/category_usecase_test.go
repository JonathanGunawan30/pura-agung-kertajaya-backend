package test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

func setupMockCategoryUsecase(t *testing.T) (usecase.CategoryUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewCategoryUsecase(gormDB, logrus.New(), validator.New())
	return u, mock
}

func TestCategoryUsecase_GetAll(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "slug", "created_at", "updated_at"}).
		AddRow("c1", "Adat", "adat", time.Now(), time.Now()).
		AddRow("c2", "Upacara", "upacara", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` ORDER BY name ASC")).
		WillReturnRows(rows)

	list, err := u.GetAll()

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "c1", list[0].ID)
	assert.Equal(t, "Adat", list[0].Name)
}

func TestCategoryUsecase_Create_Simple(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)

	req := model.CreateCategoryRequest{Name: "Upacara Besar"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE slug = ?")).
		WithArgs("upacara-besar").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories`")).
		WithArgs(
			sqlmock.AnyArg(),
			req.Name,
			"upacara-besar",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "upacara-besar", created.Slug)
}

func TestCategoryUsecase_Create_SlugCollision(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)

	req := model.CreateCategoryRequest{Name: "Upacara"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE slug = ?")).
		WithArgs("upacara").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE slug = ?")).
		WithArgs("upacara-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE slug = ?")).
		WithArgs("upacara-2").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories`")).
		WithArgs(sqlmock.AnyArg(), "Upacara", "upacara-2", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req)
	assert.NoError(t, err)
	assert.Equal(t, "upacara-2", created.Slug)
}

func TestCategoryUsecase_Update(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	targetID := "cat-123"

	req := model.UpdateCategoryRequest{Name: "Baru"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug"}).AddRow(targetID, "Lama", "lama"))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE slug = ? AND id != ?")).
		WithArgs("baru", targetID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories`")).
		WithArgs(
			"Baru",
			"baru",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)

	assert.NoError(t, err)

	if updated == nil {
		t.FailNow()
	}

	assert.NotNil(t, updated)
	assert.Equal(t, "Baru", updated.Name)
	assert.Equal(t, "baru", updated.Slug)
}

func TestCategoryUsecase_Delete(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	targetID := "cat-delete"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "To Delete"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}
