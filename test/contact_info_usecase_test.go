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

	u := usecase.NewContactInfoUsecase(gormDB, logrus.New(), validator.New())
	return u, mock
}

func TestContactInfoUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)

	req := model.ContactInfoRequest{
		Address:       "Jl. Contoh No.1",
		Phone:         "+62 8123456789",
		Email:         "info@example.com",
		VisitingHours: "08:00 - 17:00",
		MapEmbedURL:   "https://maps.google.com/?q=x",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `contact_info`")).
		WithArgs(sqlmock.AnyArg(), "Jl. Contoh No.1", "+62 8123456789", "info@example.com", "08:00 - 17:00", "https://maps.google.com/?q=x", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "address", "phone", "email", "visiting_hours", "map_embed_url"}).
		AddRow("ci-1", "Jl. Contoh No.1", "+62 8123456789", "info@example.com", "08:00 - 17:00", "https://maps.google.com/?q=x")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info`")).
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
	assert.Equal(t, req.Address, res.Address)
	assert.Equal(t, req.Email, res.Email)
}

func TestContactInfoUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockContactInfoUsecase(t)

	req := model.ContactInfoRequest{}

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContactInfoUsecase_GetAll(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "address", "email"}).
		AddRow("1", "A", "email1").
		AddRow("2", "B", "email2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info`")).
		WillReturnRows(rows)

	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestContactInfoUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ?")).
		WithArgs("not-exists", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID("not-exists")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContactInfoUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	targetID := "ci-1"

	req := model.ContactInfoRequest{Address: "New Addr", Email: "new@example.com", Phone: "000"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "email"}).AddRow(targetID, "Old Addr", "old@example.com"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `contact_info`")).
		WithArgs("New Addr", "000", "new@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
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
	assert.Equal(t, "New Addr", res.Address)
	assert.Equal(t, "new@example.com", res.Email)
}

func TestContactInfoUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockContactInfoUsecase(t)
	targetID := "to-del"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contact_info` WHERE id = ?")).
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
