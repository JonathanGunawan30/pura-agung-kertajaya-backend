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

	rows := sqlmock.NewRows([]string{"id", "name", "images", "is_active", "order_index"}).
		AddRow("g3", "C", []byte(`{"lg":"https://img3"}`), true, 1).
		AddRow("g1", "B", []byte(`{"lg":"https://img1"}`), true, 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE entity_type = ? AND is_active = ? ORDER BY order_index ASC")).
		WithArgs("pura", true).
		WillReturnRows(rows)

	list, err := u.GetPublic("pura")

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "g3", list[0].ID)
	assert.Equal(t, "g1", list[1].ID)
	assert.Equal(t, "https://img3", list[0].Images.Lg)
}

func TestFacilityUsecase_Create(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)

	req := model.CreateFacilityRequest{
		EntityType: "pura",
		Name:       "Name",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
		IsActive:   true,
		OrderIndex: 3,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `facilities`")).
		WithArgs(
			sqlmock.AnyArg(),
			req.EntityType,
			"Name",
			"",
			sqlmock.AnyArg(),
			3,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "entity_type", "name", "description", "images", "order_index", "is_active"}).
		AddRow("uuid-1", "pura", "Name", "", []byte(`{"lg":"https://img.com/lg.jpg"}`), 3, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities`")).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	created, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "https://img.com/lg.jpg", created.Images.Lg)
}

func TestFacilityUsecase_Update(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "g1"

	req := model.UpdateFacilityRequest{
		Name:       "New Name",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
		IsActive:   false,
		OrderIndex: 5,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "images"}).
			AddRow(targetID, "Old Name", []byte(`{"lg":"old.jpg"}`)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `facilities`")).
		WithArgs(
			sqlmock.AnyArg(),
			"New Name",
			"",
			sqlmock.AnyArg(),
			5,
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, updated)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, false, updated.IsActive)
	assert.Equal(t, 5, updated.OrderIndex)
}

func TestFacilityUsecase_Delete(t *testing.T) {
	u, mock := setupMockFacilityUsecase(t)
	targetID := "g1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `facilities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "images"}).
			AddRow(targetID, "Name", []byte(`{}`)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `facilities` WHERE `facilities`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}
