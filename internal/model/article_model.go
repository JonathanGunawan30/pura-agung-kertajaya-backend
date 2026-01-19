package model

import "time"

type ArticleResponse struct {
	ID          string            `json:"id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	AuthorName  string            `json:"author_name"`
	AuthorRole  string            `json:"author_role"`
	Excerpt     string            `json:"excerpt"`
	Content     string            `json:"content"`
	Images      ImageVariants     `json:"images"`
	Status      string            `json:"status"`
	IsFeatured  bool              `json:"is_featured"`
	PublishedAt *time.Time        `json:"published_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type CreateArticleRequest struct {
	CategoryID  string            `json:"category_id"`
	Title       string            `json:"title" validate:"required,min=5,max=200"`
	AuthorName  string            `json:"author_name" validate:"required,min=2,max=100"`
	AuthorRole  string            `json:"author_role" validate:"omitempty,max=100"`
	Excerpt     string            `json:"excerpt" validate:"required,max=200"`
	Content     string            `json:"content" validate:"required,min=10"`
	Images      map[string]string `json:"images" validate:"required"`
	IsFeatured  bool              `json:"is_featured"`
	Status      string            `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	PublishedAt *time.Time        `json:"published_at"`
}

type UpdateArticleRequest struct {
	CategoryID  string            `json:"category_id"`
	Title       string            `json:"title" validate:"required,min=5,max=200"`
	AuthorName  string            `json:"author_name" validate:"required,min=2,max=100"`
	AuthorRole  string            `json:"author_role" validate:"omitempty,max=100"`
	Excerpt     string            `json:"excerpt" validate:"required,min=10"`
	Content     string            `json:"content" validate:"required,min=10"`
	Images      map[string]string `json:"images" validate:"required"`
	IsFeatured  bool              `json:"is_featured"`
	Status      string            `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	PublishedAt *time.Time        `json:"published_at"`
}
