package entity

import "time"

type Gallery struct {
	ID          string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType  string    `gorm:"column:entity_type;type:enum('pura','yayasan','pasraman');default:pura';not null;index"`
	Title       string    `gorm:"column:title;type:varchar(150);not null"`
	Description string    `gorm:"column:description;type:text"`
	ImageURL    string    `gorm:"column:image_url;type:text;not null"`
	OrderIndex  int       `gorm:"column:order_index;default:1"`
	IsActive    bool      `gorm:"column:is_active;"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Gallery) TableName() string { return "galleries" }
