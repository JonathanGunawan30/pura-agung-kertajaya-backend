package http

import (
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrganizationDetailController struct {
	UseCase usecase.OrganizationDetailUsecase
	Log     *logrus.Logger
}

func NewOrganizationDetailController(usecase usecase.OrganizationDetailUsecase, log *logrus.Logger) *OrganizationDetailController {
	return &OrganizationDetailController{
		UseCase: usecase,
		Log:     log,
	}
}

func (c *OrganizationDetailController) GetPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")

	data, err := c.UseCase.GetByEntityType(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public organization details")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationDetailController) GetAdmin(ctx *fiber.Ctx) error {
	entityType := ctx.Locals(middleware.CtxEntityType).(string)

	data, err := c.UseCase.GetByEntityType(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch organization details for admin")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationDetailController) Update(ctx *fiber.Ctx) error {
	entityType := ctx.Locals(middleware.CtxEntityType).(string)

	var req model.UpdateOrganizationDetailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(entityType, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update organization details")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: data})
}
