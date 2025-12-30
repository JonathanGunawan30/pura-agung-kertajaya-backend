package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToCategoryResponse(c *entity.Category) model.CategoryResponse {
	return model.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func ToCategoryResponses(categories []entity.Category) []model.CategoryResponse {
	var responses []model.CategoryResponse
	for _, category := range categories {
		responses = append(responses, ToCategoryResponse(&category))
	}
	return responses
}
