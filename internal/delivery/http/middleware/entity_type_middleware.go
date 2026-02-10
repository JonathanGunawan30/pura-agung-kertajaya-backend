package middleware

import "github.com/gofiber/fiber/v2"

const CtxEntityType = "entity_type"

type entityTypeBody struct {
	EntityType string `json:"entity_type"`
}

func EntityTypeMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := GetUser(ctx)
		if user == nil {
			return fiber.ErrUnauthorized
		}

		entityType := ctx.Query("entity_type")

		if entityType == "" && (ctx.Method() == fiber.MethodPost || ctx.Method() == fiber.MethodPut || ctx.Method() == fiber.MethodPatch) {
			var payloadBody entityTypeBody
			if err := ctx.BodyParser(&payloadBody); err == nil {
				entityType = payloadBody.EntityType
			}
		}

		valid := map[string]bool{
			"pura":     true,
			"yayasan":  true,
			"pasraman": true,
		}

		if user.Role == "super" {
			if !valid[entityType] {
				entityType = "pura"
			}

			ctx.Locals(CtxEntityType, entityType)
			return ctx.Next()
		}

		ctx.Locals(CtxEntityType, user.Role)
		return ctx.Next()
	}
}
