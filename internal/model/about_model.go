package model

import "time"

type AboutValueRequest struct {
    Title      string `json:"title" validate:"required,min=1,max=100"`
    Value      string `json:"value" validate:"required,min=1,max=100"`
    OrderIndex int    `json:"order_index"`
}

type AboutSectionRequest struct {
    Title       string               `json:"title" validate:"required,min=1,max=150"`
    Description string               `json:"description" validate:"required"`
    ImageURL    string               `json:"image_url" validate:"omitempty,url|uri|startswith=http"`
    IsActive    bool                 `json:"is_active"`
    Values      []AboutValueRequest  `json:"values" validate:"dive"`
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
    ID          string                `json:"id"`
    Title       string                `json:"title"`
    Description string                `json:"description"`
    ImageURL    string                `json:"image_url"`
    IsActive    bool                  `json:"is_active"`
    CreatedAt   time.Time             `json:"created_at"`
    UpdatedAt   time.Time             `json:"updated_at"`
    Values      []AboutValueResponse  `json:"values"`
}
