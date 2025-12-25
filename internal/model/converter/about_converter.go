package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToAboutValueResponse(v entity.AboutValue) model.AboutValueResponse {
	return model.AboutValueResponse{
		ID:         v.ID,
		Title:      v.Title,
		Value:      v.Value,
		OrderIndex: v.OrderIndex,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
}

func ToAboutSectionResponse(a entity.AboutSection) model.AboutSectionResponse {
	values := make([]model.AboutValueResponse, 0, len(a.Values))
	for _, v := range a.Values {
		values = append(values, ToAboutValueResponse(v))
	}
	return model.AboutSectionResponse{
		ID:          a.ID,
		EntityType:  a.EntityType,
		Title:       a.Title,
		Description: a.Description,
		ImageURL:    a.ImageURL,
		IsActive:    a.IsActive,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Values:      values,
	}
}
