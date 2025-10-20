package converter

import (
    "pura-agung-kertajaya-backend/internal/entity"
    "pura-agung-kertajaya-backend/internal/model"
)

// ToContactInfoResponse converts entity.ContactInfo to model.ContactInfoResponse
func ToContactInfoResponse(e entity.ContactInfo) model.ContactInfoResponse {
    return model.ContactInfoResponse{
        ID:            e.ID,
        Address:       e.Address,
        Phone:         e.Phone,
        Email:         e.Email,
        VisitingHours: e.VisitingHours,
        MapEmbedURL:   e.MapEmbedURL,
        CreatedAt:     e.CreatedAt,
        UpdatedAt:     e.UpdatedAt,
    }
}
