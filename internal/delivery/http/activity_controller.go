package http

import (
    "pura-agung-kertajaya-backend/internal/model"
    "pura-agung-kertajaya-backend/internal/usecase"

    "github.com/gofiber/fiber/v2"
    "github.com/sirupsen/logrus"
)

type ActivityController struct {
    UseCase usecase.ActivityUsecase
    Log     *logrus.Logger
}

func NewActivityController(usecase usecase.ActivityUsecase, log *logrus.Logger) *ActivityController {
    return &ActivityController{UseCase: usecase, Log: log}
}

func (c *ActivityController) GetAll(ctx *fiber.Ctx) error {
    data, err := c.UseCase.GetAll()
    if err != nil {
        c.Log.WithError(err).Error("failed to fetch activities")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

// GetAllPublic returns only active activities for public consumption
func (c *ActivityController) GetAllPublic(ctx *fiber.Ctx) error {
    data, err := c.UseCase.GetPublic()
    if err != nil {
        c.Log.WithError(err).Error("failed to fetch public activities")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ActivityController) GetByID(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
    }
    data, err := c.UseCase.GetByID(id)
    if err != nil {
        c.Log.WithError(err).Error("failed to get activity by id")
        return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ActivityController) Create(ctx *fiber.Ctx) error {
    var req model.ActivityRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
    }
    data, err := c.UseCase.Create(req)
    if err != nil {
        c.Log.WithError(err).Error("failed to create activity")
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *ActivityController) Update(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
    }
    var req model.ActivityRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
    }
    data, err := c.UseCase.Update(id, req)
    if err != nil {
        c.Log.WithError(err).Error("failed to update activity")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *ActivityController) Delete(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
    }
    if err := c.UseCase.Delete(id); err != nil {
        c.Log.WithError(err).Error("failed to delete activity")
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{Errors: err.Error()})
    }
    return ctx.JSON(model.WebResponse[string]{Data: "Activity deleted successfully"})
}
