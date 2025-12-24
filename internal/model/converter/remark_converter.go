package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToRemarkResponse(r *entity.Remark) model.RemarkResponse {
	return model.RemarkResponse{
		ID:         r.ID,
		EntityType: r.EntityType,
		Name:       r.Name,
		Position:   r.Position,
		ImageURL:   r.ImageURL,
		Content:    r.Content,
		IsActive:   r.IsActive,
		OrderIndex: r.OrderIndex,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func ToRemarkResponses(remarks []entity.Remark) []model.RemarkResponse {
	var responses []model.RemarkResponse
	for _, remark := range remarks {
		responses = append(responses, ToRemarkResponse(&remark))
	}
	return responses
}
