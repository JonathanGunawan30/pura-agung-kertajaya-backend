package model

import "time"

// CreateContactInfoRequest defines the payload for creating contact info
type CreateContactInfoRequest struct {
	EntityType    string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Address       string `json:"address" validate:"required,max=1000"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Email         string `json:"email" validate:"omitempty,email,max=100"`
	VisitingHours string `json:"visiting_hours" validate:"omitempty,max=100"`
	MapEmbedURL   string `json:"map_embed_url" validate:"omitempty,url"`
}

// UpdateContactInfoRequest defines the payload for updating contact info
type UpdateContactInfoRequest struct {
	Address       string `json:"address" validate:"required,max=1000"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Email         string `json:"email" validate:"omitempty,email,max=100"`
	VisitingHours string `json:"visiting_hours" validate:"omitempty,max=100"`
	MapEmbedURL   string `json:"map_embed_url" validate:"omitempty,url"`
}

// ContactInfoResponse defines the API response for contact info
type ContactInfoResponse struct {
	ID            string    `json:"id"`
	EntityType    string    `json:"entity_type"`
	Address       string    `json:"address"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	VisitingHours string    `json:"visiting_hours"`
	MapEmbedURL   string    `json:"map_embed_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
