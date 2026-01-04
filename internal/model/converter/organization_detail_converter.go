package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToOrganizationDetailResponse(od *entity.OrganizationDetail) model.OrganizationDetailResponse {
	return model.OrganizationDetailResponse{
		ID:                    od.ID,
		EntityType:            od.EntityType,
		Vision:                od.Vision,
		Mission:               od.Mission,
		Rules:                 od.Rules,
		WorkProgram:           od.WorkProgram,
		VisionMissionImageURL: od.VisionMissionImageURL,
		WorkProgramImageURL:   od.WorkProgramImageURL,
		RulesImageURL:         od.RulesImageURL,
		StructureImageURL:     od.StructureImageURL,
		CreatedAt:             od.CreatedAt,
		UpdatedAt:             od.UpdatedAt,
	}
}

func ToOrganizationDetailResponses(details []entity.OrganizationDetail) []model.OrganizationDetailResponse {
	var responses []model.OrganizationDetailResponse
	for _, detail := range details {
		responses = append(responses, ToOrganizationDetailResponse(&detail))
	}
	return responses
}
