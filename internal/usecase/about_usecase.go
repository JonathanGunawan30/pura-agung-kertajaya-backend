package usecase

import (
    "pura-agung-kertajaya-backend/internal/entity"
    "pura-agung-kertajaya-backend/internal/model"
    "pura-agung-kertajaya-backend/internal/model/converter"
    "pura-agung-kertajaya-backend/internal/repository"

    "github.com/go-playground/validator/v10"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
)

type AboutUsecase interface {
    GetAll() ([]model.AboutSectionResponse, error)
    GetPublic() ([]model.AboutSectionResponse, error)
    GetByID(id string) (*model.AboutSectionResponse, error)
    Create(req model.AboutSectionRequest) (*model.AboutSectionResponse, error)
    Update(id string, req model.AboutSectionRequest) (*model.AboutSectionResponse, error)
    Delete(id string) error
}

type aboutUsecase struct {
    db       *gorm.DB
    repoAbout *repository.Repository[entity.AboutSection]
    repoValue *repository.Repository[entity.AboutValue]
    log      *logrus.Logger
    validate *validator.Validate
}

func NewAboutUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) AboutUsecase {
    return &aboutUsecase{
        db:        db,
        repoAbout: &repository.Repository[entity.AboutSection]{DB: db},
        repoValue: &repository.Repository[entity.AboutValue]{DB: db},
        log:       log,
        validate:  validate,
    }
}

func preloadValuesOrdered(tx *gorm.DB) *gorm.DB {
    return tx.Preload("Values", func(db *gorm.DB) *gorm.DB { return db.Order("order_index ASC") })
}

func (u *aboutUsecase) GetAll() ([]model.AboutSectionResponse, error) {
    var list []entity.AboutSection
    if err := preloadValuesOrdered(u.db).Order("created_at ASC").Find(&list).Error; err != nil {
        return nil, err
    }
    resp := make([]model.AboutSectionResponse, 0, len(list))
    for _, a := range list {
        resp = append(resp, converter.ToAboutSectionResponse(a))
    }
    return resp, nil
}

func (u *aboutUsecase) GetPublic() ([]model.AboutSectionResponse, error) {
    var list []entity.AboutSection
    if err := preloadValuesOrdered(u.db).Where("is_active = ?", true).Order("created_at ASC").Find(&list).Error; err != nil {
        return nil, err
    }
    resp := make([]model.AboutSectionResponse, 0, len(list))
    for _, a := range list {
        resp = append(resp, converter.ToAboutSectionResponse(a))
    }
    return resp, nil
}

func (u *aboutUsecase) GetByID(id string) (*model.AboutSectionResponse, error) {
    var a entity.AboutSection
    if err := preloadValuesOrdered(u.db).Where("id = ?", id).Take(&a).Error; err != nil {
        return nil, err
    }
    r := converter.ToAboutSectionResponse(a)
    return &r, nil
}

func (u *aboutUsecase) Create(req model.AboutSectionRequest) (*model.AboutSectionResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }

    id := uuid.New().String()
    a := entity.AboutSection{
        ID:          id,
        Title:       req.Title,
        Description: req.Description,
        ImageURL:    req.ImageURL,
        IsActive:    req.IsActive,
    }

    err := u.db.Transaction(func(tx *gorm.DB) error {
        if err := u.repoAbout.Create(tx, &a); err != nil {
            return err
        }
        if len(req.Values) > 0 {
            values := make([]entity.AboutValue, 0, len(req.Values))
            for _, v := range req.Values {
                values = append(values, entity.AboutValue{
                    ID:         uuid.New().String(),
                    AboutID:    a.ID,
                    Title:      v.Title,
                    Value:      v.Value,
                    OrderIndex: v.OrderIndex,
                })
            }
            if err := tx.Create(&values).Error; err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    // reload with values
    if err := preloadValuesOrdered(u.db).Where("id = ?", a.ID).Take(&a).Error; err != nil {
        return nil, err
    }
    r := converter.ToAboutSectionResponse(a)
    return &r, nil
}

func (u *aboutUsecase) Update(id string, req model.AboutSectionRequest) (*model.AboutSectionResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }

    var a entity.AboutSection
    if err := u.repoAbout.FindById(u.db, &a, id); err != nil {
        return nil, err
    }

    a.Title = req.Title
    a.Description = req.Description
    a.ImageURL = req.ImageURL
    a.IsActive = req.IsActive

    err := u.db.Transaction(func(tx *gorm.DB) error {
        if err := u.repoAbout.Update(tx, &a); err != nil {
            return err
        }
        // simple sync strategy: delete existing values then recreate in order
        if err := tx.Where("about_id = ?", a.ID).Delete(&entity.AboutValue{}).Error; err != nil {
            return err
        }
        if len(req.Values) > 0 {
            values := make([]entity.AboutValue, 0, len(req.Values))
            for _, v := range req.Values {
                values = append(values, entity.AboutValue{
                    ID:         uuid.New().String(),
                    AboutID:    a.ID,
                    Title:      v.Title,
                    Value:      v.Value,
                    OrderIndex: v.OrderIndex,
                })
            }
            if err := tx.Create(&values).Error; err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    if err := preloadValuesOrdered(u.db).Where("id = ?", a.ID).Take(&a).Error; err != nil {
        return nil, err
    }
    r := converter.ToAboutSectionResponse(a)
    return &r, nil
}

func (u *aboutUsecase) Delete(id string) error {
    var a entity.AboutSection
    if err := u.repoAbout.FindById(u.db, &a, id); err != nil {
        return err
    }
    return u.repoAbout.Delete(u.db, &a)
}
