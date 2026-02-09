package route

import (
	"pura-agung-kertajaya-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                          *fiber.App
	UserController               *http.UserController
	StorageController            *http.StorageController
	TestimonialController        *http.TestimonialController
	HeroSlideController          *http.HeroSlideController
	GalleryController            *http.GalleryController
	FacilityController           *http.FacilityController
	ContactInfoController        *http.ContactInfoController
	ActivityController           *http.ActivityController
	SiteIdentityController       *http.SiteIdentityController
	AboutController              *http.AboutController
	OrganizationController       *http.OrganizationController
	OrganizationDetailController *http.OrganizationDetailController
	RemarkController             *http.RemarkController
	CategoryController           *http.CategoryController
	ArticleController            *http.ArticleController
	AuthMiddleware               fiber.Handler
	EntityTypeMiddleware         fiber.Handler

	PublicRateLimiter   fiber.Handler
	AuthRateLimiter     fiber.Handler
	CMSReadRateLimiter  fiber.Handler
	CMSWriteRateLimiter fiber.Handler
	StorageRateLimiter  fiber.Handler
	DeleteRateLimiter   fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	public := c.App.Group("/api/public", c.PublicRateLimiter)

	public.Get("/testimonials", c.TestimonialController.GetAllPublic)
	public.Get("/hero-slides", c.HeroSlideController.GetAllPublic)
	public.Get("/galleries", c.GalleryController.GetAllPublic)
	public.Get("/facilities", c.FacilityController.GetAllPublic)
	public.Get("/contact-info", c.ContactInfoController.GetAll)
	public.Get("/activities", c.ActivityController.GetAllPublic)
	public.Get("/site-identity", c.SiteIdentityController.GetPublic)
	public.Get("/about", c.AboutController.GetAllPublic)
	public.Get("/organization-members", c.OrganizationController.GetAllPublic)
	public.Get("/organization-details", c.OrganizationDetailController.GetPublic)
	public.Get("/remarks", c.RemarkController.GetAllPublic)
	public.Get("/categories", c.CategoryController.GetAllPublic)
	public.Get("/articles", c.ArticleController.GetPublic)
	public.Get("/articles/:slug", c.ArticleController.GetBySlug)

	c.App.Post("/api/users/_login", c.AuthRateLimiter, c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	auth := c.App.Group("/api", c.AuthMiddleware, c.EntityTypeMiddleware)

	auth.Post("/users/_logout", c.CMSWriteRateLimiter, c.UserController.Logout)
	auth.Patch("/users/_current", c.CMSWriteRateLimiter, c.UserController.UpdateProfile)
	auth.Get("/users/_current", c.CMSReadRateLimiter, c.UserController.Current)

	storage := auth.Group("/storage", c.StorageRateLimiter)
	storage.Post("/upload", c.StorageController.Upload)
	storage.Post("/upload/single", c.StorageController.UploadSingle)
	storage.Delete("/delete", c.StorageController.Delete)
	storage.Get("/presigned-url", c.StorageController.GetPresignedURL)

	auth.Get("/testimonials", c.CMSReadRateLimiter, c.TestimonialController.GetAll)
	auth.Get("/testimonials/:id", c.CMSReadRateLimiter, c.TestimonialController.GetByID)
	auth.Post("/testimonials", c.CMSWriteRateLimiter, c.TestimonialController.Create)
	auth.Put("/testimonials/:id", c.CMSWriteRateLimiter, c.TestimonialController.Update)
	auth.Delete("/testimonials/:id", c.DeleteRateLimiter, c.TestimonialController.Delete)

	auth.Get("/hero-slides", c.CMSReadRateLimiter, c.HeroSlideController.GetAll)
	auth.Get("/hero-slides/:id", c.CMSReadRateLimiter, c.HeroSlideController.GetByID)
	auth.Post("/hero-slides", c.CMSWriteRateLimiter, c.HeroSlideController.Create)
	auth.Put("/hero-slides/:id", c.CMSWriteRateLimiter, c.HeroSlideController.Update)
	auth.Delete("/hero-slides/:id", c.DeleteRateLimiter, c.HeroSlideController.Delete)

	auth.Get("/galleries", c.CMSReadRateLimiter, c.GalleryController.GetAll)
	auth.Get("/galleries/:id", c.CMSReadRateLimiter, c.GalleryController.GetByID)
	auth.Post("/galleries", c.CMSWriteRateLimiter, c.GalleryController.Create)
	auth.Put("/galleries/:id", c.CMSWriteRateLimiter, c.GalleryController.Update)
	auth.Delete("/galleries/:id", c.DeleteRateLimiter, c.GalleryController.Delete)

	auth.Get("/facilities", c.CMSReadRateLimiter, c.FacilityController.GetAll)
	auth.Get("/facilities/:id", c.CMSReadRateLimiter, c.FacilityController.GetByID)
	auth.Post("/facilities", c.CMSWriteRateLimiter, c.FacilityController.Create)
	auth.Put("/facilities/:id", c.CMSWriteRateLimiter, c.FacilityController.Update)
	auth.Delete("/facilities/:id", c.DeleteRateLimiter, c.FacilityController.Delete)

	auth.Get("/contact-info", c.CMSReadRateLimiter, c.ContactInfoController.GetAll)
	auth.Get("/contact-info/:id", c.CMSReadRateLimiter, c.ContactInfoController.GetByID)
	auth.Post("/contact-info", c.CMSWriteRateLimiter, c.ContactInfoController.Create)
	auth.Put("/contact-info/:id", c.CMSWriteRateLimiter, c.ContactInfoController.Update)
	auth.Delete("/contact-info/:id", c.DeleteRateLimiter, c.ContactInfoController.Delete)

	auth.Get("/activities", c.CMSReadRateLimiter, c.ActivityController.GetAll)
	auth.Get("/activities/:id", c.CMSReadRateLimiter, c.ActivityController.GetByID)
	auth.Post("/activities", c.CMSWriteRateLimiter, c.ActivityController.Create)
	auth.Put("/activities/:id", c.CMSWriteRateLimiter, c.ActivityController.Update)
	auth.Delete("/activities/:id", c.DeleteRateLimiter, c.ActivityController.Delete)

	auth.Get("/site-identity", c.CMSReadRateLimiter, c.SiteIdentityController.GetAll)
	auth.Get("/site-identity/:id", c.CMSReadRateLimiter, c.SiteIdentityController.GetByID)
	auth.Post("/site-identity", c.CMSWriteRateLimiter, c.SiteIdentityController.Create)
	auth.Put("/site-identity/:id", c.CMSWriteRateLimiter, c.SiteIdentityController.Update)
	auth.Delete("/site-identity/:id", c.DeleteRateLimiter, c.SiteIdentityController.Delete)

	auth.Get("/about", c.CMSReadRateLimiter, c.AboutController.GetAll)
	auth.Get("/about/:id", c.CMSReadRateLimiter, c.AboutController.GetByID)
	auth.Post("/about", c.CMSWriteRateLimiter, c.AboutController.Create)
	auth.Put("/about/:id", c.CMSWriteRateLimiter, c.AboutController.Update)
	auth.Delete("/about/:id", c.DeleteRateLimiter, c.AboutController.Delete)

	auth.Get("/organization-members", c.CMSReadRateLimiter, c.OrganizationController.GetAll)
	auth.Get("/organization-members/:id", c.CMSReadRateLimiter, c.OrganizationController.GetByID)
	auth.Post("/organization-members", c.CMSWriteRateLimiter, c.OrganizationController.Create)
	auth.Put("/organization-members/:id", c.CMSWriteRateLimiter, c.OrganizationController.Update)
	auth.Delete("/organization-members/:id", c.DeleteRateLimiter, c.OrganizationController.Delete)

	auth.Get("/remarks", c.CMSReadRateLimiter, c.RemarkController.GetAll)
	auth.Get("/remarks/:id", c.CMSReadRateLimiter, c.RemarkController.GetByID)
	auth.Post("/remarks", c.CMSWriteRateLimiter, c.RemarkController.Create)
	auth.Put("/remarks/:id", c.CMSWriteRateLimiter, c.RemarkController.Update)
	auth.Delete("/remarks/:id", c.DeleteRateLimiter, c.RemarkController.Delete)

	auth.Get("/organization-details", c.CMSReadRateLimiter, c.OrganizationDetailController.GetAdmin)
	auth.Put("/organization-details", c.CMSWriteRateLimiter, c.OrganizationDetailController.Update)

	auth.Get("/categories", c.CMSReadRateLimiter, c.CategoryController.GetAll)
	auth.Get("/categories/:id", c.CMSReadRateLimiter, c.CategoryController.GetByID)
	auth.Post("/categories", c.CMSWriteRateLimiter, c.CategoryController.Create)
	auth.Put("/categories/:id", c.CMSWriteRateLimiter, c.CategoryController.Update)
	auth.Delete("/categories/:id", c.DeleteRateLimiter, c.CategoryController.Delete)

	auth.Get("/articles", c.CMSReadRateLimiter, c.ArticleController.GetAll)
	auth.Get("/articles/:id", c.CMSReadRateLimiter, c.ArticleController.GetByID)
	auth.Post("/articles", c.CMSWriteRateLimiter, c.ArticleController.Create)
	auth.Put("/articles/:id", c.CMSWriteRateLimiter, c.ArticleController.Update)
	auth.Delete("/articles/:id", c.DeleteRateLimiter, c.ArticleController.Delete)
}
