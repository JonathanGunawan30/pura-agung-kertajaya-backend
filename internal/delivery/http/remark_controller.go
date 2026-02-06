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

type RemarkController struct {
	UseCase usecase.RemarkUsecase
	Log     *logrus.Logger
}

func NewRemarkController(usecase usecase.RemarkUsecase, log *logrus.Logger) *RemarkController {
	return &RemarkController{
		UseCase: usecase,
		Log:     log,
	}
}

func (c *RemarkController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
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

func (c *RemarkController) GetAll(ctx *fiber.Ctx) error {
	val := ctx.Locals(middleware.CtxEntityType)
	entityType, ok := val.(string)
	if !ok {
		c.getLogger(ctx).Error("entity_type missing from context locals")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: "Internal Configuration Error"})
	}

	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch remarks")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) GetAllPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")

	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch public remarks")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	data, err := c.UseCase.GetByID(id)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("remark_id", id).Warn("remark not found")
		} else {
			c.getLogger(ctx).WithField("remark_id", id).WithError(err).Error("failed to get remark by id")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) Create(ctx *fiber.Ctx) error {
	var req model.CreateRemarkRequest

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
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to create remark: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to create remark")
		}
		return err
	}

	c.getLogger(ctx).WithField("remark_id", data.ID).Info("remark created successfully")
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.UpdateRemarkRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("remark_id", id).Warn("attempted update on non-existent remark")
		} else {
			c.getLogger(ctx).WithFields(logrus.Fields{
				"remark_id": id,
				"payload":   req,
			}).WithError(err).Error("failed to update remark")
		}
		return err
	}

	c.getLogger(ctx).WithField("remark_id", data.ID).Info("remark updated successfully")
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("remark_id", id).Warn("attempted delete non-existent remark")
		} else {
			c.getLogger(ctx).WithField("remark_id", id).WithError(err).Error("failed to delete remark")
		}
		return err
	}

	c.getLogger(ctx).WithField("remark_id", id).Info("remark deleted successfully")
	return ctx.JSON(model.WebResponse[string]{Data: "Remark deleted successfully"})
}
