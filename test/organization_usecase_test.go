package test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

// Helper function
func boolPtr(b bool) *bool {
	return &b
}

func newOrganizationMemberUsecase() usecase.OrganizationUsecase {
	return usecase.NewOrganizationRequest(db, logrus.New(), validator.New()) // Assuming constructor exists
}

// ClearOrganizationMembers clears the table before each test
func ClearOrganizationMembers() {
	if err := db.Where("id IS NOT NULL").Delete(&entity.OrganizationMember{}).Error; err != nil {
		log.Fatalf("Failed clear organization_members: %+v", err)
	}
}

func TestOrganizationMemberUsecase_Create_Success(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	req := model.OrganizationRequest{
		Name:          "Ketut Test",
		Position:      "Bendahara",
		PositionOrder: 3,
		OrderIndex:    1,
		IsActive:      true,
	}

	created, err := u.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, req.Name, created.Name)
	assert.Equal(t, req.Position, created.Position)
	assert.Equal(t, req.PositionOrder, created.PositionOrder)
	assert.Equal(t, req.IsActive, created.IsActive)

	// Verify in DB
	var count int64
	db.Model(&entity.OrganizationMember{}).Where("id = ?", created.ID).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestOrganizationMemberUsecase_Create_ValidationError(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	req := model.OrganizationRequest{
		Name:          "", // Invalid: Name is required
		Position:      "Tester",
		PositionOrder: 1,
		IsActive:      true,
	}

	created, err := u.Create(req)

	assert.Error(t, err) // Expect validation error
	assert.Nil(t, created)
}

func TestOrganizationMemberUsecase_GetAll_OrderedCorrectly(t *testing.T) {
	ClearOrganizationMembers()

	// Seed data with different orders
	m1 := entity.OrganizationMember{ID: "m1", Name: "Ketua", Position: "Ketua", PositionOrder: 1, OrderIndex: 1, IsActive: true}
	m2 := entity.OrganizationMember{ID: "m2", Name: "Wakil", Position: "Wakil Ketua", PositionOrder: 2, OrderIndex: 1, IsActive: true}
	m3 := entity.OrganizationMember{ID: "m3", Name: "Sekre 2", Position: "Sekretaris", PositionOrder: 3, OrderIndex: 2, IsActive: true} // Higher order index
	m4 := entity.OrganizationMember{ID: "m4", Name: "Sekre 1", Position: "Sekretaris", PositionOrder: 3, OrderIndex: 1, IsActive: true} // Lower order index, same position order
	m5 := entity.OrganizationMember{ID: "m5", Name: "NonAktif", Position: "Anggota", PositionOrder: 4, OrderIndex: 1, IsActive: false}  // Should not appear in GetPublic

	db.Create(&m1)
	db.Create(&m2)
	db.Create(&m3)
	db.Create(&m4)
	db.Create(&m5) // Order matters if timestamps are too close

	u := newOrganizationMemberUsecase()
	list, err := u.GetAll() // Test GetAll (includes inactive)

	assert.NoError(t, err)
	assert.Len(t, list, 5) // Should fetch all 5

	// Check order: position_order ASC, then order_index ASC
	assert.Equal(t, "m1", list[0].ID) // Ketua (Pos 1)
	assert.Equal(t, "m2", list[1].ID) // Wakil (Pos 2)
	assert.Equal(t, "m4", list[2].ID) // Sekre 1 (Pos 3, Index 1)
	assert.Equal(t, "m3", list[3].ID) // Sekre 2 (Pos 3, Index 2)
	assert.Equal(t, "m5", list[4].ID) // NonAktif (Pos 4)
}

func TestOrganizationMemberUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	ClearOrganizationMembers()

	// Seed data similar to GetAll test
	m1 := entity.OrganizationMember{ID: "m1", Name: "Ketua", Position: "Ketua", PositionOrder: 1, OrderIndex: 1, IsActive: true}
	m2 := entity.OrganizationMember{ID: "m2", Name: "Wakil", Position: "Wakil Ketua", PositionOrder: 2, OrderIndex: 1, IsActive: true}
	m3 := entity.OrganizationMember{ID: "m3", Name: "Sekre 2", Position: "Sekretaris", PositionOrder: 3, OrderIndex: 2, IsActive: true}
	m4 := entity.OrganizationMember{ID: "m4", Name: "Sekre 1", Position: "Sekretaris", PositionOrder: 3, OrderIndex: 1, IsActive: true}
	m5 := entity.OrganizationMember{ID: "m5", Name: "NonAktif", Position: "Anggota", PositionOrder: 4, OrderIndex: 1, IsActive: false} // Inactive

	db.Create(&m1)
	db.Create(&m2)
	db.Create(&m3)
	db.Create(&m4)
	db.Create(&m5)

	u := newOrganizationMemberUsecase()
	list, err := u.GetPublic() // Test GetPublic (active only)

	assert.NoError(t, err)
	assert.Len(t, list, 4) // Should only fetch 4 active members

	// Check order: position_order ASC, then order_index ASC
	assert.Equal(t, "m1", list[0].ID) // Ketua (Pos 1)
	assert.Equal(t, "m2", list[1].ID) // Wakil (Pos 2)
	assert.Equal(t, "m4", list[2].ID) // Sekre 1 (Pos 3, Index 1)
	assert.Equal(t, "m3", list[3].ID) // Sekre 2 (Pos 3, Index 2)
}

func TestOrganizationMemberUsecase_GetByID_NotFound(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	found, err := u.GetByID("non-existent-id")

	assert.Error(t, err) // Expect an error (e.g., gorm.ErrRecordNotFound)
	assert.Nil(t, found)
}

func TestOrganizationMemberUsecase_Update_Success(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	// Create initial member
	initial := entity.OrganizationMember{ID: "update-me", Name: "Old Name", Position: "Old Pos", PositionOrder: 10, IsActive: true}
	db.Create(&initial)

	req := model.OrganizationRequest{
		Name:          "New Name",
		Position:      "New Pos",
		PositionOrder: 5,
		OrderIndex:    2,
		IsActive:      false,
	}

	updated, err := u.Update("update-me", req)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "update-me", updated.ID)
	assert.Equal(t, req.Name, updated.Name)
	assert.Equal(t, req.Position, updated.Position)
	assert.Equal(t, req.PositionOrder, updated.PositionOrder)
	assert.Equal(t, req.OrderIndex, updated.OrderIndex)
	assert.Equal(t, req.IsActive, updated.IsActive)

	// Verify changes in DB
	var final entity.OrganizationMember
	db.First(&final, "id = ?", "update-me")
	assert.Equal(t, "New Name", final.Name)
	assert.Equal(t, false, final.IsActive) // Check pointer value
}

func TestOrganizationMemberUsecase_Update_NotFound(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	req := model.OrganizationRequest{Name: "Update"}
	updated, err := u.Update("non-existent-id", req)

	assert.Error(t, err)
	assert.Nil(t, updated)
}

func TestOrganizationMemberUsecase_Delete_Success(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	// Create member to delete
	member := entity.OrganizationMember{ID: "delete-me", Name: "To Delete", Position: "X", PositionOrder: 1}
	db.Create(&member)

	err := u.Delete("delete-me")

	assert.NoError(t, err)

	// Verify deletion in DB
	var count int64
	db.Model(&entity.OrganizationMember{}).Where("id = ?", "delete-me").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestOrganizationMemberUsecase_Delete_NotFound(t *testing.T) {
	ClearOrganizationMembers()
	u := newOrganizationMemberUsecase()

	err := u.Delete("non-existent-id")

	assert.Error(t, err) // Expect gorm.ErrRecordNotFound or similar
}
