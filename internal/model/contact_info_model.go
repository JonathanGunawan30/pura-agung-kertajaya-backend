package model

import "time"

// ContactInfoRequest defines the payload for updating contact info
type ContactInfoRequest struct {
	Address       string `json:"address" validate:"required,max=1000"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Email         string `json:"email" validate:"omitempty,email,max=100"`
	VisitingHours string `json:"visiting_hours" validate:"omitempty,max=100"`
	MapEmbedURL   string `json:"map_embed_url" validate:"omitempty,url"`
}

// ContactInfoResponse defines the API response for contact info
type ContactInfoResponse struct {
	ID            string    `json:"id"`
	Address       string    `json:"address"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	VisitingHours string    `json:"visiting_hours"`
	MapEmbedURL   string    `json:"map_embed_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
