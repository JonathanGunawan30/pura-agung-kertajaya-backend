package model

import "time"

type AboutValueRequest struct {
	Title      string `json:"title" validate:"required,min=1,max=100"`
	Value      string `json:"value" validate:"required,min=1,max=100"`
	OrderIndex int    `json:"order_index"`
}

type AboutSectionRequest struct {
	EntityType  string              `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Title       string              `json:"title" validate:"required,min=1,max=150"`
	Description string              `json:"description" validate:"required"`
	Images      map[string]string   `json:"images" validate:"required"`
	IsActive    bool                `json:"is_active"`
	Values      []AboutValueRequest `json:"values" validate:"dive"`
}

type AboutValueResponse struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Value      string    `json:"value"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AboutSectionResponse struct {
	ID          string               `json:"id"`
	EntityType  string               `json:"entity_type"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Images      ImageVariants        `json:"images"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Values      []AboutValueResponse `json:"values"`
}
