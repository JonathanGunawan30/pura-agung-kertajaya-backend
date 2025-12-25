package model

import "time"

type CreateActivityRequest struct {
	EntityType  string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"required"`
	TimeInfo    string `json:"time_info" validate:"omitempty,max=100"`
	Location    string `json:"location" validate:"omitempty,max=100"`
	OrderIndex  int    `json:"order_index" validate:"required,min=1"`
	IsActive    bool   `json:"is_active" validate:"boolean"`
}

type UpdateActivityRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"required"`
	TimeInfo    string `json:"time_info" validate:"omitempty,max=100"`
	Location    string `json:"location" validate:"omitempty,max=100"`
	OrderIndex  int    `json:"order_index" validate:"required,min=1"`
	IsActive    bool   `json:"is_active" validate:"boolean"`
}
type ActivityResponse struct {
	ID          string    `json:"id"`
	EntityType  string    `json:"entity_type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TimeInfo    string    `json:"time_info"`
	Location    string    `json:"location"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
