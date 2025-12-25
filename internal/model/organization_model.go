package model

import "time"

type CreateOrganizationRequest struct {
	EntityType    string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Name          string `json:"name" validate:"required,min=1,max=100"`
	Position      string `json:"position" validate:"required,min=1,max=100"`
	PositionOrder int    `json:"position_order" validate:"required,min=1"`
	OrderIndex    int    `json:"order_index"`
	IsActive      bool   `json:"is_active"`
}

type UpdateOrganizationRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=100"`
	Position      string `json:"position" validate:"required,min=1,max=100"`
	PositionOrder int    `json:"position_order" validate:"required,min=1"`
	OrderIndex    int    `json:"order_index"`
	IsActive      bool   `json:"is_active"`
}

type OrganizationResponse struct {
	ID            string    `json:"id"`
	EntityType    string    `json:"entity_type"`
	Name          string    `json:"name"`
	Position      string    `json:"position"`
	PositionOrder int       `json:"position_order"`
	OrderIndex    int       `json:"order_index"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
