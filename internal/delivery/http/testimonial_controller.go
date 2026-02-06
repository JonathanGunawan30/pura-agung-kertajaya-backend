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

type TestimonialController struct {
	UseCase usecase.TestimonialUsecase
	Log     *logrus.Logger
}

func NewTestimonialController(usecase usecase.TestimonialUsecase, log *logrus.Logger) *TestimonialController {
	return &TestimonialController{
		UseCase: usecase,
		Log:     log,
	}
}

func (c *TestimonialController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
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

func (c *TestimonialController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch testimonials")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) GetAllPublic(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetPublic()
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to fetch public testimonials")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	data, err := c.UseCase.GetByID(id)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("testimonial_id", id).Warn("testimonial not found")
		} else {
			c.getLogger(ctx).WithField("testimonial_id", id).WithError(err).Error("failed to get testimonial by id")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) Create(ctx *fiber.Ctx) error {
	var req model.TestimonialRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code < fiber.StatusInternalServerError {
			c.getLogger(ctx).WithField("payload", req).Warnf("failed to create testimonial: %s", e.Message)
		} else {
			c.getLogger(ctx).WithField("payload", req).WithError(err).Error("failed to create testimonial")
		}
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.TestimonialRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.getLogger(ctx).Warnf("invalid request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("testimonial_id", id).Warn("attempted update on non-existent testimonial")
		} else {
			c.getLogger(ctx).WithFields(logrus.Fields{
				"testimonial_id": id,
				"payload":        req,
			}).WithError(err).Error("failed to update testimonial")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).WithField("testimonial_id", id).Warn("attempted delete non-existent testimonial")
		} else {
			c.getLogger(ctx).WithField("testimonial_id", id).WithError(err).Error("failed to delete testimonial")
		}
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Testimonial deleted successfully"})
}
