package model

import "time"

type HeroSlideRequest struct {
	EntityType string            `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Images     map[string]string `json:"images" validate:"required"`
	OrderIndex int               `json:"order_index"`
	IsActive   bool              `json:"is_active"`
}

type HeroSlideResponse struct {
	ID         string            `json:"id"`
	EntityType string            `json:"entity_type"`
	Images     map[string]string `json:"images"`
	OrderIndex int               `json:"order_index"`
	IsActive   bool              `json:"is_active"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
