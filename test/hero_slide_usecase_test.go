package test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

func setupMockHeroSlideUsecase(t *testing.T) (usecase.HeroSlideUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewHeroSlideUsecase(gormDB, validator.New())
	return u, mock
}

func TestHeroSlideUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "entity_type", "images", "order_index", "is_active"}).
		AddRow("id-1", "pura", []byte(`{"lg":"https://img1.jpg"}`), 1, true).
		AddRow("id-2", "pura", []byte(`{"lg":"https://img2.jpg"}`), 2, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE is_active = ? ORDER BY order_index ASC")).
		WithArgs(true).
		WillReturnRows(rows)

	list, err := u.GetPublic("")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "https://img1.jpg", list[0].Images.Lg)
}

func TestHeroSlideUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)

	req := model.HeroSlideRequest{
		EntityType: "pura",
		Images:     map[string]string{"lg": "https://example.com/image.jpg"},
		OrderIndex: 1,
		IsActive:   true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `hero_slides`")).
		WithArgs(
			sqlmock.AnyArg(),
			"pura",
			sqlmock.AnyArg(),
			1,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	if assert.NotNil(t, created) {
		assert.Equal(t, req.EntityType, created.EntityType)
		assert.Equal(t, req.Images["lg"], created.Images.Lg)
	}
}

func TestHeroSlideUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockHeroSlideUsecase(t)
	req := model.HeroSlideRequest{}
	res, err := u.Create("pura", req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestHeroSlideUsecase_GetAll_OrderedByIndex(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "entity_type", "images", "order_index", "is_active"}).
		AddRow("id-1", "pura", []byte(`{"lg":"https://example.com/1.jpg"}`), 1, true).
		AddRow("id-2", "pura", []byte(`{"lg":"https://example.com/2.jpg"}`), 2, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE entity_type = ? ORDER BY order_index ASC")).
		WithArgs("pura").
		WillReturnRows(rows)

	list, err := u.GetAll("pura")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestHeroSlideUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)
	id := "id-1"

	rows := sqlmock.NewRows([]string{"id", "entity_type", "images"}).
		AddRow(id, "pura", []byte(`{"lg":"img.jpg"}`))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestHeroSlideUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)
	id := "not-exist"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "hero slide not found", e.Message)
	}
}

func TestHeroSlideUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)
	targetID := "slide-1"

	req := model.HeroSlideRequest{
		EntityType: "yayasan",
		Images:     map[string]string{"lg": "https://new.jpg"},
		OrderIndex: 5,
		IsActive:   false,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "images"}).
			AddRow(targetID, "pura", []byte(`{"lg":"https://old.jpg"}`)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `hero_slides`")).
		WithArgs(
			"yayasan",
			sqlmock.AnyArg(),
			5,
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, req.EntityType, res.EntityType)
		assert.Equal(t, req.Images["lg"], res.Images.Lg)
	}
}

func TestHeroSlideUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)
	targetID := "missing"

	req := model.HeroSlideRequest{
		EntityType: "pura",
		Images:     map[string]string{"lg": "img.jpg"},
		OrderIndex: 1,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "hero slide not found", e.Message)
	}
}

func TestHeroSlideUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)
	targetID := "to-delete"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "images"}).
			AddRow(targetID, []byte(`{"lg":"https://img.jpg"}`)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `hero_slides` WHERE `hero_slides`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestHeroSlideUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockHeroSlideUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hero_slides` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "hero slide not found", e.Message)
	}
}
