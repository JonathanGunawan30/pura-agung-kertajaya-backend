package route

import (
	"pura-agung-kertajaya-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	UserController        *http.UserController
	StorageController     *http.StorageController
	TestimonialController *http.TestimonialController
	HeroSlideController   *http.HeroSlideController
	AuthMiddleware        fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Post("/api/users/_logout", c.UserController.Logout)
	c.App.Patch("/api/users/_current", c.UserController.UpdateProfile)
	c.App.Get("/api/users/_current", c.UserController.Current)

	c.App.Post("/api/storage/upload", c.StorageController.Upload)
	c.App.Delete("/api/storage/delete", c.StorageController.Delete)
	c.App.Get("/api/storage/presigned-url", c.StorageController.GetPresignedURL)

	c.App.Get("/api/testimonials", c.TestimonialController.GetAll)
	c.App.Get("/api/testimonials/:id", c.TestimonialController.GetByID)
	c.App.Post("/api/testimonials", c.TestimonialController.Create)
	c.App.Put("/api/testimonials/:id", c.TestimonialController.Update)
	c.App.Delete("/api/testimonials/:id", c.TestimonialController.Delete)

	c.App.Get("/api/hero-slides", c.HeroSlideController.GetAll)
	c.App.Get("/api/hero-slides/:id", c.HeroSlideController.GetByID)
	c.App.Post("/api/hero-slides", c.HeroSlideController.Create)
	c.App.Put("/api/hero-slides/:id", c.HeroSlideController.Update)
	c.App.Delete("/api/hero-slides/:id", c.HeroSlideController.Delete)

}
