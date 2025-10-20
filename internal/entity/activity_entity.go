package entity

import "time"

// Activity Mirrors the DB schema from migrations: Activities
type Activity struct {
	ID          string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	Title       string    `gorm:"column:title;type:varchar(150);not null"`
	Description string    `gorm:"column:description;type:text;not null"`
	TimeInfo    string    `gorm:"column:time_info;type:varchar(100)"`
	Location    string    `gorm:"column:location;type:varchar(100)"`
	OrderIndex  int       `gorm:"column:order_index;not null;default:1"`
	IsActive    bool      `gorm:"column:is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Activity) TableName() string {
	return "activities"
}
