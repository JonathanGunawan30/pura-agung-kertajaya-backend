package http

import (
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

func (c *TestimonialController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.UseCase.GetAll()
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch testimonials")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
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
		c.Log.WithError(err).Error("failed to get testimonial by id")
		return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) Create(ctx *fiber.Ctx) error {
	var req model.TestimonialRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create testimonial")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: err.Error()})
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
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update testimonial")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *TestimonialController) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete testimonial")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Testimonial deleted successfully"})
}
