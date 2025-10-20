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

func newContactInfoUsecase() usecase.ContactInfoUsecase {
	return usecase.NewContactInfoUsecase(db, logrus.New(), validator.New())
}

func ClearContactInfo() {
	err := db.Where("id IS NOT NULL").Delete(&entity.ContactInfo{}).Error
	if err != nil {
		log.Fatalf("Failed clear contact info: %+v", err)
	}
}

func TestContactInfoUsecase_Create_Success(t *testing.T) {
	ClearContactInfo()
	u := newContactInfoUsecase()

	req := model.ContactInfoRequest{
		Address:       "Jl. Contoh No.1",
		Phone:         "+62 8123456789",
		Email:         "info@example.com",
		VisitingHours: "08:00 - 17:00",
		MapEmbedURL:   "https://maps.google.com/?q=x",
	}

	res, err := u.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Address, res.Address)
	assert.Equal(t, req.Email, res.Email)
}

func TestContactInfoUsecase_Create_ValidationError(t *testing.T) {
	ClearContactInfo()
	u := newContactInfoUsecase()

	req := model.ContactInfoRequest{} // address required

	res, err := u.Create(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContactInfoUsecase_GetAll(t *testing.T) {
	ClearContactInfo()

	c1 := entity.ContactInfo{ID: "1", Address: "A"}
	c2 := entity.ContactInfo{ID: "2", Address: "B"}
	assert.NoError(t, db.Create(&c1).Error)
	assert.NoError(t, db.Create(&c2).Error)

	u := newContactInfoUsecase()
	list, err := u.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestContactInfoUsecase_GetByID_NotFound(t *testing.T) {
	ClearContactInfo()
	u := newContactInfoUsecase()

	res, err := u.GetByID("not-exists")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContactInfoUsecase_Update_Success(t *testing.T) {
	ClearContactInfo()
	seed := entity.ContactInfo{ID: "ci-1", Address: "Old Addr", Email: "old@example.com"}
	assert.NoError(t, db.Create(&seed).Error)

	u := newContactInfoUsecase()
	req := model.ContactInfoRequest{Address: "New Addr", Email: "new@example.com", Phone: "000"}
	res, err := u.Update(seed.ID, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New Addr", res.Address)
	assert.Equal(t, "new@example.com", res.Email)
}

func TestContactInfoUsecase_Delete_Success(t *testing.T) {
	ClearContactInfo()
	seed := entity.ContactInfo{ID: "to-del", Address: "Addr"}
	assert.NoError(t, db.Create(&seed).Error)

	u := newContactInfoUsecase()
	err := u.Delete(seed.ID)
	assert.NoError(t, err)

	var count int64
	assert.NoError(t, db.Model(&entity.ContactInfo{}).Where("id = ?", seed.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}
