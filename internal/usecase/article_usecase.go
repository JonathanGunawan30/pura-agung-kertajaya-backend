package usecase

import (
	"fmt"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticleUsecase interface {
	GetAll(filter string) ([]model.ArticleResponse, error)
	GetPublic(limit int) ([]model.ArticleResponse, error)
	GetByID(id string) (*model.ArticleResponse, error)
	GetBySlug(slug string) (*model.ArticleResponse, error)
	Create(req model.CreateArticleRequest) (*model.ArticleResponse, error)
	Update(id string, req model.UpdateArticleRequest) (*model.ArticleResponse, error)
	Delete(id string) error
}

type articleUsecase struct {
	db       *gorm.DB
	repo     *repository.Repository[entity.Article]
	log      *logrus.Logger
	validate *validator.Validate
}

func NewArticleUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate) ArticleUsecase {
	return &articleUsecase{
		db:       db,
		repo:     &repository.Repository[entity.Article]{DB: db},
		log:      log,
		validate: validate,
	}
}

func (u *articleUsecase) GetPublic(limit int) ([]model.ArticleResponse, error) {
	var articles []entity.Article

	query := u.db.Preload("Category").
		Where("status = ?", entity.ArticleStatusPublished).
		Order("is_featured DESC, published_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := u.repo.FindAll(query, &articles); err != nil {
		return nil, err
	}

	return converter.ToArticleResponses(articles), nil
}

func (u *articleUsecase) GetAll(filter string) ([]model.ArticleResponse, error) {
	var articles []entity.Article

	query := u.db.Preload("Category").Order("created_at DESC")

	if filter != "" {
		query = query.Where("title LIKE ?", "%"+filter+"%")
	}

	if err := u.repo.FindAll(query, &articles); err != nil {
		return nil, err
	}

	return converter.ToArticleResponses(articles), nil
}

func (u *articleUsecase) GetByID(id string) (*model.ArticleResponse, error) {
	var article entity.Article
	query := u.db.Preload("Category").Where("id = ?", id)

	if err := query.First(&article).Error; err != nil {
		return nil, err
	}

	resp := converter.ToArticleResponse(&article)
	return &resp, nil
}

func (u *articleUsecase) GetBySlug(slug string) (*model.ArticleResponse, error) {
	var article entity.Article

	if err := u.db.Preload("Category").
		Where("slug = ? AND status = ?", slug, entity.ArticleStatusPublished).
		First(&article).Error; err != nil {
		return nil, err
	}

	resp := converter.ToArticleResponse(&article)
	return &resp, nil
}

func (u *articleUsecase) Create(req model.CreateArticleRequest) (*model.ArticleResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	baseSlug := slug.Make(req.Title)
	finalSlug := baseSlug
	counter := 1
	for {
		count, err := u.repo.CountBySlug(u.db, finalSlug)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			break
		}
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	var catID *string
	if req.CategoryID != "" {
		catID = &req.CategoryID
	}

	pubTime := req.PublishedAt
	if req.Status == "PUBLISHED" && pubTime == nil {
		now := time.Now()
		pubTime = &now
	}

	article := entity.Article{
		CategoryID:  catID,
		Title:       req.Title,
		Slug:        finalSlug,
		AuthorName:  req.AuthorName,
		AuthorRole:  req.AuthorRole,
		Excerpt:     req.Excerpt,
		Content:     req.Content,
		Images:      util.ImageMap(req.Images),
		Status:      entity.ArticleStatus(req.Status),
		IsFeatured:  req.IsFeatured,
		PublishedAt: pubTime,
	}

	if err := u.repo.Create(u.db, &article); err != nil {
		return nil, err
	}

	resp := converter.ToArticleResponse(&article)
	return &resp, nil
}

func (u *articleUsecase) Update(id string, req model.UpdateArticleRequest) (*model.ArticleResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	var article entity.Article
	if err := u.repo.FindById(u.db, &article, id); err != nil {
		return nil, err
	}

	if article.Title != req.Title {
		baseSlug := slug.Make(req.Title)
		finalSlug := baseSlug
		counter := 1
		for {
			count, err := u.repo.CountBySlugIgnoringID(u.db, finalSlug, id)
			if err != nil {
				return nil, err
			}
			if count == 0 {
				break
			}
			finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
			counter++
		}
		article.Slug = finalSlug
	}

	article.Title = req.Title
	article.AuthorName = req.AuthorName
	article.AuthorRole = req.AuthorRole
	article.Excerpt = req.Excerpt
	article.Content = req.Content
	article.Images = util.ImageMap(req.Images)
	article.Status = entity.ArticleStatus(req.Status)
	article.IsFeatured = req.IsFeatured

	if req.CategoryID != "" {
		article.CategoryID = &req.CategoryID
	} else {
		article.CategoryID = nil
	}

	if req.PublishedAt != nil {
		article.PublishedAt = req.PublishedAt
	}

	if err := u.repo.Update(u.db, &article); err != nil {
		return nil, err
	}

	resp := converter.ToArticleResponse(&article)
	return &resp, nil
}

func (u *articleUsecase) Delete(id string) error {
	var article entity.Article
	if err := u.repo.FindById(u.db, &article, id); err != nil {
		return err
	}
	return u.repo.Delete(u.db, &article)
}
