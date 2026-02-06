package test

import (
	"io"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewTestApp() (*fiber.App, *logrus.Logger, ut.Translator) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	app := fiber.New()
	app.Use(middleware.ErrorHandlerMiddleware(logger, trans))
	return app, logger, trans
}
