package model

import "time"

type FacilityRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" validate:"required,url|uri|startswith=http"`
	OrderIndex  int    `json:"order_index"`
	IsActive    bool   `json:"is_active"`
}

type FacilityResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
