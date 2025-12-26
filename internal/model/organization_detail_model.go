package model

import "time"

type OrganizationDetailResponse struct {
	ID          string    `json:"id"`
	EntityType  string    `json:"entity_type"`
	Vision      string    `json:"vision"`
	Mission     string    `json:"mission"`
	Rules       string    `json:"rules"`
	WorkProgram string    `json:"work_program"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrganizationDetailRequest struct {
	EntityType  string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Vision      string `json:"vision"`
	Mission     string `json:"mission"`
	Rules       string `json:"rules"`
	WorkProgram string `json:"work_program"`
	ImageURL    string `json:"image_url"`
}

type UpdateOrganizationDetailRequest struct {
	Vision      string `json:"vision"`
	Mission     string `json:"mission"`
	Rules       string `json:"rules"`
	WorkProgram string `json:"work_program"`
	ImageURL    string `json:"image_url"`
}
