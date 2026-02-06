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

func setupMockActivityUsecase(t *testing.T) (usecase.ActivityUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewActivityUsecase(gormDB, validator.New())
	return u, mock
}

func TestActivityUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)

	req := model.CreateActivityRequest{
		EntityType:  "pura",
		Title:       "Upacara",
		Description: "Deskripsi",
		TimeInfo:    "08:00",
		Location:    "Pura",
		EventDate:   "2023-10-10",
		OrderIndex:  1,
		IsActive:    true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `activities`")).
		WithArgs(
			sqlmock.AnyArg(),
			req.EntityType,
			"Upacara",
			"Deskripsi",
			"08:00",
			"Pura",
			sqlmock.AnyArg(),
			1,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Create(req.EntityType, req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, req.Title, res.Title)
		assert.Equal(t, req.Description, res.Description)
	}
}

func TestActivityUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockActivityUsecase(t)
	req := model.CreateActivityRequest{}
	res, err := u.Create("pura", req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestActivityUsecase_Create_BadRequest_Date(t *testing.T) {
	u, _ := setupMockActivityUsecase(t)

	req := model.CreateActivityRequest{
		EntityType:  "pura",
		Title:       "Title",
		Description: "Desc",
		EventDate:   "invalid-date",
		OrderIndex:  1,
	}
	res, err := u.Create("pura", req)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestActivityUsecase_GetAll_OrderedByIndex(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "order_index", "is_active"}).
		AddRow("a1", "A", "d", 1, true).
		AddRow("a2", "B", "d", 2, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE entity_type = ? ORDER BY event_date DESC,order_index ASC")).
		WithArgs("pura").
		WillReturnRows(rows)

	list, err := u.GetAll("pura")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "A", list[0].Title)
	assert.Equal(t, "B", list[1].Title)
}

func TestActivityUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "order_index", "is_active"}).
		AddRow("a3", "A", "d", 1, true).
		AddRow("a1", "B", "d", 2, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE entity_type = ? AND is_active = ? ORDER BY event_date DESC,order_index ASC")).
		WithArgs("pura", true).
		WillReturnRows(rows)

	list, err := u.GetPublic("pura")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "A", list[0].Title)
	assert.Equal(t, "B", list[1].Title)
}

func TestActivityUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)
	id := "act-1"

	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(id, "Title")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Title", res.Title)
}

func TestActivityUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE id = ? LIMIT ?")).
		WithArgs("missing", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID("missing")
	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "activity not found", e.Message)
	}
}

func TestActivityUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)
	targetID := "act-1"

	req := model.UpdateActivityRequest{
		Title:       "New",
		Description: "new d",
		TimeInfo:    "09:00",
		Location:    "Pura",
		EventDate:   "2023-12-12",
		OrderIndex:  5,
		IsActive:    false,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "title"}).AddRow(targetID, "pura", "Old"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `activities`")).
		WithArgs(
			"pura",
			"New",
			"new d",
			"09:00",
			"Pura",
			sqlmock.AnyArg(),
			5,
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "New", res.Title)
		assert.Equal(t, false, res.IsActive)
		assert.Equal(t, 5, res.OrderIndex)
	}
}

func TestActivityUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)
	targetID := "missing"

	req := model.UpdateActivityRequest{
		Title:       "New",
		Description: "Desc",
		EventDate:   "2023-01-01",
		OrderIndex:  1,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.Update(targetID, req)
	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "activity not found", e.Message)
	}
}

func TestActivityUsecase_Update_BadRequest_Date(t *testing.T) {
	u, _ := setupMockActivityUsecase(t)
	targetID := "act-1"

	req := model.UpdateActivityRequest{
		Title:       "New",
		Description: "Desc",
		EventDate:   "invalid",
		OrderIndex:  1,
	}

	res, err := u.Update(targetID, req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestActivityUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)
	targetID := "to-del"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(targetID, "Del"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `activities` WHERE `activities`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestActivityUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockActivityUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activities` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)
	assert.Error(t, err)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "activity not found", e.Message)
	}
}
