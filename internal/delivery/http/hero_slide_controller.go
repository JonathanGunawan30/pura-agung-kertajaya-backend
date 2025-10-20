package http

import (
    "pura-agung-kertajaya-backend/internal/model"
    "pura-agung-kertajaya-backend/internal/usecase"

    "github.com/gofiber/fiber/v2"
    "github.com/sirupsen/logrus"
)

type HeroSlideController struct {
    UseCase usecase.HeroSlideUsecase
    Log     *logrus.Logger
}

func NewHeroSlideController(usecase usecase.HeroSlideUsecase, log *logrus.Logger) *HeroSlideController {
    return &HeroSlideController{
        UseCase: usecase,
        Log:     log,
    }
}

func (c *HeroSlideController) GetAll(ctx *fiber.Ctx) error {
    data, err := c.UseCase.GetAll()
    if err != nil {
        c.Log.WithError(err).Error("failed to fetch hero slides")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *HeroSlideController) GetByID(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
    }

    data, err := c.UseCase.GetByID(id)
    if err != nil {
        c.Log.WithError(err).Error("failed to get hero slide by id")
        return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *HeroSlideController) Create(ctx *fiber.Ctx) error {
    var req model.HeroSlideRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
    }

    data, err := c.UseCase.Create(req)
    if err != nil {
        c.Log.WithError(err).Error("failed to create hero slide")
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *HeroSlideController) Update(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
    }

    var req model.HeroSlideRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
    }

    data, err := c.UseCase.Update(id, req)
    if err != nil {
        c.Log.WithError(err).Error("failed to update hero slide")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *HeroSlideController) Delete(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
    }

    if err := c.UseCase.Delete(id); err != nil {
        c.Log.WithError(err).Error("failed to delete hero slide")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[string]{Data: "Hero slide deleted successfully"})
}
