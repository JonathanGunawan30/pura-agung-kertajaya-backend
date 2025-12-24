package model

import "time"

type RemarkResponse struct {
	ID         string    `json:"id"`
	EntityType string    `json:"entity_type"`
	Name       string    `json:"name"`
	Position   string    `json:"position"`
	ImageURL   string    `json:"image_url"`
	Content    string    `json:"content"`
	IsActive   bool      `json:"is_active"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateRemarkRequest struct {
	EntityType string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Name       string `json:"name" validate:"required,max=100"`
	Position   string `json:"position" validate:"required,max=100"`
	ImageURL   string `json:"image_url"`
	Content    string `json:"content" validate:"required"`
	IsActive   bool   `json:"is_active"`
	OrderIndex int    `json:"order_index"`
}

type UpdateRemarkRequest struct {
	Name       string `json:"name" validate:"required,max=100"`
	Position   string `json:"position" validate:"required,max=100"`
	ImageURL   string `json:"image_url"`
	Content    string `json:"content" validate:"required"`
	IsActive   bool   `json:"is_active"`
	OrderIndex int    `json:"order_index"`
}
