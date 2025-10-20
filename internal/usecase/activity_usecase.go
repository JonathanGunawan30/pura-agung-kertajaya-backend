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

type ActivityUsecase interface {
    GetAll() ([]model.ActivityResponse, error)
    GetPublic() ([]model.ActivityResponse, error)
    GetByID(id string) (*model.ActivityResponse, error)
    Create(req model.ActivityRequest) (*model.ActivityResponse, error)
    Update(id string, req model.ActivityRequest) (*model.ActivityResponse, error)
    Delete(id string) error
}

type activityUsecase struct {
    db       *gorm.DB
    repo     *repository.Repository[entity.Activity]
    log      *logrus.Logger
    validate *validator.Validate
}

func NewActivityUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) ActivityUsecase {
    return &activityUsecase{
        db:       db,
        repo:     &repository.Repository[entity.Activity]{DB: db},
        log:      log,
        validate: validate,
    }
}

func (u *activityUsecase) GetAll() ([]model.ActivityResponse, error) {
    var items []entity.Activity
    if err := u.db.Order("order_index ASC").Find(&items).Error; err != nil {
        return nil, err
    }
    resp := make([]model.ActivityResponse, 0, len(items))
    for _, a := range items {
        resp = append(resp, converter.ToActivityResponse(a))
    }
    return resp, nil
}

func (u *activityUsecase) GetPublic() ([]model.ActivityResponse, error) {
    var items []entity.Activity
    if err := u.db.Where("is_active = ?", true).Order("order_index ASC").Find(&items).Error; err != nil {
        return nil, err
    }
    resp := make([]model.ActivityResponse, 0, len(items))
    for _, a := range items {
        resp = append(resp, converter.ToActivityResponse(a))
    }
    return resp, nil
}

func (u *activityUsecase) GetByID(id string) (*model.ActivityResponse, error) {
    var a entity.Activity
    if err := u.repo.FindById(u.db, &a, id); err != nil {
        return nil, err
    }
    r := converter.ToActivityResponse(a)
    return &r, nil
}

func (u *activityUsecase) Create(req model.ActivityRequest) (*model.ActivityResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }
    a := entity.Activity{
        ID:          uuid.New().String(),
        Title:       req.Title,
        Description: req.Description,
        TimeInfo:    req.TimeInfo,
        Location:    req.Location,
        OrderIndex:  req.OrderIndex,
        IsActive:    req.IsActive,
    }
    if err := u.repo.Create(u.db, &a); err != nil {
        return nil, err
    }
    r := converter.ToActivityResponse(a)
    return &r, nil
}

func (u *activityUsecase) Update(id string, req model.ActivityRequest) (*model.ActivityResponse, error) {
    if err := u.validate.Struct(req); err != nil {
        return nil, err
    }
    var a entity.Activity
    if err := u.repo.FindById(u.db, &a, id); err != nil {
        return nil, err
    }
    a.Title = req.Title
    a.Description = req.Description
    a.TimeInfo = req.TimeInfo
    a.Location = req.Location
    a.OrderIndex = req.OrderIndex
    a.IsActive = req.IsActive

    if err := u.repo.Update(u.db, &a); err != nil {
        return nil, err
    }
    r := converter.ToActivityResponse(a)
    return &r, nil
}

func (u *activityUsecase) Delete(id string) error {
    var a entity.Activity
    if err := u.repo.FindById(u.db, &a, id); err != nil {
        return err
    }
    return u.repo.Delete(u.db, &a)
}
