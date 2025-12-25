package model

import "time"

type CreateFacilityRequest struct {
	EntityType  string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" validate:"required,url|uri|startswith=http"`
	OrderIndex  int    `json:"order_index"`
	IsActive    bool   `json:"is_active"`
}

type UpdateFacilityRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" validate:"required,url|uri|startswith=http"`
	OrderIndex  int    `json:"order_index"`
	IsActive    bool   `json:"is_active"`
}

type FacilityResponse struct {
	ID          string    `json:"id"`
	EntityType  string    `json:"entity_type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
