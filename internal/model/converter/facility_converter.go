package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToFacilityResponse(g entity.Facility) model.FacilityResponse {
	return model.FacilityResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		ImageURL:    g.ImageURL,
		OrderIndex:  g.OrderIndex,
		IsActive:    g.IsActive,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}
