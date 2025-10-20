package model

import "time"

type GalleryRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" validate:"required,url|uri|startswith=http"`
	OrderIndex  int    `json:"order_index"`
	IsActive    bool   `json:"is_active"`
}

type GalleryResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
