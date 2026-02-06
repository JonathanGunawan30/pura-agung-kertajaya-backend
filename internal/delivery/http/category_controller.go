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

type CategoryController struct {
	UseCase usecase.CategoryUsecase
	Log     *logrus.Logger
}

func NewCategoryController(usecase usecase.CategoryUsecase, log *logrus.Logger) *CategoryController {
	return &CategoryController{UseCase: usecase, Log: log}
}

func (c *CategoryController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
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

func (c *CategoryController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch categories")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) GetAllPublic(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch public categories")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	data, err := c.UseCase.GetByID(id)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("category_id", id).Warn("category not found")
		} else {
			c.getLogger(ctx).WithField("category_id", id).WithError(err).Error("failed to get category by id")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	var req model.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code < fiber.StatusInternalServerError {
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to create category: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to create category")
		}
		return err
	}

	c.getLogger(ctx).WithField("category_id", data.ID).Info("category created successfully")
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.UpdateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("category_id", id).Warn("attempted update on non-existent category")
		} else {
			c.getLogger(ctx).WithFields(logrus.Fields{
				"category_id": id,
				"payload":     req,
			}).WithError(err).Error("failed to update category")
		}
		return err
	}

	c.getLogger(ctx).WithField("category_id", data.ID).Info("category updated successfully")
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) {
			if e.Code == fiber.StatusNotFound {
				c.getLogger(ctx).WithField("category_id", id).Warn("attempted delete non-existent category")
			} else if e.Code == fiber.StatusConflict {
				c.getLogger(ctx).WithField("category_id", id).Warn("prevented deletion of category in use")
			} else {
				c.getLogger(ctx).WithField("category_id", id).Warnf("business error during delete: %s", e.Message)
			}
		} else {
			c.getLogger(ctx).WithField("category_id", id).WithError(err).Error("failed to delete category")
		}
		return err
	}

	c.getLogger(ctx).WithField("category_id", id).Info("category deleted successfully")
	return ctx.JSON(model.WebResponse[string]{Data: "Category deleted successfully"})
}
