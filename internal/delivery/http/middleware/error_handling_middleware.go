package middleware

import (
	"errors"
	"pura-agung-kertajaya-backend/internal/model"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ErrorHandlerMiddleware(log *logrus.Logger, trans ut.Translator) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err == nil {
			return nil
		}

		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			translations := valErrs.Translate(trans)
			var errorMessages []string
			for _, msg := range translations {
				errorMessages = append(errorMessages, msg)
			}
			log.Warnf("Validation error: %v", errorMessages)
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Errors: strings.Join(errorMessages, " | "),
			})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[any]{Errors: fiberErr.Message})
		}

		var responseErr *model.ResponseError
		if errors.As(err, &responseErr) {
			if responseErr.Code >= 500 {
				log.Errorf("System error: %v", responseErr.Message)
			} else {
				log.Warnf("Business error: %v", responseErr.Message)
			}
			return ctx.Status(responseErr.Code).JSON(model.WebResponse[any]{
				Errors: responseErr.Message,
			})
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{
				Errors: "Data not found",
			})
		}

		log.Errorf("INTERNAL SERVER ERROR: %+v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: "Internal Server Error",
		})
	}
}
