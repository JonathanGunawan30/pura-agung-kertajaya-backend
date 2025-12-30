package model

import "time"

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateCategoryRequest struct {
	ID   string `json:"-"`
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
