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

type SiteIdentityController struct {
	UseCase usecase.SiteIdentityUsecase
	Log     *logrus.Logger
}

func NewSiteIdentityController(usecase usecase.SiteIdentityUsecase, log *logrus.Logger) *SiteIdentityController {
	return &SiteIdentityController{UseCase: usecase, Log: log}
}

func (c *SiteIdentityController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
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

func (c *SiteIdentityController) GetAll(ctx *fiber.Ctx) error {
	val := ctx.Locals(middleware.CtxEntityType)
	entityType, ok := val.(string)
	if !ok {
		c.getLogger(ctx).Error("entity_type missing from context locals")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: "Internal Configuration Error"})
	}

	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch site identities")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) GetPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).Warn("public site identity not found")
		} else {
			c.getLogger(ctx).WithError(err).Error("failed to fetch public site identity")
		}
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
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("site_id", id).Warn("site identity not found")
		} else {
			c.getLogger(ctx).WithField("site_id", id).WithError(err).Error("failed to get site identity by id")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) Create(ctx *fiber.Ctx) error {
	var req model.SiteIdentityRequest

	val := ctx.Locals(middleware.CtxEntityType)
	entityType, ok := val.(string)
	if !ok {
		c.getLogger(ctx).Error("entity_type missing from context locals during create")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: "Internal Configuration Error"})
	}

	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(entityType, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code < fiber.StatusInternalServerError {
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to create site identity: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to create site identity")
		}
		return err
	}

	c.getLogger(ctx).WithField("site_id", data.ID).Info("site identity created successfully")
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.SiteIdentityRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("site_id", id).Warn("attempted update on non-existent site identity")
		} else {
			c.getLogger(ctx).WithFields(logrus.Fields{
				"site_id": id,
				"payload": req,
			}).WithError(err).Error("failed to update site identity")
		}
		return err
	}

	c.getLogger(ctx).WithField("site_id", data.ID).Info("site identity updated successfully")
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *SiteIdentityController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("site_id", id).Warn("attempted delete non-existent site identity")
		} else {
			c.getLogger(ctx).WithField("site_id", id).WithError(err).Error("failed to delete site identity")
		}
		return err
	}

	c.getLogger(ctx).WithField("site_id", id).Info("site identity deleted successfully")
	return ctx.JSON(model.WebResponse[string]{Data: "Site identity deleted successfully"})
}
