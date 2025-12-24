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

func setupMockAboutUsecase(t *testing.T) (usecase.AboutUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewAboutUsecase(gormDB, logrus.New(), validator.New())
	return u, mock
}

func TestAboutUsecase_Create_WithValues(t *testing.T) {
	u, mock := setupMockAboutUsecase(t)

	req := model.AboutSectionRequest{
		Title:       "About Title",
		Description: "About Description",
		ImageURL:    "https://img",
		IsActive:    true,
		Values: []model.AboutValueRequest{
			{Title: "Vision", Value: "Be great", OrderIndex: 2},
			{Title: "Mission", Value: "Serve", OrderIndex: 1},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `about_section`")).
		WithArgs(sqlmock.AnyArg(), "About Title", "About Description", "https://img", true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `about_values`")).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 2))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "title", "description", "image_url", "is_active"}).
		AddRow("60d57d69-968b-46fb-ab59-4e84f99aa40b", "About Title", "About Description", "https://img", true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_section`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	rowsValues := sqlmock.NewRows([]string{"id", "about_id", "title", "value", "order_index"}).
		AddRow("val-1", "60d57d69-968b-46fb-ab59-4e84f99aa40b", "Vision", "Be great", 2).
		AddRow("val-2", "60d57d69-968b-46fb-ab59-4e84f99aa40b", "Mission", "Serve", 1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_values`")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rowsValues)

	res, err := u.Create(req)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, res)
	if res == nil {
		return
	}
	assert.Equal(t, 2, len(res.Values))
}

func TestAboutUsecase_GetPublic_FilterActive(t *testing.T) {
	u, mock := setupMockAboutUsecase(t)

	rowsSection := sqlmock.NewRows([]string{"id", "title", "is_active"}).
		AddRow("uuid-1", "Active", true)

	rowsValues := sqlmock.NewRows([]string{"id", "about_id", "title", "value", "order_index"}).
		AddRow("val-1", "uuid-1", "A", "val-a", 1).
		AddRow("val-2", "uuid-1", "B", "val-b", 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_section` WHERE is_active = ?")).
		WithArgs(true).
		WillReturnRows(rowsSection)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_values` WHERE `about_values`.`about_id` = ?")).
		WithArgs("uuid-1").
		WillReturnRows(rowsValues)

	list, err := u.GetPublic()

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Active", list[0].Title)
	assert.Equal(t, "A", list[0].Values[0].Title)
}

func TestAboutUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockAboutUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_section` WHERE id = ?")).
		WithArgs("missing", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID("missing")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestAboutUsecase_Update_ReplacesValues(t *testing.T) {
	u, mock := setupMockAboutUsecase(t)
	targetID := "ab-1"

	req := model.AboutSectionRequest{
		Title:       "New",
		Description: "nd",
		ImageURL:    "https://new",
		IsActive:    false,
		Values: []model.AboutValueRequest{
			{Title: "New1", Value: "n1", OrderIndex: 1},
		},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_section` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(targetID, "Old"))

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `about_section`")).
		WithArgs("New", "nd", "https://new", false, sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `about_values` WHERE about_id = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `about_values`")).
		WithArgs(sqlmock.AnyArg(), targetID, "New1", "n1", 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	rowsResult := sqlmock.NewRows([]string{"id", "title", "description", "image_url", "is_active"}).
		AddRow(targetID, "New", "nd", "https://new", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_section`")).
		WithArgs(targetID, targetID, 1).
		WillReturnRows(rowsResult)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_values`")).
		WithArgs(targetID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "about_id", "title"}).AddRow("v1", targetID, "New1"))

	res, err := u.Update(targetID, req)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, res)
	if res == nil {
		return
	}
	assert.Equal(t, "New", res.Title)
}

func TestAboutUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockAboutUsecase(t)
	targetID := "to-del"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `about_section` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `about_section` WHERE `about_section`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}
