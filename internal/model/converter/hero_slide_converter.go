package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToHeroSlideResponse(h entity.HeroSlide) model.HeroSlideResponse {
	return model.HeroSlideResponse{
		ID:         h.ID,
		EntityType: h.EntityType,
		ImageURL:   h.ImageUrl,
		OrderIndex: h.OrderIndex,
		IsActive:   h.IsActive,
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}
}
