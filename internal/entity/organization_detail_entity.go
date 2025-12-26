package entity

import "time"

type OrganizationDetail struct {
	ID          string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType  string    `gorm:"column:entity_type;type:enum('pura','yayasan','pasraman');not null;unique"`
	Vision      string    `gorm:"column:vision;type:longtext"`
	Mission     string    `gorm:"column:mission;type:longtext"`
	Rules       string    `gorm:"column:rules;type:longtext"`
	WorkProgram string    `gorm:"column:work_program;type:longtext"`
	ImageURL    string    `gorm:"column:image_url;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (OrganizationDetail) TableName() string {
	return "organization_details"
}
