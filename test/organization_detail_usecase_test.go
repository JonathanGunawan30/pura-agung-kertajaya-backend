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
	rows := sqlmock.NewRows([]string{
		"id", "entity_type", "vision", "mission", "rules", "work_program",
		"vision_mission_image_url", "work_program_image_url", "rules_image_url",
		"created_at", "updated_at",
	}).
		AddRow("uuid-1", "pura", "Visi Pura", "Misi Pura", "Aturan", "Proker",
			"img_vm.jpg", "img_wp.jpg", "img_r.jpg", now, now)

	expectedSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs("pura", 1).
		WillReturnRows(rows)

	result, err := u.GetByEntityType("pura")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Visi Pura", result.Vision)
	assert.Equal(t, "img_vm.jpg", result.VisionMissionImageURL)
	assert.Equal(t, "pura", result.EntityType)
}

func TestOrgDetailUsecase_GetByEntityType_NotFound(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	expectedSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs("pasraman", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := u.GetByEntityType("pasraman")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "pasraman", result.EntityType)
	assert.Equal(t, "", result.Vision)
}

func TestOrgDetailUsecase_Update_TriggerCreate(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	req := model.UpdateOrganizationDetailRequest{
		Vision:                "Visi Baru",
		Mission:               "Misi Baru",
		Rules:                 "Aturan",
		VisionMissionImageURL: "new_img.jpg",
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
			"new_img.jpg",
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
	assert.Equal(t, "new_img.jpg", res.VisionMissionImageURL)
}

func TestOrgDetailUsecase_Update_TriggerUpdate(t *testing.T) {
	u, mock := setupMockOrgDetailUsecase(t)

	entityType := "pura"
	existingID := "existing-uuid"
	now := time.Now()

	req := model.UpdateOrganizationDetailRequest{
		Vision:              "Visi Update",
		WorkProgramImageURL: "wp_update.jpg",
	}

	selectSQL := "SELECT * FROM `organization_details` WHERE entity_type = ? ORDER BY `organization_details`.`id` LIMIT ?"

	rows := sqlmock.NewRows([]string{
		"id", "entity_type", "vision", "mission", "rules", "work_program",
		"vision_mission_image_url", "work_program_image_url", "rules_image_url",
		"created_at", "updated_at",
	}).AddRow(existingID, entityType, "Visi Lama", "Misi Lama", "", "", "", "old_wp.jpg", "", now, now)

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
			"wp_update.jpg",
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
	assert.Equal(t, "wp_update.jpg", res.WorkProgramImageURL)

	assert.NoError(t, mock.ExpectationsWereMet())
}
