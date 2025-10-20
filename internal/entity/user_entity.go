package entity

import "time"

// User represents an admin user in the CMS
type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Email     string    `gorm:"column:email;size:100;unique;not null"`
	Password  string    `gorm:"column:password;size:100;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
