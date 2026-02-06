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

func setupMockContactInfoUsecase(t *testing.T) (usecase.ContactInfoUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewContactInfoUsecase(gormDB, validator.New())
	return u, mock
}

func TestContactInfoUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)

	req := model.CreateContactInfoRequest{
		EntityType:    "pura",
		Address:       "Jl. Contoh No.1",
		Phone:         "+62 8123456789",
		Email:         "info@example.com",
		VisitingHours: "08:00 - 17:00",
		MapEmbedURL:   "https://maps.google.com/?q=x",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `contact_info`")).
		WithArgs(sqlmock.AnyArg(), req.EntityType, "Jl. Contoh No.1", "+62 8123456789", "info@example.com", "08:00 - 17:00", "https://maps.google.com/?q=x", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Create(req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, req.Address, res.Address)
		assert.Equal(t, req.Email, res.Email)
	}
}

func TestContactInfoUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockContactInfoUsecase(t)

	req := model.CreateContactInfoRequest{}

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContactInfoUsecase_GetAll(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "address", "email", "created_at"}).
		AddRow("1", "A", "email1", time.Now()).
		AddRow("2", "B", "email2", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE entity_type = ? ORDER BY created_at ASC")).
		WithArgs("pura").
		WillReturnRows(rows)

	list, err := u.GetAll("pura")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestContactInfoUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	id := "ci-1"

	rows := sqlmock.NewRows([]string{"id", "address", "email"}).
		AddRow(id, "Addr", "email@test.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(id)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Addr", res.Address)
	}
}

func TestContactInfoUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	id := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "contact info not found", e.Message)
	}
}

func TestContactInfoUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	targetID := "ci-1"

	req := model.UpdateContactInfoRequest{
		Address:       "New Addr",
		Email:         "new@example.com",
		Phone:         "000",
		VisitingHours: "09:00 - 15:00",
		MapEmbedURL:   "http://maps.com/new",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "address", "email"}).
			AddRow(targetID, "pura", "Old Addr", "old@example.com"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `contact_info`")).
		WithArgs(
			"pura",
			"New Addr",
			"000",
			"new@example.com",
			"09:00 - 15:00",
			"http://maps.com/new",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)

	if assert.NotNil(t, res) {
		assert.Equal(t, "New Addr", res.Address)
		assert.Equal(t, "new@example.com", res.Email)
	}
}

func TestContactInfoUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	targetID := "missing"

	req := model.UpdateContactInfoRequest{
		Address:       "New Addr",
		Email:         "new@example.com",
		Phone:         "000",
		VisitingHours: "09:00 - 15:00",
		MapEmbedURL:   "http://maps.com/new",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "contact info not found", e.Message)
	}
}

func TestContactInfoUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	targetID := "to-del"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address"}).AddRow(targetID, "Addr"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `contact_info` WHERE `contact_info`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestContactInfoUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "contact info not found", e.Message)
	}
}
