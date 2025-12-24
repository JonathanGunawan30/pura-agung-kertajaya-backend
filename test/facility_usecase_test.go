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

	u := usecase.NewFacilityUsecase(gormDB, logrus.New(), validator.New())
	return u, mock
}

func TestFacilityUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "image_url", "is_active", "order_index"}).
		AddRow("g3", "C", "https://img3", true, 1).
		AddRow("g1", "B", "https://img1", true, 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE is_active = ? ORDER BY order_index ASC")).
		WithArgs(true).
		WillReturnRows(rows)

	list, err := u.GetPublic()

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "g3", list[0].ID)
	assert.Equal(t, "g1", list[1].ID)
}

func TestFacilityUsecase_Create(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)

	req := model.FacilityRequest{Name: "Name", ImageURL: "https://img", IsActive: true, OrderIndex: 3}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `facilities`")).
		WithArgs(
			sqlmock.AnyArg(),
			"Name",
			"",
			"https://img",
			3,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "image_url", "order_index", "is_active"}).
		AddRow("uuid-1", "Name", "", "https://img", 3, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	created, err := u.Create(req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, created)
	if created == nil {
		return
	}
	assert.NotEmpty(t, created.ID)
}

func TestFacilityUsecase_Update(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "g1"

	req := model.FacilityRequest{Name: "New Name", ImageURL: "https://img2", IsActive: false, OrderIndex: 5}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Old Name"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `facilities`")).
		WithArgs("New Name", "", "https://img2", 5, false, sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, updated)
	if updated == nil {
		return
	}
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, false, updated.IsActive)
	assert.Equal(t, 5, updated.OrderIndex)
}

func TestFacilityUsecase_Delete(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "g1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Name"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `facilities` WHERE `facilities`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}
