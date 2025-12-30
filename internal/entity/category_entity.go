package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Slug      string    `gorm:"column:slug;type:varchar(100;unique;not null;index"`
	CreatedAt time.Time `gorm:"created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"updated_at;autoUpdateTime"`
}

func (Category) TableName() string {
	return "categories"
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
