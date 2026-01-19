package entity

import "time"

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Email     string    `gorm:"column:email;size:100;unique;not null"`
	Password  string    `gorm:"column:password;size:100;not null"`
	Role      string    `gorm:"column:role;type:enum('pura','yayasan','pasraman','super');not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
