package http

import (
	"errors"
	"fmt"
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

func (c *OrganizationDetailController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
	user := middleware.GetUser(ctx)

	userID := "guest"
	userRole := "unknown"

	if user != nil {
		userID = fmt.Sprintf("%d", user.ID)
		userRole = user.Role
	}

	return c.Log.WithFields(logrus.Fields{
		"user_id":   userID,
		"user_role": userRole,
		"ip":        ctx.IP(),
		"req_id":    ctx.Get("X-Request-ID"),
	})
}

func (c *OrganizationDetailController) GetPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")

	data, err := c.UseCase.GetByEntityType(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch public organization details")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationDetailController) GetAdmin(ctx *fiber.Ctx) error {
	val := ctx.Locals(middleware.CtxEntityType)
	entityType, ok := val.(string)
	if !ok {
		c.getLogger(ctx).Error("entity_type missing from context locals")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: "Internal Configuration Error"})
	}

	data, err := c.UseCase.GetByEntityType(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch organization details for admin")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *OrganizationDetailController) Update(ctx *fiber.Ctx) error {
	val := ctx.Locals(middleware.CtxEntityType)
	entityType, ok := val.(string)
	if !ok {
		c.getLogger(ctx).Error("entity_type missing from context locals during update")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: "Internal Configuration Error"})
	}

	var req model.UpdateOrganizationDetailRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(entityType, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code < fiber.StatusInternalServerError {
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to update organization details: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to update organization details")
		}
		return err
	}

	c.getLogger(ctx).Info("organization details updated successfully")
	return ctx.JSON(model.WebResponse[any]{Data: data})
}
