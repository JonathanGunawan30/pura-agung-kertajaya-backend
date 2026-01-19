package http

import (
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type GalleryController struct {
	UseCase usecase.GalleryUsecase
	Log     *logrus.Logger
}

func NewGalleryController(usecase usecase.GalleryUsecase, log *logrus.Logger) *GalleryController {
	return &GalleryController{UseCase: usecase, Log: log}
}

func (c *GalleryController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Locals(middleware.CtxEntityType).(string)

	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch galleries")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *GalleryController) GetAllPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public galleries")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *GalleryController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get gallery by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *GalleryController) Create(ctx *fiber.Ctx) error {
	var req model.CreateGalleryRequest
	entityType := ctx.Locals(middleware.CtxEntityType).(string)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(entityType, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create gallery")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *GalleryController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.UpdateGalleryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update gallery")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *GalleryController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete gallery")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Gallery deleted successfully"})
}
