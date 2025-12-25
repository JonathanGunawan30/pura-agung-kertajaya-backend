package entity

import "time"

// SiteIdentity mirrors the DB schema from migrations: site_identity
// id VARCHAR(100) PRIMARY KEY,
// site_name VARCHAR(150) NOT NULL,
// logo_url TEXT,
// tagline VARCHAR(255),
// primary_button_text VARCHAR(50),
// primary_button_link VARCHAR(255),
// secondary_button_text VARCHAR(50),
// secondary_button_link VARCHAR(255),
// created_at TIMESTAMP,
// updated_at TIMESTAMP

type SiteIdentity struct {
	ID                  string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType          string    `gorm:"column:entity_type;type:enum('pura', 'yayasan', 'pasraman');not null;default:'pura'"`
	SiteName            string    `gorm:"column:site_name;type:varchar(150);not null"`
	LogoURL             string    `gorm:"column:logo_url;type:text"`
	Tagline             string    `gorm:"column:tagline;type:varchar(255)"`
	PrimaryButtonText   string    `gorm:"column:primary_button_text;type:varchar(50)"`
	PrimaryButtonLink   string    `gorm:"column:primary_button_link;type:varchar(255)"`
	SecondaryButtonText string    `gorm:"column:secondary_button_text;type:varchar(50)"`
	SecondaryButtonLink string    `gorm:"column:secondary_button_link;type:varchar(255)"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SiteIdentity) TableName() string { return "site_identity" }
