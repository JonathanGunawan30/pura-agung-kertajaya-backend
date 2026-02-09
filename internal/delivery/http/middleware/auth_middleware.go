package middleware

import (
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	ID   string
	Role string
}

func AuthMiddleware(tokenUtil *util.TokenUtil) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		auth, _, err := tokenUtil.ParseToken(c.Context(), token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user", &Auth{
			ID:   auth.ID,
			Role: auth.Role,
		})
		return c.Next()
	}
}

func GetUser(c *fiber.Ctx) *Auth {
	user, ok := c.Locals("user").(*Auth)
	if !ok {
		return nil
	}
	return user
}
