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

type SiteIdentityUsecase interface {
    GetAll() ([]model.SiteIdentityResponse, error)
    // GetPublic returns the latest SiteIdentity (by created_at DESC) for public consumption
    GetPublic() (*model.SiteIdentityResponse, error)
    GetByID(id string) (*model.SiteIdentityResponse, error)
    Create(req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error)
    Update(id string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error)
    Delete(id string) error
}

type siteIdentityUsecase struct {
    db       *gorm.DB
    repo     *repository.Repository[entity.SiteIdentity]
    log      *logrus.Logger
    validate *validator.Validate
}

func NewSiteIdentityUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) SiteIdentityUsecase {
    return &siteIdentityUsecase{
        db:       db,
        repo:     &repository.Repository[entity.SiteIdentity]{DB: db},
        log:      log,
        validate: validate,
    }
}

func (u *siteIdentityUsecase) GetAll() ([]model.SiteIdentityResponse, error) {
    var items []entity.SiteIdentity
    if err := u.db.Order("created_at ASC").Find(&items).Error; err != nil {
        return nil, err
    }
    resp := make([]model.SiteIdentityResponse, 0, len(items))
    for _, e := range items {
        resp = append(resp, converter.ToSiteIdentityResponse(e))
    }
    return resp, nil
}

func (u *siteIdentityUsecase) GetPublic() (*model.SiteIdentityResponse, error) {
    var e entity.SiteIdentity
    if err := u.db.Order("created_at DESC").Limit(1).Take(&e).Error; err != nil {
        return nil, err
    }
    r := converter.ToSiteIdentityResponse(e)
    return &r, nil
}

func (u *siteIdentityUsecase) GetByID(id string) (*model.SiteIdentityResponse, error) {
    var e entity.SiteIdentity
    if err := u.repo.FindById(u.db, &e, id); err != nil {
        return nil, err
    }
    r := converter.ToSiteIdentityResponse(e)
    return &r, nil
}

func (u *siteIdentityUsecase) Create(req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }
    e := entity.SiteIdentity{
        ID:                  uuid.New().String(),
        SiteName:            req.SiteName,
        LogoURL:             req.LogoURL,
        Tagline:             req.Tagline,
        PrimaryButtonText:   req.PrimaryButtonText,
        PrimaryButtonLink:   req.PrimaryButtonLink,
        SecondaryButtonText: req.SecondaryButtonText,
        SecondaryButtonLink: req.SecondaryButtonLink,
    }
    if err := u.repo.Create(u.db, &e); err != nil {
        return nil, err
    }
    r := converter.ToSiteIdentityResponse(e)
    return &r, nil
}

func (u *siteIdentityUsecase) Update(id string, req model.SiteIdentityRequest) (*model.SiteIdentityResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }
    var e entity.SiteIdentity
    if err := u.repo.FindById(u.db, &e, id); err != nil {
        return nil, err
    }
    e.SiteName = req.SiteName
    e.LogoURL = req.LogoURL
    e.Tagline = req.Tagline
    e.PrimaryButtonText = req.PrimaryButtonText
    e.PrimaryButtonLink = req.PrimaryButtonLink
    e.SecondaryButtonText = req.SecondaryButtonText
    e.SecondaryButtonLink = req.SecondaryButtonLink

    if err := u.repo.Update(u.db, &e); err != nil {
        return nil, err
    }
    r := converter.ToSiteIdentityResponse(e)
    return &r, nil
}

func (u *siteIdentityUsecase) Delete(id string) error {
    var e entity.SiteIdentity
    if err := u.repo.FindById(u.db, &e, id); err != nil {
        return err
    }
    return u.repo.Delete(u.db, &e)
}
