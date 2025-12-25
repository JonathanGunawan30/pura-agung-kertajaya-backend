package model

import "time"

type HeroSlideRequest struct {
	EntityType string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	ImageURL   string `json:"image_url" validate:"required,url|uri|startswith=http"`
	OrderIndex int    `json:"order_index"`
	IsActive   bool   `json:"is_active"`
}

type HeroSlideResponse struct {
	ID         string    `json:"id"`
	EntityType string    `json:"entity_type"`
	ImageURL   string    `json:"image_url"`
	OrderIndex int       `json:"order_index"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
