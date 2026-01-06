package config

import (
	"pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/delivery/http/route"
	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/usecase"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(cfg *BootstrapConfig) {
	// Setup Redis client
	redisHost := cfg.Config.GetString("redis.host")
	redisPort := cfg.Config.GetInt("redis.port")
	redisPass := cfg.Config.GetString("redis.password")
	redisDB := cfg.Config.GetInt("redis.db")
	rateLimiterDB := cfg.Config.GetInt("redis.rate_limiter_db")
	redisTLS := cfg.Config.GetBool("redis.tls")
	redisClient := NewRedisClient(redisHost, redisPort, redisPass, redisDB, redisTLS)

	// Setup TokenUtil (JWT + Redis)
	secretKey := cfg.Config.GetString("jwt.secret")
	tokenUtil := util.NewTokenUtil(secretKey, redisClient.RDB)

	// Setup RecaptchaUtil
	recaptchaUtil := util.NewRecaptchaUtil(cfg.Config)

	r2Client, err := util.NewR2Client(cfg.Config)
	if err != nil {
		cfg.Log.WithError(err).Fatal("failed to initialize R2 client")
	}

	// Setup repositories
	userRepository := repository.NewUserRepository(cfg.Log)
	storageRepository := repository.NewStorageRepository(r2Client, cfg.Config, cfg.Log)

	// Setup usecases
	userUseCase := usecase.NewUserUseCase(cfg.DB, cfg.Log, cfg.Validate, userRepository, tokenUtil, recaptchaUtil, cfg.Config)
	storageUseCase := usecase.NewStorageUsecase(storageRepository, cfg.Log, cfg.Validate)
	testimonialUseCase := usecase.NewTestimonialUsecase(cfg.DB, cfg.Log, cfg.Validate)
	heroSlideUseCase := usecase.NewHeroSlideUsecase(cfg.DB, cfg.Log, cfg.Validate)
	galleryUseCase := usecase.NewGalleryUsecase(cfg.DB, cfg.Log, cfg.Validate)
	facilityUseCase := usecase.NewFacilityUsecase(cfg.DB, cfg.Log, cfg.Validate)
	contactInfoUseCase := usecase.NewContactInfoUsecase(cfg.DB, cfg.Log, cfg.Validate)
	activityUseCase := usecase.NewActivityUsecase(cfg.DB, cfg.Log, cfg.Validate)
	siteIdentityUseCase := usecase.NewSiteIdentityUsecase(cfg.DB, cfg.Log, cfg.Validate)
	aboutUseCase := usecase.NewAboutUsecase(cfg.DB, cfg.Log, cfg.Validate)
	organizationUsecase := usecase.NewOrganizationRequest(cfg.DB, cfg.Log, cfg.Validate)
	remarkUseCase := usecase.NewRemarkUsecase(cfg.DB, cfg.Log, cfg.Validate)
	organizationDetailUsecase := usecase.NewOrganizationDetailUsecase(cfg.DB, cfg.Log, cfg.Validate)
	categoryUsecase := usecase.NewCategoryUsecase(cfg.DB, cfg.Log, cfg.Validate)
	articleUsecase := usecase.NewArticleUsecase(cfg.DB, cfg.Log, cfg.Validate)

	// Setup controllers
	userController := http.NewUserController(userUseCase, cfg.Log)
	storageController := http.NewStorageController(storageUseCase, cfg.Log)
	testimonialController := http.NewTestimonialController(testimonialUseCase, cfg.Log)
	heroSlideController := http.NewHeroSlideController(heroSlideUseCase, cfg.Log)
	galleryController := http.NewGalleryController(galleryUseCase, cfg.Log)
	facilityController := http.NewFacilityController(facilityUseCase, cfg.Log)
	contactInfoController := http.NewContactInfoController(contactInfoUseCase, cfg.Log)
	activityController := http.NewActivityController(activityUseCase, cfg.Log)
	siteIdentityController := http.NewSiteIdentityController(siteIdentityUseCase, cfg.Log)
	aboutController := http.NewAboutController(aboutUseCase, cfg.Log)
	organizationController := http.NewOrganizationController(organizationUsecase, cfg.Log)
	remarkcontroller := http.NewRemarkController(remarkUseCase, cfg.Log)
	organizationDetailController := http.NewOrganizationDetailController(organizationDetailUsecase, cfg.Log)
	categoryController := http.NewCategoryController(categoryUsecase, cfg.Log)
	articleController := http.NewArticleController(articleUsecase, cfg.Log)

	// Setup redis storage
	storage := NewFiberRedisStorage(redisHost, redisPort, redisPass, rateLimiterDB, redisTLS)

	cfg.App.Hooks().OnShutdown(func() error {
		cfg.Log.Info("Closing Redis connections...")
		if err := storage.Close(); err != nil {
			cfg.Log.WithError(err).Error("Failed to close Redis Storage")
		}
		if err := redisClient.RDB.Close(); err != nil {
			cfg.Log.WithError(err).Error("Failed to close Redis Client")
		}
		return nil
	})

	// Setup middleware
	authMiddleware := middleware.AuthMiddleware(tokenUtil)

	// Rate Limiter
	publicRateLimiter := middleware.PublicRateLimiter(storage)
	authRateLimiter := middleware.AuthRateLimiter(storage)
	cmsReadRateLimiter := middleware.CMSReadRateLimiter(storage)
	cmsWriteRateLimiter := middleware.CMSWriteRateLimiter(storage)
	storageRateLimiter := middleware.StorageRateLimiter(storage)
	deleteRateLimiter := middleware.DeleteRateLimiter(storage)

	// Setup routes
	routeConfig := route.RouteConfig{
		App:                          cfg.App,
		UserController:               userController,
		StorageController:            storageController,
		TestimonialController:        testimonialController,
		HeroSlideController:          heroSlideController,
		GalleryController:            galleryController,
		FacilityController:           facilityController,
		ContactInfoController:        contactInfoController,
		ActivityController:           activityController,
		SiteIdentityController:       siteIdentityController,
		AboutController:              aboutController,
		OrganizationController:       organizationController,
		RemarkController:             remarkcontroller,
		OrganizationDetailController: organizationDetailController,
		CategoryController:           categoryController,
		ArticleController:            articleController,

		AuthMiddleware: authMiddleware,

		PublicRateLimiter:   publicRateLimiter,
		AuthRateLimiter:     authRateLimiter,
		CMSReadRateLimiter:  cmsReadRateLimiter,
		CMSWriteRateLimiter: cmsWriteRateLimiter,
		StorageRateLimiter:  storageRateLimiter,
		DeleteRateLimiter:   deleteRateLimiter,
	}
	routeConfig.Setup()
}
