package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToOrganizationResponse(g entity.OrganizationMember) model.OrganizationResponse {
	return model.OrganizationResponse{
		ID:            g.ID,
		Name:          g.Name,
		Position:      g.Position,
		PositionOrder: g.PositionOrder,
		OrderIndex:    g.OrderIndex,
		IsActive:      g.IsActive,
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}
}
