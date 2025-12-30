package http

import (
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CategoryController struct {
	UseCase usecase.CategoryUsecase
	Log     *logrus.Logger
}

func NewCategoryController(usecase usecase.CategoryUsecase, log *logrus.Logger) *CategoryController {
	return &CategoryController{UseCase: usecase, Log: log}
}

func (c *CategoryController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch categories")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) GetAllPublic(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public categories")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get Category by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	var req model.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create Category")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.UpdateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update Category")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete Category")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Category deleted successfully"})
}
