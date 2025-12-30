package http

import (
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ArticleController struct {
	UseCase usecase.ArticleUsecase
	Log     *logrus.Logger
}

func NewArticleController(usecase usecase.ArticleUsecase, log *logrus.Logger) *ArticleController {
	return &ArticleController{UseCase: usecase, Log: log}
}

func (c *ArticleController) GetPublic(ctx *fiber.Ctx) error {
	limitQuery := ctx.Query("limit", "0")
	limit, _ := strconv.Atoi(limitQuery)

	data, err := c.UseCase.GetPublic(limit)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public articles")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ArticleController) GetBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid Slug"})
	}

	data, err := c.UseCase.GetBySlug(slug)
	if err != nil {
		c.Log.WithError(err).Error("failed to get article by slug")
		return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{Errors: "Article not found"})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ArticleController) GetAll(ctx *fiber.Ctx) error {
	filter := ctx.Query("search")

	data, err := c.UseCase.GetAll(filter)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch articles")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ArticleController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get article by id")
		return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{Errors: "Article not found"})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ArticleController) Create(ctx *fiber.Ctx) error {
	var req model.CreateArticleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create article")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *ArticleController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.UpdateArticleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update article")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ArticleController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete article")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Article deleted successfully"})
}
