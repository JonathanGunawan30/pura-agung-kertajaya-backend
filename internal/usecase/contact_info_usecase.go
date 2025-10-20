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

type ContactInfoUsecase interface {
    GetAll() ([]model.ContactInfoResponse, error)
    GetByID(id string) (*model.ContactInfoResponse, error)
    Create(req model.ContactInfoRequest) (*model.ContactInfoResponse, error)
    Update(id string, req model.ContactInfoRequest) (*model.ContactInfoResponse, error)
    Delete(id string) error
}

type contactInfoUsecase struct {
    db       *gorm.DB
    repo     *repository.Repository[entity.ContactInfo]
    log      *logrus.Logger
    validate *validator.Validate
}

func NewContactInfoUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) ContactInfoUsecase {
    return &contactInfoUsecase{
        db:       db,
        repo:     &repository.Repository[entity.ContactInfo]{DB: db},
        log:      log,
        validate: validate,
    }
}

func (u *contactInfoUsecase) GetAll() ([]model.ContactInfoResponse, error) {
    var items []entity.ContactInfo
    if err := u.db.Order("created_at ASC").Find(&items).Error; err != nil {
        return nil, err
    }
    resp := make([]model.ContactInfoResponse, 0, len(items))
    for _, e := range items {
        resp = append(resp, converter.ToContactInfoResponse(e))
    }
    return resp, nil
}

func (u *contactInfoUsecase) GetByID(id string) (*model.ContactInfoResponse, error) {
    var e entity.ContactInfo
    if err := u.repo.FindById(u.db, &e, id); err != nil {
        return nil, err
    }
    r := converter.ToContactInfoResponse(e)
    return &r, nil
}

func (u *contactInfoUsecase) Create(req model.ContactInfoRequest) (*model.ContactInfoResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }

    e := entity.ContactInfo{
        ID:            uuid.New().String(),
        Address:       req.Address,
        Phone:         req.Phone,
        Email:         req.Email,
        VisitingHours: req.VisitingHours,
        MapEmbedURL:   req.MapEmbedURL,
    }

    if err := u.repo.Create(u.db, &e); err != nil {
        return nil, err
    }

    r := converter.ToContactInfoResponse(e)
    return &r, nil
}

func (u *contactInfoUsecase) Update(id string, req model.ContactInfoRequest) (*model.ContactInfoResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }

    var e entity.ContactInfo
    if err := u.repo.FindById(u.db, &e, id); err != nil {
        return nil, err
    }

    e.Address = req.Address
    e.Phone = req.Phone
    e.Email = req.Email
    e.VisitingHours = req.VisitingHours
    e.MapEmbedURL = req.MapEmbedURL

    if err := u.repo.Update(u.db, &e); err != nil {
        return nil, err
    }

    r := converter.ToContactInfoResponse(e)
    return &r, nil
}

func (u *contactInfoUsecase) Delete(id string) error {
    var e entity.ContactInfo
    if err := u.repo.FindById(u.db, &e, id); err != nil {
        return err
    }
    return u.repo.Delete(u.db, &e)
}
