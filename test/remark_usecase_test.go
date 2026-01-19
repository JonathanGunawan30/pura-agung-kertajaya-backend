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

func setupMockUsecase(t *testing.T) (usecase.RemarkUsecase, sqlmock.Sqlmock) {
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

	logger := logrus.New()
	validate := validator.New()

	u := usecase.NewRemarkUsecase(gormDB, logger, validate)

	return u, mock
}

func TestRemarkUsecase_GetAll_Success(t *testing.T) {
	u, mock := setupMockUsecase(t)

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
	assert.Equal(t, "pura", result[0].EntityType)
}

func TestRemarkUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockUsecase(t)
	targetUUID := "uuid-123"

	rows := sqlmock.NewRows([]string{"id", "entity_type", "name", "position"}).
		AddRow(targetUUID, "pura", "Pak Bos", "Ketua")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ?")).
		WithArgs(targetUUID, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(targetUUID)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, res)
	if res == nil {
		return
	}
	assert.Equal(t, targetUUID, res.ID)
}

func TestRemarkUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockUsecase(t)

	req := model.CreateRemarkRequest{
		EntityType: "yayasan",
		Name:       "Ibu Ketua",
		Position:   "Ketua Yayasan",
		Content:    "Halo",
		OrderIndex: 1,
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

	rows := sqlmock.NewRows([]string{"id", "entity_type", "name", "position", "order_index"}).
		AddRow("uuid-new", "yayasan", "Ibu Ketua", "Ketua Yayasan", 1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks`")).
		WithArgs("uuid-new", 1).
		WillReturnRows(rows)

	res, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, res)
	if res == nil {
		return
	}
	assert.Equal(t, "Ibu Ketua", res.Name)
	assert.NotEmpty(t, res.ID)
}

func TestRemarkUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockUsecase(t)
	targetID := "uuid-delete-me"

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Deleted Guy")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `remarks` WHERE id = ?")).
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
