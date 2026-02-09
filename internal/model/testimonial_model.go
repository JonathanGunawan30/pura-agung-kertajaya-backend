package model

import "time"

type TestimonialRequest struct {
	Name       string `json:"name" validate:"required"`
	AvatarURL  string `json:"avatar_url"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	Comment    string `json:"comment" validate:"required"`
	IsActive   bool   `json:"is_active"`
	OrderIndex int    `json:"order_index"`
}

type TestimonialResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	IsActive   bool      `json:"is_active"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
