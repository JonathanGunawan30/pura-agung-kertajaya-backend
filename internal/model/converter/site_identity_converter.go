package converter

import (
    "pura-agung-kertajaya-backend/internal/entity"
    "pura-agung-kertajaya-backend/internal/model"
)

func ToSiteIdentityResponse(e entity.SiteIdentity) model.SiteIdentityResponse {
    return model.SiteIdentityResponse{
        ID:                  e.ID,
        SiteName:            e.SiteName,
        LogoURL:             e.LogoURL,
        Tagline:             e.Tagline,
        PrimaryButtonText:   e.PrimaryButtonText,
        PrimaryButtonLink:   e.PrimaryButtonLink,
        SecondaryButtonText: e.SecondaryButtonText,
        SecondaryButtonLink: e.SecondaryButtonLink,
        CreatedAt:           e.CreatedAt,
        UpdatedAt:           e.UpdatedAt,
    }
}
