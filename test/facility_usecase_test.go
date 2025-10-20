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

func newFacilityUsecase() usecase.FacilityUsecase {
	return usecase.NewFacilityUsecase(db, logrus.New(), validator.New())
}

func ClearFacilities() {
	if err := db.Where("id IS NOT NULL").Delete(&entity.Facility{}).Error; err != nil {
		log.Fatalf("Failed clear facilities: %+v", err)
	}
}

func TestFacilityUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
	ClearFacilities()

	// Seed: active and inactive, with different order_index
	g1 := entity.Facility{ID: "g1", Name: "B", ImageURL: "https://img1", IsActive: true, OrderIndex: 2}
	g2 := entity.Facility{ID: "g2", Name: "A", ImageURL: "https://img2", IsActive: false, OrderIndex: 1}
	g3 := entity.Facility{ID: "g3", Name: "C", ImageURL: "https://img3", IsActive: true, OrderIndex: 1}
	assert.NoError(t, db.Create(&g1).Error)
	assert.NoError(t, db.Create(&g2).Error)
	assert.NoError(t, db.Create(&g3).Error)

	u := newFacilityUsecase()
	list, err := u.GetPublic()
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// order_index ascending: g3 (1) then g1 (2)
	assert.Equal(t, "g3", list[0].ID)
	assert.Equal(t, "g1", list[1].ID)
}

func TestFacilityUsecase_Create_Update_Delete_Basic(t *testing.T) {
	ClearFacilities()
	u := newFacilityUsecase()

	// Create
	req := model.FacilityRequest{Name: "Name", ImageURL: "https://img", IsActive: true, OrderIndex: 3}
	created, err := u.Create(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)

	// Update
	updReq := model.FacilityRequest{Name: "New Name", ImageURL: "https://img2", IsActive: false, OrderIndex: 5}
	updated, err := u.Update(created.ID, updReq)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, false, updated.IsActive)
	assert.Equal(t, 5, updated.OrderIndex)

	// Delete
	err = u.Delete(created.ID)
	assert.NoError(t, err)

	var count int64
	assert.NoError(t, db.Model(&entity.Facility{}).Where("id = ?", created.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}
