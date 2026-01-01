package entity

import (
	"pura-agung-kertajaya-backend/internal/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "DRAFT"
	ArticleStatusPublished ArticleStatus = "PUBLISHED"
	ArticleStatusArchived  ArticleStatus = "ARCHIVED"
)

type Article struct {
	ID          string        `gorm:"column:id;primaryKey;type:varchar(100)"`
	CategoryID  *string       `gorm:"column:category_id;type:varchar(100)"`
	Category    *Category     `gorm:"foreignKey:CategoryID"`
	Title       string        `gorm:"column:title;type:varchar(255);not null"`
	Slug        string        `gorm:"column:slug;type:varchar(255);unique;not null;index"`
	AuthorName  string        `gorm:"column:author_name;type:varchar(100);not null"`
	AuthorRole  string        `gorm:"column:author_role;type:varchar(100)"`
	Excerpt     string        `gorm:"column:excerpt;type:text"`
	Content     string        `gorm:"column:content;type:longtext"`
	Images      util.ImageMap `gorm:"column:images;types:json"`
	Status      ArticleStatus `gorm:"column:status;type:enum('DRAFT','PUBLISHED','ARCHIVED');default:'DRAFT';index"`
	IsFeatured  bool          `gorm:"column:is_featured;default:false;index"`
	PublishedAt *time.Time    `gorm:"column:published_at"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (a *Article) TableName() string {
	return "articles"
}

func (a *Article) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return
}
