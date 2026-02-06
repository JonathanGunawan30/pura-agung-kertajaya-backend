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

	u := usecase.NewCategoryUsecase(gormDB, validator.New())
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
	if len(list) > 0 {
		assert.Equal(t, "c1", list[0].ID)
		assert.Equal(t, "Adat", list[0].Name)
	}
}

func TestCategoryUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	id := "c1"

	rows := sqlmock.NewRows([]string{"id", "name", "slug"}).
		AddRow(id, "Adat", "adat")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(id)

	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Adat", res.Name)
	}
}

func TestCategoryUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	id := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "category not found", e.Message)
	}
}

func TestCategoryUsecase_Create_Simple(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	req := model.CreateCategoryRequest{Name: "Upacara Besar"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE slug = ?")).
		WithArgs("upacara-besar").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories`")).
		WithArgs(sqlmock.AnyArg(), req.Name, "upacara-besar", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req)
	assert.NoError(t, err)
	if assert.NotNil(t, created) {
		assert.Equal(t, "upacara-besar", created.Slug)
	}
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
	if assert.NotNil(t, created) {
		assert.Equal(t, "upacara-2", created.Slug)
	}
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
		WithArgs("Baru", "baru", sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)

	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "Baru", updated.Name)
	}
}

func TestCategoryUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	targetID := "missing"
	req := model.UpdateCategoryRequest{Name: "Baru"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	updated, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, updated)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "category not found", e.Message)
	}
}

func TestCategoryUsecase_Delete(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	targetID := "cat-delete"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE id = ?")).
		WithArgs(targetID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE category_id = ?")).
		WithArgs(targetID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestCategoryUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE id = ?")).
		WithArgs(targetID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "category not found", e.Message)
	}
}

func TestCategoryUsecase_Delete_Fail_Referenced(t *testing.T) {
	u, mock := setupMockCategoryUsecase(t)
	targetID := "cat-busy"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE id = ?")).
		WithArgs(targetID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE category_id = ?")).
		WithArgs(targetID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	err := u.Delete(targetID)
	assert.Error(t, err)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 409, e.Code)
		assert.Equal(t, "category is currently in use", e.Message)
	}
}
