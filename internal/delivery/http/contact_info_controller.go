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

type ContactInfoController struct {
	UseCase usecase.ContactInfoUsecase
	Log     *logrus.Logger
}

func NewContactInfoController(usecase usecase.ContactInfoUsecase, log *logrus.Logger) *ContactInfoController {
	return &ContactInfoController{UseCase: usecase, Log: log}
}

func (c *ContactInfoController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
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

func (c *ContactInfoController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch contact info")
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
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("contact_id", id).Warn("contact info not found")
		} else {
			c.getLogger(ctx).WithField("contact_id", id).WithError(err).Error("failed to get contact info by id")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) Create(ctx *fiber.Ctx) error {
	var req model.CreateContactInfoRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code < fiber.StatusInternalServerError {
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to create contact info: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to create contact info")
		}
		return err
	}

	c.getLogger(ctx).WithField("contact_id", data.ID).Info("contact info created successfully")
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.UpdateContactInfoRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("contact_id", id).Warn("attempted update on non-existent contact info")
		} else {
			c.getLogger(ctx).WithFields(logrus.Fields{
				"contact_id": id,
				"payload":    req,
			}).WithError(err).Error("failed to update contact info")
		}
		return err
	}

	c.getLogger(ctx).WithField("contact_id", data.ID).Info("contact info updated successfully")
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ContactInfoController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("contact_id", id).Warn("attempted delete non-existent contact info")
		} else {
			c.getLogger(ctx).WithField("contact_id", id).WithError(err).Error("failed to delete contact info")
		}
		return err
	}

	c.getLogger(ctx).WithField("contact_id", id).Info("contact info deleted successfully")
	return ctx.JSON(model.WebResponse[string]{Data: "Contact info deleted successfully"})
}
