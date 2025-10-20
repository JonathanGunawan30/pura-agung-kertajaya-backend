package entity

import "time"

// AboutValue represents a single value point under an AboutSection
// Mirrors table: about_values

type AboutValue struct {
    ID         string    `gorm:"column:id;primaryKey;type:varchar(100)"`
    AboutID    string    `gorm:"column:about_id;type:varchar(100);index"`
    Title      string    `gorm:"column:title;type:varchar(100);not null"`
    Value      string    `gorm:"column:value;type:varchar(100);not null"`
    OrderIndex int       `gorm:"column:order_index;default:1"`
    CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AboutValue) TableName() string { return "about_values" }
