package converter

import (
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
)

func ToArticleResponse(a *entity.Article) model.ArticleResponse {
	var categoryResp *model.CategoryResponse
	if a.Category != nil {
		c := ToCategoryResponse(a.Category)
		categoryResp = &c
	}

	return model.ArticleResponse{
		ID:          a.ID,
		Category:    categoryResp,
		Title:       a.Title,
		Slug:        a.Slug,
		AuthorName:  a.AuthorName,
		AuthorRole:  a.AuthorRole,
		Excerpt:     a.Excerpt,
		Content:     a.Content,
		Images:      a.Images,
		Status:      string(a.Status),
		IsFeatured:  a.IsFeatured,
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func ToArticleResponses(articles []entity.Article) []model.ArticleResponse {
	var responses []model.ArticleResponse
	for _, article := range articles {
		responses = append(responses, ToArticleResponse(&article))
	}
	return responses
}
