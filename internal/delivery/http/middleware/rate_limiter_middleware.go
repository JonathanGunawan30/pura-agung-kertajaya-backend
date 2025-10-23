package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
)

// 1. Public endpoints (read-only)
func PublicRateLimiter(storage *redis.Storage) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("public:%v", c.IP())
		},
		LimitReached: limitReachedHandler,
	})
}

// 2. Auth endpoints (login, logout)
func AuthRateLimiter(storage *redis.Storage) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 15 * time.Minute,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("auth:%v", c.IP())
		},
		LimitReached: limitReachedHandler,
	})
}

// 3. CMS Read operations (authenticated)
func CMSReadRateLimiter(storage *redis.Storage) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        300,
		Expiration: 5 * time.Minute,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := c.Locals("user"); userID != nil {
				return fmt.Sprintf("cms_read:%v", userID)
			}
			return fmt.Sprintf("cms_read:%v", c.IP())
		},
		LimitReached: limitReachedHandler,
	})
}

// 4. CMS Write operations (create, update, delete)
func CMSWriteRateLimiter(storage *redis.Storage) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 5 * time.Minute,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := c.Locals("user"); userID != nil {
				return fmt.Sprintf("cms_write:%v", userID)
			}
			return fmt.Sprintf("cms_write:%v", c.IP())
		},
		LimitReached: limitReachedHandler,
	})
}

// 5. Storage/Upload
func StorageRateLimiter(storage *redis.Storage) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Hour,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := c.Locals("user"); userID != nil {
				return fmt.Sprintf("storage:%v", userID)
			}
			return fmt.Sprintf("storage:%v", c.IP())
		},
		LimitReached: limitReachedHandler,
	})
}

// 6. Delete operations
func DeleteRateLimiter(storage *redis.Storage) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30,
		Expiration: 10 * time.Minute,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := c.Locals("user"); userID != nil {
				return fmt.Sprintf("delete:%v", userID)
			}
			return fmt.Sprintf("delete:%v", c.IP())
		},
		LimitReached: limitReachedHandler,
	})
}

func limitReachedHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"data":   nil,
		"errors": "Too many requests. Please try again later.",
	})
}
