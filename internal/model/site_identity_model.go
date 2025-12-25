package model

import "time"

// SiteIdentityRequest defines payload to create/update site identity
// SiteName required; links optional but should be valid URLs when provided
// Keep validation minimal and consistent with other modules

type SiteIdentityRequest struct {
	EntityType          string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	SiteName            string `json:"site_name" validate:"required,min=1,max=150"`
	LogoURL             string `json:"logo_url" validate:"omitempty,url|uri|startswith=http"`
	Tagline             string `json:"tagline" validate:"omitempty,max=255"`
	PrimaryButtonText   string `json:"primary_button_text" validate:"omitempty,max=50"`
	PrimaryButtonLink   string `json:"primary_button_link" validate:"omitempty,url|uri|startswith=http"`
	SecondaryButtonText string `json:"secondary_button_text" validate:"omitempty,max=50"`
	SecondaryButtonLink string `json:"secondary_button_link" validate:"omitempty,url|uri|startswith=http"`
}

type SiteIdentityResponse struct {
	ID                  string    `json:"id"`
	EntityType          string    `json:"entity_type"`
	SiteName            string    `json:"site_name"`
	LogoURL             string    `json:"logo_url"`
	Tagline             string    `json:"tagline"`
	PrimaryButtonText   string    `json:"primary_button_text"`
	PrimaryButtonLink   string    `json:"primary_button_link"`
	SecondaryButtonText string    `json:"secondary_button_text"`
	SecondaryButtonLink string    `json:"secondary_button_link"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
