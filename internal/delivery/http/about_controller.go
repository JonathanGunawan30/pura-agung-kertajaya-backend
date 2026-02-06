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

type AboutController struct {
	UseCase usecase.AboutUsecase
	Log     *logrus.Logger
}

func NewAboutController(usecase usecase.AboutUsecase, log *logrus.Logger) *AboutController {
	return &AboutController{UseCase: usecase, Log: log}
}

func (c *AboutController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
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

func (c *AboutController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch about sections")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) GetAllPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch public about sections")
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
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("about_id", id).Warn("about section not found")
		} else {
			c.getLogger(ctx).WithField("about_id", id).WithError(err).Error("failed to get about section by id")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) Create(ctx *fiber.Ctx) error {
	var req model.AboutSectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code < fiber.StatusInternalServerError {
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to create about section: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to create about section")
		}
		return err
	}

	c.getLogger(ctx).WithField("about_id", data.ID).Info("about section created successfully")
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}
	var req model.AboutSectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("about_id", id).Warn("attempted update on non-existent about section")
		} else {
			c.getLogger(ctx).WithFields(logrus.Fields{
				"about_id": id,
				"payload":  req,
			}).WithError(err).Error("failed to update about section")
		}
		return err
	}

	c.getLogger(ctx).WithField("about_id", data.ID).Info("about section updated successfully")
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *AboutController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("about_id", id).Warn("attempted delete non-existent about section")
		} else {
			c.getLogger(ctx).WithField("about_id", id).WithError(err).Error("failed to delete about section")
		}
		return err
	}

	c.getLogger(ctx).WithField("about_id", id).Info("about section deleted successfully")
	return ctx.JSON(model.WebResponse[string]{Data: "About section deleted successfully"})
}
