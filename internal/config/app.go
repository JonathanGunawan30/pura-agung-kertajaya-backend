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
	"github.com/redis/go-redis/v9"
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Config.GetString("redis.host"),
		Password: cfg.Config.GetString("redis.password"),
		DB:       cfg.Config.GetInt("redis.db"),
	})

	// Setup TokenUtil (JWT + Redis)
	secretKey := cfg.Config.GetString("jwt.secret")
	tokenUtil := util.NewTokenUtil(secretKey, redisClient)

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
	userUseCase := usecase.NewUserUseCase(cfg.DB, cfg.Log, cfg.Validate, userRepository, tokenUtil, recaptchaUtil)
	storageUseCase := usecase.NewStorageUsecase(storageRepository, cfg.Log, cfg.Validate)
	testimonialUseCase := usecase.NewTestimonialUsecase(cfg.DB, cfg.Log, cfg.Validate)
	heroSlideUseCase := usecase.NewHeroSlideUsecase(cfg.DB, cfg.Log, cfg.Validate)
	galleryUseCase := usecase.NewGalleryUsecase(cfg.DB, cfg.Log, cfg.Validate)
	facilityUseCase := usecase.NewFacilityUsecase(cfg.DB, cfg.Log, cfg.Validate)
	contactInfoUseCase := usecase.NewContactInfoUsecase(cfg.DB, cfg.Log, cfg.Validate)
	activityUseCase := usecase.NewActivityUsecase(cfg.DB, cfg.Log, cfg.Validate)

	// Setup controllers
	userController := http.NewUserController(userUseCase, cfg.Log)
	storageController := http.NewStorageController(storageUseCase, cfg.Log)
	testimonialController := http.NewTestimonialController(testimonialUseCase, cfg.Log)
	heroSlideController := http.NewHeroSlideController(heroSlideUseCase, cfg.Log)
	galleryController := http.NewGalleryController(galleryUseCase, cfg.Log)
	facilityController := http.NewFacilityController(facilityUseCase, cfg.Log)
	contactInfoController := http.NewContactInfoController(contactInfoUseCase, cfg.Log)
	activityController := http.NewActivityController(activityUseCase, cfg.Log)

	// Setup middleware
	authMiddleware := middleware.AuthMiddleware(tokenUtil)

	// Setup routes
	routeConfig := route.RouteConfig{
		App:                   cfg.App,
		UserController:        userController,
		StorageController:     storageController,
		TestimonialController: testimonialController,
		HeroSlideController:   heroSlideController,
		GalleryController:     galleryController,
		FacilityController:    facilityController,
		ContactInfoController: contactInfoController,
		ActivityController:    activityController,
		AuthMiddleware:        authMiddleware,
	}
	routeConfig.Setup()
}
