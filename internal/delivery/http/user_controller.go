package http

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase usecase.UserUseCase
	Config  *viper.Viper
}

func NewUserController(useCase usecase.UserUseCase, logger *logrus.Logger, config *viper.Viper) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
		Config:  config,
	}
}

func (c *UserController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
	user := middleware.GetUser(ctx)

	userID := "guest"
	userRole := "unknown"

	if user != nil {
		userID = user.ID
		userRole = user.Role
	}

	return c.Log.WithFields(logrus.Fields{
		"user_id":   userID,
		"user_role": userRole,
		"ip":        ctx.IP(),
		"req_id":    ctx.Get("X-Request-ID"),
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	req := new(model.LoginUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.getLogger(ctx).Warnf("Failed to parse request body: %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{Errors: "Invalid request body format"})
	}

	response, tokenString, err := c.UseCase.Login(ctx.UserContext(), req)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) {
			c.getLogger(ctx).Warnf("Login failed: %s", e.Message)
		} else {
			c.getLogger(ctx).WithError(err).Error("Login system error")
		}
		return err
	}

	domain := c.Config.GetString("cookie.domain")

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
		Domain:   domain,
		MaxAge:   86400,
	})

	c.getLogger(ctx).WithFields(logrus.Fields{
		"user_id":   response.ID,
		"user_role": response.Role,
	}).Info("User logged in successfully")

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	tokenString := ctx.Cookies("access_token")

	err := c.UseCase.Logout(ctx.UserContext(), tokenString)
	if err != nil {
		c.getLogger(ctx).WithError(err).Warn("Failed to revoke token during logout")
	}

	ctx.ClearCookie("access_token")

	c.getLogger(ctx).Info("User logged out")
	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	if auth == nil {
		c.getLogger(ctx).Warn("Auth user not found in context")
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[any]{Errors: fiber.ErrUnauthorized.Message})
	}

	response, err := c.UseCase.Current(ctx.UserContext(), auth.ID)
	if err != nil {
		var e *model.ResponseError
		if errors.As(err, &e) && e.Code == fiber.StatusNotFound {
			c.getLogger(ctx).Warn("Current user not found in DB")
		} else {
			c.getLogger(ctx).WithError(err).Error("Failed to get current user")
		}
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	if auth == nil {
		c.getLogger(ctx).Warn("Auth user not found in context for update")
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[any]{Errors: fiber.ErrUnauthorized.Message})
	}

	req := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.getLogger(ctx).Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.UpdateProfile(ctx.UserContext(), auth.ID, req)
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("Failed to update profile")
		return err
	}

	c.getLogger(ctx).Info("User profile updated")
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
