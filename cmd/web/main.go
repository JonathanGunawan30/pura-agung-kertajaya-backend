package main

import (
	"fmt"
	"pura-agung-kertajaya-backend/db/seeder"
	"pura-agung-kertajaya-backend/internal/config"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/entity"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	app.Use(middleware.ErrorHandlerMiddleware(log))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     viperConfig.GetString("cors.allow_origins"),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
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
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Info("Database migration completed")

	seeder.SeedUsers(db)

	webPort := viperConfig.GetInt("web.port")
	err = app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
