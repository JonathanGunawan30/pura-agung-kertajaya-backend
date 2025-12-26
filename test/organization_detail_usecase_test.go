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

func setupMockOrgDetailUsecase(t *testing.T) (usecase.OrganizationDetailUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewOrganizationDetailUsecase(gormDB, logger, validate)

	return u, mock
}

func TestOrgDetailUsecase_GetByEntityType_Success(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "entity_type", "vision", "mission", "created_at", "updated_at"}).
		AddRow("uuid-1", "pura", "Visi Pura", "Misi Pura", now, now)

	expectedSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs("pura", 1).
		WillReturnRows(rows)

	result, err := u.GetByEntityType("pura")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Visi Pura", result.Vision)
	assert.Equal(t, "pura", result.EntityType)
}

func TestOrgDetailUsecase_GetByEntityType_NotFound(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	expectedSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs("pasraman", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := u.GetByEntityType("pasraman")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestOrgDetailUsecase_Update_TriggerCreate(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	req := model.UpdateOrganizationDetailRequest{
		Vision:  "Visi Baru",
		Mission: "Misi Baru",
		Rules:   "Aturan",
	}
	entityType := "yayasan"

	selectSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"
	mock.ExpectQuery(regexp.QuoteMeta(selectSQL)).
		WithArgs(entityType, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `organization_details`")).
		WithArgs(
			sqlmock.AnyArg(),
			entityType,
			"Visi Baru",
			"Misi Baru",
			"Aturan",
			"",
			"",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	res, err := u.Update(entityType, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Visi Baru", res.Vision)
	assert.Equal(t, entityType, res.EntityType)
	assert.NotEmpty(t, res.ID)
}

func TestOrgDetailUsecase_Update_TriggerUpdate(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	entityType := "pura"
	existingID := "existing-uuid"

	req := model.UpdateOrganizationDetailRequest{
		Vision:  "Visi Update",
		Mission: "Misi Lama",
	}

	selectSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"
	rows := sqlmock.NewRows([]string{"id", "entity_type", "vision", "mission"}).
		AddRow(existingID, entityType, "Visi Lama", "Misi Lama")

	mock.ExpectQuery(regexp.QuoteMeta(selectSQL)).
		WithArgs(entityType, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `organization_details` SET")).
		WithArgs(
			entityType,
			"Visi Update",
			"Misi Lama",
			"",
			"",
			"",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			existingID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	res, err := u.Update(entityType, req)

	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.NotNil(t, res)
	assert.Equal(t, "Visi Update", res.Vision)
	assert.Equal(t, existingID, res.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
