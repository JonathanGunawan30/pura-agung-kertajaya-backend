package http

import (
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrganizationController struct {
	UseCase usecase.OrganizationUsecase
	Log     *logrus.Logger
}

func NewOrganizationController(usecase usecase.OrganizationUsecase, log *logrus.Logger) *OrganizationController {
	return &OrganizationController{UseCase: usecase, Log: log}
}

func (c *OrganizationController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch organization members")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

// GetAllPublic returns only active organization members for public consumption
func (c *OrganizationController) GetAllPublic(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetPublic()
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public organization members")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	data, err := c.UseCase.GetByID(id)
	if err != nil {
		c.Log.WithError(err).Error("failed to get organization member by id")
		return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationController) Create(ctx *fiber.Ctx) error {
	var req model.OrganizationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create organization member")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.OrganizationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}
	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update organization member")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete organization member")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Gallery deleted successfully"})
}
