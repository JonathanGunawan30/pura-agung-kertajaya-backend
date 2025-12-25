package entity

import "time"

// ContactInfo represents the contact information for the website
// Mirrors the DB schema from migrations: contact_info
type ContactInfo struct {
	ID            string    `gorm:"column:id;primaryKey;type:varchar(100)"`
	EntityType    string    `gorm:"column:entity_type;type:enum('pura','yayasan','pasraman');default:pura';not null;index"`
	Address       string    `gorm:"column:address;type:text;not null"`
	Phone         string    `gorm:"column:phone;type:varchar(50)"`
	Email         string    `gorm:"column:email;type:varchar(100)"`
	VisitingHours string    `gorm:"column:visiting_hours;type:varchar(100)"`
	MapEmbedURL   string    `gorm:"column:map_embed_url;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ContactInfo) TableName() string {
	return "contact_info"
}
