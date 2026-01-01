package entity

import (
	"pura-agung-kertajaya-backend/internal/util"
	"time"
)

type Facility struct {
	ID          string        `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType  string        `gorm:"column:entity_type;type:enum('pura','yayasan','pasraman');default:pura';not null;index"`
	Name        string        `gorm:"column:name;type:text;not null"`
	Description string        `gorm:"column:description;type:text"`
	Images      util.ImageMap `gorm:"column:images;type:json"`
	OrderIndex  int           `gorm:"column:order_index;default:1"`
	IsActive    bool          `gorm:"column:is_active;"`
	CreatedAt   time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoUpdateTime"`
}

func (Facility) TableName() string { return "facilities" }
