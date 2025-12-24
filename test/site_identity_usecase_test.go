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

func setupMockSiteIdentityUsecase(t *testing.T) (usecase.SiteIdentityUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewSiteIdentityUsecase(gormDB, logrus.New(), validator.New())
	return u, mock
}

func TestSiteIdentityUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	req := model.SiteIdentityRequest{SiteName: "Pura", LogoURL: "https://logo", Tagline: "Tag"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `site_identity`")).
		WithArgs(sqlmock.AnyArg(), "Pura", "https://logo", "Tag", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "site_name", "logo_url", "tagline"}).
		AddRow("uuid-1", "Pura", "https://logo", "Tag")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity`")).
		WithArgs(sqlmock.AnyArg(), 1).
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
	assert.Equal(t, "Pura", res.SiteName)
}

func TestSiteIdentityUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockSiteIdentityUsecase(t)

	req := model.SiteIdentityRequest{}

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSiteIdentityUsecase_GetAll(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "site_name"}).
		AddRow("s1", "A").
		AddRow("s2", "B")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity`")).
		WillReturnRows(rows)

	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestSiteIdentityUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ?")).
		WithArgs("missing", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID("missing")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSiteIdentityUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "sid-1"

	req := model.SiteIdentityRequest{SiteName: "New", Tagline: "New", PrimaryButtonText: "Go"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "site_name", "tagline"}).AddRow(targetID, "Old", "Old"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `site_identity`")).
		WithArgs("New", sqlmock.AnyArg(), "New", "Go", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
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
	assert.Equal(t, "New", res.SiteName)
	assert.Equal(t, "Go", res.PrimaryButtonText)
}

func TestSiteIdentityUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "del-1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "site_name"}).AddRow(targetID, "Name"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `site_identity` WHERE `site_identity`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestSiteIdentityUsecase_GetPublic_ReturnsLatest(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "site_name"}).
		AddRow("new", "New")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` ORDER BY created_at DESC LIMIT ?")).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := u.GetPublic()

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, result)
	if result == nil {
		return
	}
	assert.Equal(t, "New", result.SiteName)
}
