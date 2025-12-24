package entity

import "time"

type Remark struct {
	ID         string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType string    `gorm:"column:entity_type;type:enum('pura','yayasan','pasraman');default:'pura';not null;index"`
	Name       string    `gorm:"column:name;type:varchar(100);not null"`
	Position   string    `gorm:"column:position;type:varchar(100);not null"`
	ImageURL   string    `gorm:"column:image_url;type:text"`
	Content    string    `gorm:"column:content;type:text;not null"`
	IsActive   bool      `gorm:"column:is_active;default:true"`
	OrderIndex int       `gorm:"column:order_index;default:1"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Remark) TableName() string {
	return "remarks"
}
