package http

import (
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FacilityController struct {
	UseCase usecase.FacilityUsecase
	Log     *logrus.Logger
}

func NewFacilityController(usecase usecase.FacilityUsecase, log *logrus.Logger) *FacilityController {
	return &FacilityController{UseCase: usecase, Log: log}
}

func (c *FacilityController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Locals(middleware.CtxEntityType).(string)
	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch facilities")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *FacilityController) GetAllPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public facilities")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *FacilityController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get Facility by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *FacilityController) Create(ctx *fiber.Ctx) error {
	var req model.CreateFacilityRequest
	entityType := ctx.Locals(middleware.CtxEntityType).(string)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(entityType, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create Facility")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *FacilityController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.UpdateFacilityRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update Facility")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *FacilityController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete Facility")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Facility deleted successfully"})
}
