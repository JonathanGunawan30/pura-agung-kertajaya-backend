package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

// ToContactInfoResponse converts entity.ContactInfo to model.ContactInfoResponse
func ToContactInfoResponse(e entity.ContactInfo) model.ContactInfoResponse {
	return model.ContactInfoResponse{
		ID:            e.ID,
		EntityType:    e.EntityType,
		Address:       e.Address,
		Phone:         e.Phone,
		Email:         e.Email,
		VisitingHours: e.VisitingHours,
		MapEmbedURL:   e.MapEmbedURL,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// ToContactInfoResponses converts a slice of entity.ContactInfo to a slice of model.ContactInfoResponse
func ToContactInfoResponses(items []entity.ContactInfo) []model.ContactInfoResponse {
	var responses []model.ContactInfoResponse
	for _, item := range items {
		responses = append(responses, ToContactInfoResponse(item))
	}
	return responses
}
