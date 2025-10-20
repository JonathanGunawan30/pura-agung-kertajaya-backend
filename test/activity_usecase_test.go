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

func newActivityUsecase() usecase.ActivityUsecase {
    return usecase.NewActivityUsecase(db, logrus.New(), validator.New())
}

func ClearActivities() {
    if err := db.Where("id IS NOT NULL").Delete(&entity.Activity{}).Error; err != nil {
        log.Fatalf("Failed clear activities: %+v", err)
    }
}

func TestActivityUsecase_Create_Success(t *testing.T) {
    ClearActivities()
    u := newActivityUsecase()

    req := model.ActivityRequest{
        Title:       "Upacara",
        Description: "Deskripsi",
        TimeInfo:    "08:00",
        Location:    "Pura",
        OrderIndex:  1,
        IsActive:    true,
    }

    res, err := u.Create(req)
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, req.Title, res.Title)
    assert.Equal(t, req.Description, res.Description)
}

func TestActivityUsecase_Create_ValidationError(t *testing.T) {
    ClearActivities()
    u := newActivityUsecase()

    req := model.ActivityRequest{} // missing required fields
    res, err := u.Create(req)
    assert.Error(t, err)
    assert.Nil(t, res)
}

func TestActivityUsecase_GetAll_OrderedByIndex(t *testing.T) {
    ClearActivities()

    a1 := entity.Activity{ID: "a2", Title: "B", Description: "d", OrderIndex: 2, IsActive: true}
    a2 := entity.Activity{ID: "a1", Title: "A", Description: "d", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&a1).Error)
    assert.NoError(t, db.Create(&a2).Error)

    u := newActivityUsecase()
    list, err := u.GetAll()
    assert.NoError(t, err)
    assert.Len(t, list, 2)
    assert.Equal(t, "A", list[0].Title)
    assert.Equal(t, "B", list[1].Title)
}

func TestActivityUsecase_GetPublic_FilterActiveAndOrder(t *testing.T) {
    ClearActivities()

    a1 := entity.Activity{ID: "a1", Title: "B", Description: "d", OrderIndex: 2, IsActive: true}
    a2 := entity.Activity{ID: "a2", Title: "X", Description: "d", OrderIndex: 1, IsActive: false}
    a3 := entity.Activity{ID: "a3", Title: "A", Description: "d", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&a1).Error)
    assert.NoError(t, db.Create(&a2).Error)
    assert.NoError(t, db.Create(&a3).Error)

    u := newActivityUsecase()
    list, err := u.GetPublic()
    assert.NoError(t, err)
    assert.Len(t, list, 2)
    assert.Equal(t, "A", list[0].Title)
    assert.Equal(t, "B", list[1].Title)
}

func TestActivityUsecase_GetByID_NotFound(t *testing.T) {
    ClearActivities()
    u := newActivityUsecase()

    res, err := u.GetByID("missing")
    assert.Error(t, err)
    assert.Nil(t, res)
}

func TestActivityUsecase_Update_Success(t *testing.T) {
    ClearActivities()

    seed := entity.Activity{ID: "act-1", Title: "Old", Description: "d", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&seed).Error)

    u := newActivityUsecase()

    req := model.ActivityRequest{Title: "New", Description: "new d", TimeInfo: "09:00", Location: "Pura", OrderIndex: 5, IsActive: false}
    res, err := u.Update(seed.ID, req)
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, "New", res.Title)
    assert.Equal(t, false, res.IsActive)
    assert.Equal(t, 5, res.OrderIndex)
}

func TestActivityUsecase_Delete_Success(t *testing.T) {
    ClearActivities()
    seed := entity.Activity{ID: "to-del", Title: "Del", Description: "d", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&seed).Error)

    u := newActivityUsecase()
    err := u.Delete(seed.ID)
    assert.NoError(t, err)

    var count int64
    assert.NoError(t, db.Model(&entity.Activity{}).Where("id = ?", seed.ID).Count(&count).Error)
    assert.Equal(t, int64(0), count)
}
