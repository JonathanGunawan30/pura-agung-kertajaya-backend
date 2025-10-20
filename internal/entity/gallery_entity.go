package entity

import "time"

// Gallery represents a gallery item
// Mirrors the DB schema from migrations: galleries
// id VARCHAR(100) PRIMARY KEY,
// title VARCHAR(150) NOT NULL,
// description TEXT,
// image_url TEXT NOT NULL,
// order_index INT DEFAULT 1,
// is_active BOOLEAN DEFAULT TRUE,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

type Gallery struct {
	ID          string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	Title       string    `gorm:"column:title;type:varchar(150);not null"`
	Description string    `gorm:"column:description;type:text"`
	ImageURL    string    `gorm:"column:image_url;type:text;not null"`
	OrderIndex  int       `gorm:"column:order_index;default:1"`
	IsActive    bool      `gorm:"column:is_active;"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Gallery) TableName() string { return "galleries" }
