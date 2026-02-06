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

func setupMockFacilityUsecase(t *testing.T) (usecase.FacilityUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewFacilityUsecase(gormDB, validator.New())
	return u, mock
}

func TestFacilityUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "images", "is_active", "order_index"}).
		AddRow("g3", "C", []byte(`{"lg":"https://img3"}`), true, 1).
		AddRow("g1", "B", []byte(`{"lg":"https://img1"}`), true, 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE entity_type = ? AND is_active = ? ORDER BY order_index ASC")).
		WithArgs("pura", true).
		WillReturnRows(rows)

	list, err := u.GetPublic("pura")

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "g3", list[0].ID)
	assert.Equal(t, "https://img3", list[0].Images.Lg)
}

func TestFacilityUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	id := "uuid-1"

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Facility Name")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Facility Name", res.Name)
}

func TestFacilityUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	id := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "facility not found", e.Message)
	}
}

func TestFacilityUsecase_Create(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)

	req := model.CreateFacilityRequest{
		EntityType: "pura",
		Name:       "Name",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
		IsActive:   true,
		OrderIndex: 3,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `facilities`")).
		WithArgs(
			sqlmock.AnyArg(),
			req.EntityType,
			"Name",
			"",
			sqlmock.AnyArg(),
			3,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	if assert.NotNil(t, created) {
		assert.Equal(t, "https://img.com/lg.jpg", created.Images.Lg)
	}
}

func TestFacilityUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockFacilityUsecase(t)
	req := model.CreateFacilityRequest{} // Empty
	res, err := u.Create("pura", req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestFacilityUsecase_Update(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "g1"

	req := model.UpdateFacilityRequest{
		Name:       "New Name",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
		IsActive:   false,
		OrderIndex: 5,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "name", "images"}).
			AddRow(targetID, "pura", "Old Name", []byte(`{"lg":"old.jpg"}`)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `facilities`")).
		WithArgs(
			"pura",
			"New Name",
			"",
			sqlmock.AnyArg(),
			5,
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, false, updated.IsActive)
		assert.Equal(t, 5, updated.OrderIndex)
	}
}

func TestFacilityUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "missing"

	req := model.UpdateFacilityRequest{
		Name:       "New Name",
		Images:     map[string]string{"lg": "img.jpg"},
		OrderIndex: 1,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	updated, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, updated)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "facility not found", e.Message)
	}
}

func TestFacilityUsecase_Delete(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "g1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "images"}).
			AddRow(targetID, "Name", []byte(`{}`)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `facilities` WHERE `facilities`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestFacilityUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "facility not found", e.Message)
	}
}
