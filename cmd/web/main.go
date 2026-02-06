package main

import (
	"fmt"
	"pura-agung-kertajaya-backend/db/seeder"
	"pura-agung-kertajaya-backend/internal/config"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/entity"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)

	sentryDSN := viperConfig.GetString("sentry.dsn")
	appEnv := viperConfig.GetString("app.env")

	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			EnableTracing:    true,
			TracesSampleRate: 0.2,
			Environment:      appEnv,
		})
		if err != nil {
			logger.Fatalf("sentry.Init: %v", err)
		}
		defer sentry.Flush(2 * time.Second)
		logger.Info("Sentry initialized successfully")
	} else {
		logger.Warn("Sentry DSN not found, skipping Sentry initialization")
	}

	db := config.NewDatabase(viperConfig, logger)
	validate, trans := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	if sentryDSN != "" {
		app.Use(fibersentry.New(fibersentry.Config{
			Repanic:         true,
			WaitForDelivery: true,
		}))
	}

	app.Use(middleware.ErrorHandlerMiddleware(logger, trans))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     viperConfig.GetString("cors.allow_origins"),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      logger,
		Validate: validate,
		Config:   viperConfig,
	})

	err := db.AutoMigrate(
		&entity.User{},
		&entity.Testimonial{},
		&entity.HeroSlide{},
		&entity.Gallery{},
		&entity.Facility{},
		&entity.ContactInfo{},
		&entity.Activity{},
		&entity.SiteIdentity{},
		&entity.AboutSection{},
		&entity.AboutValue{},
		&entity.OrganizationMember{},
		&entity.OrganizationDetail{},
		&entity.Remark{},
		&entity.Category{},
		&entity.Article{},
	)
	if err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}
	logger.Info("Database migration completed")

	if viperConfig.GetBool("db.seed_on_start") {
		seeder.SeedUsers(db)
		logger.Info("Seeding completed")
	}

	webPort := viperConfig.GetInt("web.port")
	logger.Infof("Starting server on port %d", webPort)
	err = app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
