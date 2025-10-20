package test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

func newSiteIdentityUsecase() usecase.SiteIdentityUsecase {
	return usecase.NewSiteIdentityUsecase(db, logrus.New(), validator.New())
}

func ClearSiteIdentity() {
	if err := db.Where("id IS NOT NULL").Delete(&entity.SiteIdentity{}).Error; err != nil {
		log.Fatalf("Failed clear site identity: %+v", err)
	}
}

func TestSiteIdentityUsecase_Create_Success(t *testing.T) {
	ClearSiteIdentity()
	u := newSiteIdentityUsecase()

	req := model.SiteIdentityRequest{SiteName: "Pura", LogoURL: "https://logo", Tagline: "Tag"}
	res, err := u.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Pura", res.SiteName)
}

func TestSiteIdentityUsecase_Create_ValidationError(t *testing.T) {
	ClearSiteIdentity()
	u := newSiteIdentityUsecase()

	req := model.SiteIdentityRequest{} // missing site_name
	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSiteIdentityUsecase_GetAll(t *testing.T) {
	ClearSiteIdentity()

	// seed two
	s1 := entity.SiteIdentity{ID: "s1", SiteName: "A"}
	s2 := entity.SiteIdentity{ID: "s2", SiteName: "B"}
	assert.NoError(t, db.Create(&s1).Error)
	assert.NoError(t, db.Create(&s2).Error)

	u := newSiteIdentityUsecase()
	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestSiteIdentityUsecase_GetByID_NotFound(t *testing.T) {
	ClearSiteIdentity()
	u := newSiteIdentityUsecase()
	res, err := u.GetByID("missing")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSiteIdentityUsecase_Update_Success(t *testing.T) {
	ClearSiteIdentity()
	seed := entity.SiteIdentity{ID: "sid-1", SiteName: "Old", Tagline: "Old"}
	assert.NoError(t, db.Create(&seed).Error)

	u := newSiteIdentityUsecase()
	req := model.SiteIdentityRequest{SiteName: "New", Tagline: "New", PrimaryButtonText: "Go"}
	res, err := u.Update(seed.ID, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New", res.SiteName)
	assert.Equal(t, "Go", res.PrimaryButtonText)
}

func TestSiteIdentityUsecase_Delete_Success(t *testing.T) {
	ClearSiteIdentity()
	seed := entity.SiteIdentity{ID: "del-1", SiteName: "Name"}
	assert.NoError(t, db.Create(&seed).Error)

	u := newSiteIdentityUsecase()
	err := u.Delete(seed.ID)
	assert.NoError(t, err)

	var count int64
	assert.NoError(t, db.Model(&entity.SiteIdentity{}).Where("id = ?", seed.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestSiteIdentityUsecase_GetPublic_ReturnsLatest(t *testing.T) {
	ClearSiteIdentity()

	idOld := entity.SiteIdentity{ID: "old", SiteName: "Old"}
	assert.NoError(t, db.Create(&idOld).Error)

	time.Sleep(time.Millisecond * 10)

	idNew := entity.SiteIdentity{ID: "new", SiteName: "New"}
	assert.NoError(t, db.Create(&idNew).Error)

	u := newSiteIdentityUsecase()
	result, err := u.GetPublic()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New", result.SiteName)
}
