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

	u := usecase.NewOrganizationRequest(gormDB, logrus.New(), validator.New())
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
		WithArgs(sqlmock.AnyArg(), req.EntityType, "Ketut Test", "Bendahara", 3, 1, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"id", "name", "position", "position_order", "order_index", "is_active"}).
		AddRow("uuid-1", "Ketut Test", "Bendahara", 3, 1, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members`")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	created, err := u.Create(req.EntityType, req)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, created)
	if created == nil {
		return
	}
	assert.NotEmpty(t, created.ID)
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

func TestOrganizationMemberUsecase_GetAll_OrderedCorrectly(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "position", "position_order", "order_index", "is_active"}).
		AddRow("m1", "Ketua", "Ketua", 1, 1, true).
		AddRow("m2", "Wakil", "Wakil Ketua", 2, 1, true).
		AddRow("m4", "Sekre 1", "Sekretaris", 3, 1, true).
		AddRow("m3", "Sekre 2", "Sekretaris", 3, 2, true).
		AddRow("m5", "NonAktif", "Anggota", 4, 1, false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE entity_type = ? ORDER BY position_order ASC, order_index ASC")).
		WithArgs("pura").
		WillReturnRows(rows)

	list, err := u.GetAll("pura")

	assert.NoError(t, err)
	assert.Len(t, list, 5)
	assert.Equal(t, "m1", list[0].ID)
	assert.Equal(t, "m2", list[1].ID)
	assert.Equal(t, "m4", list[2].ID)
	assert.Equal(t, "m3", list[3].ID)
	assert.Equal(t, "m5", list[4].ID)
}

func TestOrganizationMemberUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "name", "position", "position_order", "order_index", "is_active"}).
		AddRow("m1", "Ketua", "Ketua", 1, 1, true).
		AddRow("m2", "Wakil", "Wakil Ketua", 2, 1, true).
		AddRow("m4", "Sekre 1", "Sekretaris", 3, 1, true).
		AddRow("m3", "Sekre 2", "Sekretaris", 3, 2, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE entity_type = ? AND is_active = ? ORDER BY order_index ASC")).
		WithArgs("pura", true).
		WillReturnRows(rows)

	list, err := u.GetPublic("pura")

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Len(t, list, 4)
	if len(list) == 4 {
		assert.Equal(t, "m1", list[0].ID)
		assert.Equal(t, "m2", list[1].ID)
		assert.Equal(t, "m4", list[2].ID)
		assert.Equal(t, "m3", list[3].ID)
	}
}

func TestOrganizationMemberUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockOrganizationUsecase(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `organization_members` WHERE id = ? LIMIT ?")).
		WithArgs("non-existent-id", 1).
		WillReturnRows(sqlmock.NewRows(nil))

	found, err := u.GetByID("non-existent-id")

	assert.Error(t, err)
	assert.Nil(t, found)
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(targetID, "Old Name"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `organization_members`")).
		WithArgs(sqlmock.AnyArg(), "New Name", "New Pos", 5, 2, false, sqlmock.AnyArg(), sqlmock.AnyArg(), targetID).
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
