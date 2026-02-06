package test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
)

func setupMockArticleUsecase(t *testing.T) (usecase.ArticleUsecase, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub db: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	u := usecase.NewArticleUsecase(gormDB, validator.New())
	return u, mock
}

func TestArticleUsecase_GetPublic(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)

	rows := sqlmock.NewRows([]string{"id", "title", "status", "published_at", "images"}).
		AddRow("uuid-1", "Berita 1", "PUBLISHED", time.Now(), []byte(`{"lg":"img1.jpg"}`)).
		AddRow("uuid-2", "Berita 2", "PUBLISHED", time.Now(), []byte(`{"lg":"img2.jpg"}`))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE status = ? ORDER BY is_featured DESC, published_at DESC LIMIT ?")).
		WithArgs(entity.ArticleStatusPublished, 10).
		WillReturnRows(rows)

	results, err := u.GetPublic(10)

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	if len(results) > 0 {
		assert.Equal(t, "Berita 1", results[0].Title)
		assert.Equal(t, "img1.jpg", results[0].Images.Lg)
	}
}

func TestArticleUsecase_GetBySlug(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)

	slug := "upacara-ngaben"

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "status", "images"}).
		AddRow("uuid-1", "Upacara Ngaben", slug, "PUBLISHED", []byte(`{"lg":"img1.jpg"}`))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE slug = ? AND status = ? ORDER BY `articles`.`id` LIMIT ?")).
		WithArgs(slug, entity.ArticleStatusPublished, 1).
		WillReturnRows(rows)

	res, err := u.GetBySlug(slug)

	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Upacara Ngaben", res.Title)
		assert.Equal(t, "img1.jpg", res.Images.Lg)
	}
}

func TestArticleUsecase_GetBySlug_NotFound(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)
	slug := "missing-slug"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE slug = ? AND status = ? ORDER BY `articles`.`id` LIMIT ?")).
		WithArgs(slug, entity.ArticleStatusPublished, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetBySlug(slug)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "article not found", e.Message)
	}
}

func TestArticleUsecase_GetByID_NotFound(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)
	id := "missing-id"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE id = ? ORDER BY `articles`.`id` LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	res, err := u.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, res)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "article not found", e.Message)
	}
}

func TestArticleUsecase_Create_AutoSlugAndExcerpt(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)

	req := model.CreateArticleRequest{
		Title:      "Judul Berita Keren",
		AuthorName: "Admin",
		Content:    "Ini adalah konten yang sangat panjang sekali...",
		Excerpt:    "Ini adalah konten yang sangat panjang sekali...",
		Status:     "PUBLISHED",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE slug = ?")).
		WithArgs("judul-berita-keren").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `articles`")).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			req.Title,
			"judul-berita-keren",
			req.AuthorName,
			"",
			"Ini adalah konten yang sangat panjang sekali...",
			req.Content,
			sqlmock.AnyArg(),
			"PUBLISHED",
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req)

	assert.NoError(t, err)
	if assert.NotNil(t, created) {
		assert.Equal(t, "judul-berita-keren", created.Slug)
		assert.Equal(t, req.Content, created.Excerpt)
		assert.Equal(t, req.Images["lg"], created.Images.Lg)
	}
}

func TestArticleUsecase_Create_SlugCollision(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)

	req := model.CreateArticleRequest{
		Title:      "Berita Sama",
		AuthorName: "Budi",
		Content:    "Isi konten ini harus cukup panjang ya",
		Excerpt:    "Isi konten ini harus cukup panjang ya",
		Status:     "DRAFT",
		Images:     map[string]string{"lg": "https://img.com/lg.jpg"},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE slug = ?")).
		WithArgs("berita-sama").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE slug = ?")).
		WithArgs("berita-sama-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `articles`")).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			req.Title,
			"berita-sama-1",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			"Isi konten ini harus cukup panjang ya",
			"Isi konten ini harus cukup panjang ya",
			sqlmock.AnyArg(),
			"DRAFT",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := u.Create(req)

	assert.NoError(t, err)
	if assert.NotNil(t, created) {
		assert.Equal(t, "berita-sama-1", created.Slug)
	}
}

func TestArticleUsecase_Update(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)
	id := "art-1"

	req := model.UpdateArticleRequest{
		Title:      "Judul Baru",
		Content:    "Konten baru",
		Excerpt:    "Konten baru",
		AuthorName: "Author",
		Status:     "PUBLISHED",
		Images:     map[string]string{"lg": "https://new.jpg"},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "images"}).
			AddRow(id, "Judul Lama", "judul-lama", []byte(`{"lg":"old.jpg"}`)))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE slug = ? AND id != ?")).
		WithArgs("judul-baru", id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `articles`")).
		WithArgs(
			sqlmock.AnyArg(),
			"Judul Baru",
			"judul-baru",
			req.AuthorName,
			"",
			"Konten baru",
			req.Content,
			sqlmock.AnyArg(),
			"PUBLISHED",
			false,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := u.Update(id, req)
	assert.NoError(t, err)

	if assert.NotNil(t, updated) {
		assert.Equal(t, "Judul Baru", updated.Title)
		assert.Equal(t, "judul-baru", updated.Slug)
		assert.Equal(t, "https://new.jpg", updated.Images.Lg)
	}
}

func TestArticleUsecase_Update_NotFound(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)
	id := "missing-id"

	req := model.UpdateArticleRequest{
		Title:      "Valid Title",
		Content:    "Ini konten yang panjangnya harus mencukupi sesuai aturan validasi min tag",
		Excerpt:    "Ini excerpt valid",
		AuthorName: "Valid Author",
		Status:     "DRAFT",
		Images:     map[string]string{"lg": "https://img.com/valid.jpg"},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE id = ? LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	updated, err := u.Update(id, req)

	assert.Error(t, err)
	assert.Nil(t, updated)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "article not found", e.Message)
	}
}

func TestArticleUsecase_Delete(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)
	id := "del-1"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE id = ?")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `articles` WHERE `articles`.`id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := u.Delete(id)
	assert.NoError(t, err)
}

func TestArticleUsecase_Delete_NotFound(t *testing.T) {
	u, mock := setupMockArticleUsecase(t)
	id := "missing-id"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `articles` WHERE id = ?")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err := u.Delete(id)
	assert.Error(t, err)

	var e *model.ResponseError
	if assert.ErrorAs(t, err, &e) {
		assert.Equal(t, 404, e.Code)
		assert.Equal(t, "article not found", e.Message)
	}
}
