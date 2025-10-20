package entity

import (
	"time"
)

type Testimonial struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	AvatarURL  string    `gorm:"type:text" json:"avatar_url"`
	Rating     int       `gorm:"type:int;not null;check:rating>=1 AND rating<=5" json:"rating" validate:"required,min=1,max=5"`
	Comment    string    `gorm:"type:text;not null" json:"comment" validate:"required"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	OrderIndex int       `gorm:"default:1" json:"order_index"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Testimonial) TableName() string {
	return "testimonials"
}
