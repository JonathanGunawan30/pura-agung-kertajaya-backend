package entity

import "time"

type HeroSlide struct {
	ID         string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType string    `gorm:"column:entity_type;type:enum('pura', 'yayasan', 'pasraman');not null;default:'pura'"`
	ImageUrl   string    `gorm:"column:image_url;type:text;not null"`
	OrderIndex int       `gorm:"column:order_index;not null;default:1"`
	IsActive   bool      `gorm:"column:is_active;default:true"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (HeroSlide) TableName() string {
	return "hero_slides"
}
