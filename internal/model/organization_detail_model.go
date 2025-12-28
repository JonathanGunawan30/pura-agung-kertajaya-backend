package model

import "time"

type OrganizationDetailResponse struct {
	ID                    string `json:"id"`
	EntityType            string `json:"entity_type"`
	Vision                string `json:"vision"`
	Mission               string `json:"mission"`
	Rules                 string `json:"rules"`
	WorkProgram           string `json:"work_program"`
	VisionMissionImageURL string `json:"vision_mission_image_url"`
	WorkProgramImageURL   string `json:"work_program_image_url"`
	RulesImageURL         string `json:"rules_image_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrganizationDetailRequest struct {
	EntityType            string `json:"entity_type" validate:"required,oneof=pura yayasan pasraman"`
	Vision                string `json:"vision"`
	Mission               string `json:"mission"`
	Rules                 string `json:"rules"`
	WorkProgram           string `json:"work_program"`
	VisionMissionImageURL string `json:"vision_mission_image_url"`
	WorkProgramImageURL   string `json:"work_program_image_url"`
	RulesImageURL         string `json:"rules_image_url"`
}

type UpdateOrganizationDetailRequest struct {
	Vision                string `json:"vision"`
	Mission               string `json:"mission"`
	Rules                 string `json:"rules"`
	WorkProgram           string `json:"work_program"`
	VisionMissionImageURL string `json:"vision_mission_image_url"`
	WorkProgramImageURL   string `json:"work_program_image_url"`
	RulesImageURL         string `json:"rules_image_url"`
}
