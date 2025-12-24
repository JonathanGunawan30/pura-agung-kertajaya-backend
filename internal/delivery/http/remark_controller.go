package http

import (
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

func (c *RemarkController) GetAll(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")

	data, err := c.UseCase.GetAll(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch remarks")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) GetAllPublic(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")

	data, err := c.UseCase.GetPublic(entityType)
	if err != nil {
		c.Log.WithError(err).Error("failed to fetch public remarks")
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
		c.Log.WithError(err).Error("failed to get remark by id")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) Create(ctx *fiber.Ctx) error {
	var req model.CreateRemarkRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Create(req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create remark")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	var req model.UpdateRemarkRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body"})
	}

	data, err := c.UseCase.Update(id, req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update remark")
		return err
	}
	return ctx.JSON(model.WebResponse[any]{Data: data})
}

func (c *RemarkController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid ID"})
	}

	if err := c.UseCase.Delete(id); err != nil {
		c.Log.WithError(err).Error("failed to delete remark")
		return err
	}
	return ctx.JSON(model.WebResponse[string]{Data: "Remark deleted successfully"})
}
