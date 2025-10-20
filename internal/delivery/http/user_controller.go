package http

import (
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	req := new(model.LoginUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), req, ctx)
	if err != nil {
		c.Log.Warnf("Failed to login user: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	response, err := c.UseCase.Current(ctx.UserContext(), auth.ID)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	success, err := c.UseCase.Logout(ctx.UserContext(), ctx)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to logout user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: success})
}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	if auth == nil {
		return fiber.ErrUnauthorized
	}

	req := new(model.UpdateUserRequest)
	req.ID = auth.ID
	if err := ctx.BodyParser(req); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.UpdateProfile(ctx.UserContext(), auth.ID, req)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to update profile")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
