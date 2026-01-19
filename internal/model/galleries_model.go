package model

import "time"

type CreateGalleryRequest struct {
	EntityType  string            `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Title       string            `json:"title" validate:"required,min=1,max=150"`
	Description string            `json:"description"`
	Images      map[string]string `json:"images" validate:"required"`
	OrderIndex  int               `json:"order_index"`
	IsActive    bool              `json:"is_active"`
}

type UpdateGalleryRequest struct {
	Title       string            `json:"title" validate:"required,min=1,max=150"`
	Description string            `json:"description"`
	Images      map[string]string `json:"images" validate:"required"`
	OrderIndex  int               `json:"order_index"`
	IsActive    bool              `json:"is_active"`
}

type GalleryResponse struct {
	ID          string        `json:"id"`
	EntityType  string        `json:"entity_type"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Images      ImageVariants `json:"images"`
	OrderIndex  int           `json:"order_index"`
	IsActive    bool          `json:"is_active"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
