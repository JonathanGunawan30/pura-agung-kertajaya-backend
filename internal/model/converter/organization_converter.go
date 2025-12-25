package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToOrganizationResponse(g *entity.OrganizationMember) model.OrganizationResponse {
	return model.OrganizationResponse{
		ID:            g.ID,
		EntityType:    g.EntityType,
		Name:          g.Name,
		Position:      g.Position,
		PositionOrder: g.PositionOrder,
		OrderIndex:    g.OrderIndex,
		IsActive:      g.IsActive,
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}
}

func ToOrganizationResponses(organizations []entity.OrganizationMember) []model.OrganizationResponse {
	var responses []model.OrganizationResponse
	for _, organization := range organizations {
		responses = append(responses, ToOrganizationResponse(&organization))
	}

	return responses
}
