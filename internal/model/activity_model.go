package model

import "time"

// ActivityRequest defines the payload for creating or updating an activity
type ActivityRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"required"`
	TimeInfo    string `json:"time_info" validate:"omitempty,max=100"`
	Location    string `json:"location" validate:"omitempty,max=100"`
	OrderIndex  int    `json:"order_index" validate:"required,min=1"`
	IsActive    bool   `json:"is_active" validate:"boolean"`
}

// ActivityResponse defines the API response for an activity
type ActivityResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TimeInfo    string    `json:"time_info"`
	Location    string    `json:"location"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
