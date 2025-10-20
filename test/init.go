package test

import (
	"pura-agung-kertajaya-backend/internal/config"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App
var db *gorm.DB
var viperConfig *viper.Viper
var log *logrus.Logger
var validate *validator.Validate
var redisClient *redis.Client
var tokenUtil *util.TokenUtil

func init() {
	viperConfig = config.NewViper()
	viperConfig.Set("app.env", "development")
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log)

	// Setup Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     viperConfig.GetString("redis.host"),
		Password: viperConfig.GetString("redis.password"),
		DB:       viperConfig.GetInt("redis.db"),
	})

	secret := viperConfig.GetString("jwt.secret")
	tokenUtil = util.NewTokenUtil(secret, redisClient)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})
}
