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

func newHeroSlideUsecase() usecase.HeroSlideUsecase {
    return usecase.NewHeroSlideUsecase(db, logrus.New(), validator.New())
}

func ClearHeroSlides() {
    err := db.Where("id IS NOT NULL").Delete(&entity.HeroSlide{}).Error
    if err != nil {
        log.Fatalf("Failed clear hero slides: %+v", err)
    }
}

func TestHeroSlideUsecase_Create_Success(t *testing.T) {
    ClearHeroSlides()
    u := newHeroSlideUsecase()

    req := model.HeroSlideRequest{
        ImageURL:   "https://example.com/image.jpg",
        OrderIndex: 1,
        IsActive:   true,
    }

    res, err := u.Create(req)
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, req.ImageURL, res.ImageURL)
    assert.Equal(t, req.OrderIndex, res.OrderIndex)
}

func TestHeroSlideUsecase_Create_ValidationError(t *testing.T) {
    ClearHeroSlides()
    u := newHeroSlideUsecase()

    req := model.HeroSlideRequest{} // missing ImageURL

    res, err := u.Create(req)
    assert.Error(t, err)
    assert.Nil(t, res)
}

func TestHeroSlideUsecase_GetAll_OrderedByIndex(t *testing.T) {
    ClearHeroSlides()

    // seed records with different order_index
    s1 := entity.HeroSlide{ID: "id-2", ImageUrl: "https://example.com/2.jpg", OrderIndex: 2, IsActive: true}
    s2 := entity.HeroSlide{ID: "id-1", ImageUrl: "https://example.com/1.jpg", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&s1).Error)
    assert.NoError(t, db.Create(&s2).Error)

    u := newHeroSlideUsecase()
    list, err := u.GetAll()
    assert.NoError(t, err)
    assert.Len(t, list, 2)
    assert.Equal(t, "https://example.com/1.jpg", list[0].ImageURL) // order_index 1 comes first
    assert.Equal(t, "https://example.com/2.jpg", list[1].ImageURL)
}

func TestHeroSlideUsecase_GetByID_NotFound(t *testing.T) {
    ClearHeroSlides()
    u := newHeroSlideUsecase()

    res, err := u.GetByID("not-exist")
    assert.Error(t, err)
    assert.Nil(t, res)
}

func TestHeroSlideUsecase_Update_Success(t *testing.T) {
    ClearHeroSlides()

    seed := entity.HeroSlide{ID: "slide-1", ImageUrl: "https://old.jpg", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&seed).Error)

    u := newHeroSlideUsecase()

    req := model.HeroSlideRequest{
        ImageURL:   "https://new.jpg",
        OrderIndex: 5,
        IsActive:   false,
    }

    res, err := u.Update(seed.ID, req)
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, req.ImageURL, res.ImageURL)
    assert.Equal(t, req.OrderIndex, res.OrderIndex)
    assert.Equal(t, req.IsActive, res.IsActive)
}

func TestHeroSlideUsecase_Delete_Success(t *testing.T) {
    ClearHeroSlides()
    seed := entity.HeroSlide{ID: "to-delete", ImageUrl: "https://img.jpg", OrderIndex: 1, IsActive: true}
    assert.NoError(t, db.Create(&seed).Error)

    u := newHeroSlideUsecase()
    err := u.Delete(seed.ID)
    assert.NoError(t, err)

    var count int64
    assert.NoError(t, db.Model(&entity.HeroSlide{}).Where("id = ?", seed.ID).Count(&count).Error)
    assert.Equal(t, int64(0), count)
}
