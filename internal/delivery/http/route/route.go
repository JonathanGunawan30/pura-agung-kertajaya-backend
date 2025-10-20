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
	GalleryController     *http.GalleryController
	FacilityController    *http.FacilityController
	ContactInfoController *http.ContactInfoController
	ActivityController    *http.ActivityController
	AuthMiddleware        fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users/_login", c.UserController.Login)
	c.App.Get("/api/public/testimonials", c.TestimonialController.GetAllPublic)
	c.App.Get("/api/public/hero-slides", c.HeroSlideController.GetAllPublic)
	c.App.Get("/api/public/galleries", c.GalleryController.GetAllPublic)
	c.App.Get("/api/public/facilities", c.FacilityController.GetAllPublic)
	c.App.Get("/api/public/contact-info", c.ContactInfoController.GetAll)
	c.App.Get("/api/public/activities", c.ActivityController.GetAllPublic)
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

	c.App.Get("/api/galleries", c.GalleryController.GetAll)
	c.App.Get("/api/galleries/:id", c.GalleryController.GetByID)
	c.App.Post("/api/galleries", c.GalleryController.Create)
	c.App.Put("/api/galleries/:id", c.GalleryController.Update)
	c.App.Delete("/api/galleries/:id", c.GalleryController.Delete)

	c.App.Get("/api/facilities", c.FacilityController.GetAll)
	c.App.Get("/api/facilities/:id", c.FacilityController.GetByID)
	c.App.Post("/api/facilities", c.FacilityController.Create)
	c.App.Put("/api/facilities/:id", c.FacilityController.Update)
	c.App.Delete("/api/facilities/:id", c.FacilityController.Delete)

	c.App.Get("/api/contact-info", c.ContactInfoController.GetAll)
	c.App.Get("/api/contact-info/:id", c.ContactInfoController.GetByID)
	c.App.Post("/api/contact-info", c.ContactInfoController.Create)
	c.App.Put("/api/contact-info/:id", c.ContactInfoController.Update)
	c.App.Delete("/api/contact-info/:id", c.ContactInfoController.Delete)

	c.App.Get("/api/activities", c.ActivityController.GetAll)
	c.App.Get("/api/activities/:id", c.ActivityController.GetByID)
	c.App.Post("/api/activities", c.ActivityController.Create)
	c.App.Put("/api/activities/:id", c.ActivityController.Update)
	c.App.Delete("/api/activities/:id", c.ActivityController.Delete)

}
