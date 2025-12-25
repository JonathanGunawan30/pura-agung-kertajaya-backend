package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToActivityResponse(a *entity.Activity) model.ActivityResponse {
	return model.ActivityResponse{
		ID:          a.ID,
		EntityType:  a.EntityType,
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

func ToActivityResponses(activities []entity.Activity) []model.ActivityResponse {
	var responses []model.ActivityResponse
	for _, activity := range activities {
		responses = append(responses, ToActivityResponse(&activity))
	}

	return responses
}
