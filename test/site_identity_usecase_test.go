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

	u := usecase.NewSiteIdentityUsecase(gormDB, validator.New())
	return u, mock
}

func TestSiteIdentityUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	req := model.SiteIdentityRequest{
		EntityType: "pura",
		SiteName:   "Pura Agung",
		LogoURL:    "https://logo.com/img.png",
		Tagline:    "Tagline",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `site_identity`")).
		WithArgs(
			sqlmock.AnyArg(),
			"pura",
			"Pura Agung",
			"https://logo.com/img.png",
			"Tagline",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Pura Agung", res.SiteName)
}

func TestSiteIdentityUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockSiteIdentityUsecase(t)
	req := model.SiteIdentityRequest{}
	res, err := u.Create("pura", req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSiteIdentityUsecase_GetAll(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "entity_type", "site_name"}).
		AddRow("s1", "pura", "A").
		AddRow("s2", "pura", "B")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE entity_type = ? ORDER BY created_at ASC")).
		WithArgs("pura").
		WillReturnRows(rows)

	list, err := u.GetAll("pura")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestSiteIdentityUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "sid-123"

	rows := sqlmock.NewRows([]string{"id", "site_name"}).AddRow(targetID, "My Site")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(targetID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestSiteIdentityUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ? LIMIT ?")).
		WithArgs("missing", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID("missing")

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "site identity not found", e.Message)
	}
}

func TestSiteIdentityUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "sid-1"

	req := model.SiteIdentityRequest{
		EntityType:        "yayasan",
		SiteName:          "New Site",
		Tagline:           "New Tag",
		PrimaryButtonText: "Go",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "site_name", "tagline"}).AddRow(targetID, "pura", "Old", "Old"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `site_identity`")).
		WithArgs(
			"yayasan",
			"New Site",
			sqlmock.AnyArg(),
			"New Tag",
			"Go",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "New Site", res.SiteName)
		assert.Equal(t, "Go", res.PrimaryButtonText)
	}
}

func TestSiteIdentityUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "missing"

	req := model.SiteIdentityRequest{
		EntityType: "pura",
		SiteName:   "Valid Name",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "site identity not found", e.Message)
	}
}

func TestSiteIdentityUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "del-1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ? LIMIT ?")).
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

func TestSiteIdentityUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "site identity not found", e.Message)
	}
}

func TestSiteIdentityUsecase_GetPublic_ReturnsLatest(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "entity_type", "site_name"}).
		AddRow("new", "pura", "New")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE entity_type = ? ORDER BY created_at DESC LIMIT ?")).
		WithArgs("pura", 1).
		WillReturnRows(rows)

	result, err := u.GetPublic("pura")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New", result.SiteName)
}

func TestSiteIdentityUsecase_GetPublic_NotFound(t *testing.T) {
	u, mock := setupMockSiteIdentityUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `site_identity` WHERE entity_type = ? ORDER BY created_at DESC LIMIT ?")).
		WithArgs("pura", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	result, err := u.GetPublic("pura")

	assert.Error(t, err)
	assert.Nil(t, result)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "site identity not found", e.Message)
	}
}
