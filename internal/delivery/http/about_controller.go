package http

import (
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AboutController struct {
	UseCase usecase.AboutUsecase
	Log     *logrus.Logger
}

func NewAboutController(usecase usecase.AboutUsecase, log *logrus.Logger) *AboutController {
	return &AboutController{UseCase: usecase, Log: log}
}

func (c *AboutController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch about sections")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

// GetAllPublic returns only active about sections with values
func (c *AboutController) GetAllPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public about sections")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get about by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) Create(ctx *fiber.Ctx) error {
	var req model.AboutSectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create about section")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.AboutSectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update about section")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete about section")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "About section deleted successfully"})
}
