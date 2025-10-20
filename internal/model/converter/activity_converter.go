package converter

import (
    "pura-agung-kertajaya-backend/internal/entity"
    "pura-agung-kertajaya-backend/internal/model"
)

func ToActivityResponse(a entity.Activity) model.ActivityResponse {
    return model.ActivityResponse{
        ID:          a.ID,
        Title:       a.Title,
        Description: a.Description,
        TimeInfo:    a.TimeInfo,
        Location:    a.Location,
        OrderIndex:  a.OrderIndex,
        IsActive:    a.IsActive,
        CreatedAt:   a.CreatedAt,
        UpdatedAt:   a.UpdatedAt,
    }
}
