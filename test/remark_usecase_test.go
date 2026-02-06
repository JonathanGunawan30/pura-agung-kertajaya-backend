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

func setupMockRemarkUsecase(t *testing.T) (usecase.RemarkUsecase, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	u := usecase.NewRemarkUsecase(gormDB, validator.New())

	return u, mock
}

func TestRemarkUsecase_GetAll_Success(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "entity_type", "name", "position", "order_index", "created_at", "updated_at"}).
		AddRow("uuid-1", "pura", "Pak Ketua", "Ketua", 1, now, now).
		AddRow("uuid-2", "pura", "Pak Wakil", "Wakil", 2, now, now)

	expectedSQL := "SELECT * FROM `remarks` WHERE entity_type = ? ORDER BY order_index ASC"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs("pura").
		WillReturnRows(rows)

	result, err := u.GetAll("pura")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Pak Ketua", result[0].Name)
}

func TestRemarkUsecase_GetPublic_Success(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "entity_type", "name", "is_active"}).
		AddRow("uuid-1", "pura", "Pak Ketua", true)

	expectedSQL := "SELECT * FROM `remarks` WHERE is_active = ? AND entity_type = ? ORDER BY order_index ASC"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(true, "pura").
		WillReturnRows(rows)

	result, err := u.GetPublic("pura")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestRemarkUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)
	targetUUID := "uuid-123"

	rows := sqlmock.NewRows([]string{"id", "entity_type", "name", "position"}).
		AddRow(targetUUID, "pura", "Pak Bos", "Ketua")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ? LIMIT ?")).
		WithArgs(targetUUID, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(targetUUID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, targetUUID, res.ID)
}

func TestRemarkUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)
	targetUUID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ? LIMIT ?")).
		WithArgs(targetUUID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(targetUUID)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "remark not found", e.Message)
	}
}

func TestRemarkUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)

	req := model.CreateRemarkRequest{
		EntityType: "yayasan",
		Name:       "Ibu Ketua",
		Position:   "Ketua Yayasan",
		Content:    "Halo",
		OrderIndex: 1,
		IsActive:   true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `remarks`")).
		WithArgs(
			sqlmock.AnyArg(),
			"yayasan",
			"Ibu Ketua",
			"Ketua Yayasan",
			"",
			"Halo",
			true,
			1,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Ibu Ketua", res.Name)
}

func TestRemarkUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockRemarkUsecase(t)
	req := model.CreateRemarkRequest{}
	res, err := u.Create("pura", req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestRemarkUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)
	targetID := "uuid-1"

	req := model.UpdateRemarkRequest{
		Name:       "New Name",
		Position:   "New Pos",
		Content:    "New Content",
		OrderIndex: 2,
		IsActive:   false,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "name"}).AddRow(targetID, "pura", "Old Name"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `remarks`")).
		WithArgs(
			"pura",
			"New Name",
			"New Pos",
			"",
			"New Content",
			false,
			2,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New Name", res.Name)
}

func TestRemarkUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)
	targetID := "missing"

	req := model.UpdateRemarkRequest{
		Name:     "Valid Name",
		Position: "Valid Pos",
		Content:  "Valid",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "remark not found", e.Message)
	}
}

func TestRemarkUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)
	targetID := "uuid-delete-me"

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Deleted Guy")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `remarks` WHERE `remarks`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestRemarkUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockRemarkUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "remark not found", e.Message)
	}
}
