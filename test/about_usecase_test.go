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

func newAboutUsecase() usecase.AboutUsecase {
    return usecase.NewAboutUsecase(db, logrus.New(), validator.New())
}

func clearAboutTables() {
    // delete children first then parents to be safe
    if err := db.Where("id IS NOT NULL").Delete(&entity.AboutValue{}).Error; err != nil {
        log.Fatalf("Failed clear about_values: %+v", err)
    }
    if err := db.Where("id IS NOT NULL").Delete(&entity.AboutSection{}).Error; err != nil {
        log.Fatalf("Failed clear about_section: %+v", err)
    }
}

func TestAboutUsecase_Create_WithValues(t *testing.T) {
    clearAboutTables()
    u := newAboutUsecase()

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

    res, err := u.Create(req)
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 2, len(res.Values))
    // Ensure ordered by order_index asc: Mission first
    assert.Equal(t, "Mission", res.Values[0].Title)
}

func TestAboutUsecase_GetPublic_FilterActive(t *testing.T) {
    clearAboutTables()

    // Seed two sections: active and inactive
    a1 := entity.AboutSection{ID: "a1", Title: "Active", Description: "d", IsActive: true}
    a2 := entity.AboutSection{ID: "a2", Title: "Inactive", Description: "d", IsActive: false}
    assert.NoError(t, db.Create(&a1).Error)
    assert.NoError(t, db.Create(&a2).Error)
    // Seed some values for a1 with different order indexes
    v1 := entity.AboutValue{ID: "v1", AboutID: "a1", Title: "B", Value: "b", OrderIndex: 2}
    v2 := entity.AboutValue{ID: "v2", AboutID: "a1", Title: "A", Value: "a", OrderIndex: 1}
    assert.NoError(t, db.Create(&v1).Error)
    assert.NoError(t, db.Create(&v2).Error)

    u := newAboutUsecase()
    list, err := u.GetPublic()
    assert.NoError(t, err)
    assert.Len(t, list, 1)
    assert.Equal(t, "Active", list[0].Title)
    // Values ordered by order_index
    assert.Equal(t, "A", list[0].Values[0].Title)
}

func TestAboutUsecase_GetByID_NotFound(t *testing.T) {
    clearAboutTables()
    u := newAboutUsecase()
    res, err := u.GetByID("missing")
    assert.Error(t, err)
    assert.Nil(t, res)
}

func TestAboutUsecase_Update_ReplacesValues(t *testing.T) {
    clearAboutTables()
    // Seed section with one value
    a := entity.AboutSection{ID: "ab-1", Title: "Old", Description: "d", IsActive: true}
    assert.NoError(t, db.Create(&a).Error)
    ov := entity.AboutValue{ID: "ov1", AboutID: a.ID, Title: "OldV", Value: "x", OrderIndex: 1}
    assert.NoError(t, db.Create(&ov).Error)

    u := newAboutUsecase()
    req := model.AboutSectionRequest{
        Title:       "New",
        Description: "nd",
        ImageURL:    "https://new",
        IsActive:    false,
        Values: []model.AboutValueRequest{
            {Title: "New1", Value: "n1", OrderIndex: 2},
            {Title: "New0", Value: "n0", OrderIndex: 1},
        },
    }
    res, err := u.Update(a.ID, req)
    assert.NoError(t, err)
    assert.Equal(t, "New", res.Title)
    assert.Equal(t, false, res.IsActive)
    assert.Equal(t, 2, len(res.Values))
    assert.Equal(t, "New0", res.Values[0].Title)

    // Ensure only new values remain
    var count int64
    assert.NoError(t, db.Model(&entity.AboutValue{}).Where("about_id = ? AND title = ?", a.ID, "OldV").Count(&count).Error)
    assert.Equal(t, int64(0), count)
}

func TestAboutUsecase_Delete_Success(t *testing.T) {
    clearAboutTables()
    a := entity.AboutSection{ID: "to-del", Title: "T", Description: "d", IsActive: true}
    assert.NoError(t, db.Create(&a).Error)
    v := entity.AboutValue{ID: "vv", AboutID: a.ID, Title: "V", Value: "v", OrderIndex: 1}
    assert.NoError(t, db.Create(&v).Error)

    u := newAboutUsecase()
    err := u.Delete(a.ID)
    assert.NoError(t, err)

    var c1, c2 int64
    assert.NoError(t, db.Model(&entity.AboutSection{}).Where("id = ?", a.ID).Count(&c1).Error)
    assert.NoError(t, db.Model(&entity.AboutValue{}).Where("about_id = ?", a.ID).Count(&c2).Error)
    assert.Equal(t, int64(0), c1)
    assert.Equal(t, int64(0), c2) // cascade ensured
}
