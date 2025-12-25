package entity

import "time"

// Mirrors table: organization_members

type OrganizationMember struct {
	ID            string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType    string    `gorm:"column:entity_type;type:enum('pura','yayasan','pasraman');default:pura';not null;index"`
	Name          string    `gorm:"column:name;type:varchar(100);not null"`
	Position      string    `gorm:"column:position;type:varchar(100);not null;"`
	PositionOrder int       `gorm:"column:position_order;not null;default:99"`
	OrderIndex    int       `gorm:"column:order_index;default:1"`
	IsActive      bool      `gorm:"column:is_active"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}
