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

func setupMockOrganizationUsecase(t *testing.T) (usecase.OrganizationUsecase, sqlmock.Sqlmock) {
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

	u := usecase.NewOrganizationUsecase(gormDB, validator.New())
	return u, mock
}

func TestOrganizationMemberUsecase_Create_Success(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)

	req := model.CreateOrganizationRequest{
		EntityType:    "pura",
		Name:          "Ketut Test",
		Position:      "Bendahara",
		PositionOrder: 3,
		OrderIndex:    1,
		IsActive:      true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `organization_members`")).
		WithArgs(
			sqlmock.AnyArg(),
			req.EntityType,
			"Ketut Test",
			"Bendahara",
			3,
			1,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req.EntityType, req)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, req.Name, created.Name)
}

func TestOrganizationMemberUsecase_Create_ValidationError(t *testing.T) {
	u, _ := setupMockOrganizationUsecase(t)

	req := model.CreateOrganizationRequest{
		Name:          "",
		Position:      "Tester",
		PositionOrder: 1,
		IsActive:      true,
	}

	created, err := u.Create("pura", req)

	assert.Error(t, err)
	assert.Nil(t, created)
}

func TestOrganizationMemberUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "position", "position_order", "order_index", "is_active"}).
		AddRow("m1", "Ketua", "Ketua", 1, 1, true).
		AddRow("m2", "Wakil", "Wakil Ketua", 2, 1, true).
		AddRow("m4", "Sekre 1", "Sekretaris", 3, 1, true).
		AddRow("m3", "Sekre 2", "Sekretaris", 3, 2, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE entity_type = ? AND is_active = ? ORDER BY position_order ASC, order_index ASC")).
		WithArgs("pura", true).
		WillReturnRows(rows)

	list, err := u.GetPublic("pura")

	assert.NoError(t, err)
	assert.Len(t, list, 4)
	assert.Equal(t, "m1", list[0].ID)
}

func TestOrganizationMemberUsecase_GetByID_Success(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)
	id := "uuid-1"

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Member Name")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := u.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestOrganizationMemberUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)
	id := "non-existent-id"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	found, err := u.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, found)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "organization member not found", e.Message)
	}
}

func TestOrganizationMemberUsecase_Update_Success(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)
	targetID := "update-me"

	req := model.UpdateOrganizationRequest{
		Name:          "New Name",
		Position:      "New Pos",
		PositionOrder: 5,
		OrderIndex:    2,
		IsActive:      false,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "name"}).AddRow(targetID, "pura", "Old Name"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `organization_members`")).
		WithArgs(
			"pura",
			"New Name",
			"New Pos",
			5,
			2,
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			targetID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(targetID, req)
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, false, updated.IsActive)
	}
}

func TestOrganizationMemberUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)
	targetID := "missing"

	req := model.UpdateOrganizationRequest{
		Name:          "Valid Name",
		Position:      "Valid Pos",
		PositionOrder: 1,
		OrderIndex:    1,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	updated, err := u.Update(targetID, req)

	assert.Error(t, err)
	assert.Nil(t, updated)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "organization member not found", e.Message)
	}
}

func TestOrganizationMemberUsecase_Delete_Success(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)
	targetID := "delete-me"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "To Delete"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `organization_members` WHERE `organization_members`.`id` = ?")).
		WithArgs(targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(targetID)
	assert.NoError(t, err)
}

func TestOrganizationMemberUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)
	targetID := "missing"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs(targetID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	err := u.Delete(targetID)

	assert.Error(t, err)
	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "organization member not found", e.Message)
	}
}
