package entity

import (
	"pura-agung-kertajaya-backend/internal/util"
	"time"
)

type AboutSection struct {
	ID          string        `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType  string        `gorm:"column:entity_type;type:enum('pura', 'yayasan', 'pasraman');not null;default:'pura'"`
	Title       string        `gorm:"column:title;type:varchar(150);not null"`
	Description string        `gorm:"column:description;type:text;not null"`
	Images      util.ImageMap `gorm:"column:images;type:json"`
	IsActive    bool          `gorm:"column:is_active"`
	CreatedAt   time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoUpdateTime"`

	Values []AboutValue `gorm:"foreignKey:AboutID;references:ID;constraint:OnDelete:CASCADE"`
}

func (AboutSection) TableName() string { return "about_section" }
