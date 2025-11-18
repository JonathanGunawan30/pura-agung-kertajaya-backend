package main

import (
	"fmt"
	"pura-agung-kertajaya-backend/db/seeder"
	"pura-agung-kertajaya-backend/internal/config"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting app...")

	viperConfig := config.NewViper()
	logrus.Info("Viper loaded")

	log := config.NewLogger(viperConfig)
	logrus.Info("Logger loaded")

	db := config.NewDatabase(viperConfig, log)
	logrus.Info("Database loaded")

	validate := config.NewValidator(viperConfig)
	logrus.Info("Validator loaded")

	app := config.NewFiber(viperConfig)
	logrus.Info("Fiber loaded")

	app.Use(middleware.ErrorHandlerMiddleware(log))
	logrus.Info("Middleware loaded")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     viperConfig.GetString("cors.allow_origins"),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	logrus.Info("Bootstrap done")

	seeder.SeedUsers(db)
	logrus.Info("Seeder done")

	webPort := viperConfig.GetInt("web.port")
	logrus.Infof("Web port: %d", webPort)

	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
