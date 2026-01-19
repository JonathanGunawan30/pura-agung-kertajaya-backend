package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToFacilityResponse(g entity.Facility) model.FacilityResponse {
	return model.FacilityResponse{
		ID:          g.ID,
		EntityType:  g.EntityType,
		Name:        g.Name,
		Description: g.Description,
		Images:      ToImageVariants(g.Images),
		OrderIndex:  g.OrderIndex,
		IsActive:    g.IsActive,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

func ToFacilityResponses(facilities []entity.Facility) []model.FacilityResponse {
	var responses []model.FacilityResponse
	for _, facility := range facilities {
		responses = append(responses, ToFacilityResponse(facility))
	}

	return responses
}
