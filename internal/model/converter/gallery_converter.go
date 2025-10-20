package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToGalleryResponse(g entity.Gallery) model.GalleryResponse {
	return model.GalleryResponse{
		ID:          g.ID,
		Title:       g.Title,
		Description: g.Description,
		ImageURL:    g.ImageURL,
		OrderIndex:  g.OrderIndex,
		IsActive:    g.IsActive,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}
