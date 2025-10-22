package http

import (
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ContactInfoController struct {
	UseCase usecase.ContactInfoUsecase
	Log     *logrus.Logger
}

func NewContactInfoController(usecase usecase.ContactInfoUsecase, log *logrus.Logger) *ContactInfoController {
	return &ContactInfoController{UseCase: usecase, Log: log}
}

func (c *ContactInfoController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch contact info")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get contact info by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) Create(ctx *fiber.Ctx) error {
	var req model.ContactInfoRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create contact info")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.ContactInfoRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update contact info")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete contact info")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Contact info deleted successfully"})
}
