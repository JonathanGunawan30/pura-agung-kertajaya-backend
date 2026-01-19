package http

import (
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SiteIdentityController struct {
	UseCase usecase.SiteIdentityUsecase
	Log     *logrus.Logger
}

func NewSiteIdentityController(usecase usecase.SiteIdentityUsecase, log *logrus.Logger) *SiteIdentityController {
	return &SiteIdentityController{UseCase: usecase, Log: log}
}

func (c *SiteIdentityController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Locals(middleware.CtxEntityType).(string)
	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch site identities")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

// GetPublic returns the latest site identity for public consumption
func (c *SiteIdentityController) GetPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public site identity")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get site identity by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) Create(ctx *fiber.Ctx) error {
	var req model.SiteIdentityRequest
	entityType := ctx.Locals(middleware.CtxEntityType).(string)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(entityType, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create site identity")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.SiteIdentityRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update site identity")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete site identity")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Site identity deleted successfully"})
}
