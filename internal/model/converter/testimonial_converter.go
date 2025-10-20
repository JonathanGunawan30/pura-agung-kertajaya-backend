package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToTestimonialResponse(t entity.Testimonial) model.TestimonialResponse {
	return model.TestimonialResponse{
		ID:         t.ID,
		Name:       t.Name,
		AvatarURL:  t.AvatarURL,
		Rating:     t.Rating,
		Comment:    t.Comment,
		IsActive:   t.IsActive,
		OrderIndex: t.OrderIndex,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
