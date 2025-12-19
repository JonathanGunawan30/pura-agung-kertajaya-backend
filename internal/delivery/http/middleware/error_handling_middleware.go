package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"pura-agung-kertajaya-backend/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ErrorHandlerMiddleware(log *logrus.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err == nil {
			return nil
		}

		log.WithError(err).Error("Error occurred during request processing")

		code := fiber.StatusInternalServerError
		message := "An internal server error occurred. Please try again later."

		var fiberErr *fiber.Error
		var validationErrs validator.ValidationErrors
		var syntaxErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError
		var numErr *strconv.NumError

		switch {
		case errors.As(err, &fiberErr):
			code = fiberErr.Code
			message = fiberErr.Message
			if code == fiber.StatusUnauthorized {
				message = "Unauthorized access."
			}

		case errors.Is(err, fiber.ErrUnauthorized):
			code = fiber.StatusUnauthorized
			message = "Unauthorized access."

		case errors.As(err, &validationErrs):
			code = fiber.StatusBadRequest
			message = formatValidationErrors(validationErrs)

		case errors.As(err, &syntaxErr):
			code = fiber.StatusBadRequest
			message = fmt.Sprintf("Invalid JSON format at position %d", syntaxErr.Offset)

		case errors.As(err, &unmarshalErr):
			code = fiber.StatusBadRequest
			message = fmt.Sprintf("Invalid data type for field '%s': expected %s but got %s",
				unmarshalErr.Field, unmarshalErr.Type.String(), unmarshalErr.Value)

		case errors.As(err, &numErr):
			code = fiber.StatusBadRequest
			message = fmt.Sprintf("Invalid number format: '%s' cannot be converted to a number", numErr.Num)

		case errors.Is(err, gorm.ErrRecordNotFound):
			code = fiber.StatusNotFound
			message = "The requested resource was not found."

		case errors.Is(err, gorm.ErrInvalidTransaction):
			code = fiber.StatusInternalServerError
			message = "Database transaction error occurred."

		case errors.Is(err, gorm.ErrInvalidField):
			code = fiber.StatusBadRequest
			message = "Invalid field in request."

		case errors.Is(err, gorm.ErrInvalidData):
			code = fiber.StatusBadRequest
			message = "Invalid data provided."

		case errors.Is(err, gorm.ErrDuplicatedKey):
			code = fiber.StatusConflict
			message = "A record with this data already exists."

		case errors.Is(err, gorm.ErrForeignKeyViolated):
			code = fiber.StatusBadRequest
			message = "Cannot perform this operation due to related records."

		case errors.Is(err, gorm.ErrCheckConstraintViolated):
			code = fiber.StatusBadRequest
			message = "Data validation constraint failed."

		default:
			code, message = detectErrorFromMessage(err.Error())
		}

		if sendErr := ctx.Status(code).JSON(model.WebResponse[any]{
			Errors: message,
			Data:   nil,
		}); sendErr != nil {
			log.WithError(sendErr).Error("FATAL: Failed to send error JSON response")
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}
}

func detectErrorFromMessage(errMsg string) (int, string) {
	lowerErr := strings.ToLower(errMsg)

	if strings.Contains(lowerErr, "key:") && strings.Contains(lowerErr, "error:field validation") {
		return fiber.StatusBadRequest, parseRawValidationError(errMsg)
	}

	if strings.Contains(lowerErr, "unique constraint") ||
		strings.Contains(lowerErr, "duplicate key") ||
		strings.Contains(lowerErr, "violates unique") {
		return fiber.StatusConflict, extractDuplicateFieldError(errMsg)
	}

	if strings.Contains(lowerErr, "foreign key constraint") ||
		strings.Contains(lowerErr, "violates foreign key") {
		return fiber.StatusBadRequest, "Cannot perform this operation due to related records."
	}

	if strings.Contains(lowerErr, "check constraint") {
		return fiber.StatusBadRequest, "Data validation constraint failed."
	}

	if strings.Contains(lowerErr, "not null constraint") ||
		strings.Contains(lowerErr, "cannot be null") {
		return fiber.StatusBadRequest, extractNotNullFieldError(errMsg)
	}

	if strings.Contains(lowerErr, "validation failed") ||
		strings.Contains(lowerErr, "invalid request") ||
		strings.Contains(lowerErr, "bad request") ||
		strings.Contains(lowerErr, "invalid input") ||
		strings.Contains(lowerErr, "field validation") ||
		strings.Contains(lowerErr, "failed on the") {
		return fiber.StatusBadRequest, errMsg
	}

	if strings.Contains(lowerErr, "too long") ||
		strings.Contains(lowerErr, "exceeds maximum length") ||
		strings.Contains(lowerErr, "maximum length exceeded") {
		return fiber.StatusBadRequest, errMsg
	}

	if strings.Contains(lowerErr, "too short") ||
		strings.Contains(lowerErr, "minimum length") {
		return fiber.StatusBadRequest, errMsg
	}

	if strings.Contains(lowerErr, "unauthorized") ||
		strings.Contains(lowerErr, "token expired") ||
		strings.Contains(lowerErr, "invalid token") ||
		strings.Contains(lowerErr, "authentication failed") {
		return fiber.StatusUnauthorized, "Unauthorized access."
	}

	if strings.Contains(lowerErr, "forbidden") ||
		strings.Contains(lowerErr, "permission denied") ||
		strings.Contains(lowerErr, "access denied") {
		return fiber.StatusForbidden, "You do not have permission to access this resource."
	}

	if strings.Contains(lowerErr, "not found") ||
		strings.Contains(lowerErr, "missing") ||
		strings.Contains(lowerErr, "no such record") ||
		strings.Contains(lowerErr, "does not exist") {
		return fiber.StatusNotFound, "The requested resource was not found."
	}

	if strings.Contains(lowerErr, "conflict") ||
		strings.Contains(lowerErr, "already exists") {
		return fiber.StatusConflict, "The resource already exists or conflicts with existing data."
	}

	if strings.Contains(lowerErr, "timeout") ||
		strings.Contains(lowerErr, "deadline exceeded") ||
		strings.Contains(lowerErr, "context deadline") {
		return fiber.StatusRequestTimeout, "The request took too long to process."
	}

	if strings.Contains(lowerErr, "too many requests") ||
		strings.Contains(lowerErr, "rate limit") {
		return fiber.StatusTooManyRequests, "Too many requests. Please slow down."
	}

	if strings.Contains(lowerErr, "service unavailable") ||
		strings.Contains(lowerErr, "temporarily unavailable") {
		return fiber.StatusServiceUnavailable, "Service temporarily unavailable. Please try again later."
	}

	if strings.Contains(lowerErr, "connection refused") ||
		strings.Contains(lowerErr, "connection reset") ||
		strings.Contains(lowerErr, "broken pipe") {
		return fiber.StatusServiceUnavailable, "Service connection error. Please try again later."
	}

	if strings.Contains(lowerErr, "db") ||
		strings.Contains(lowerErr, "sql") ||
		strings.Contains(lowerErr, "query") ||
		strings.Contains(lowerErr, "database") ||
		strings.Contains(lowerErr, "driver") ||
		strings.Contains(lowerErr, "update failed") ||
		strings.Contains(lowerErr, "insert failed") ||
		strings.Contains(lowerErr, "save failed") ||
		strings.Contains(lowerErr, "delete failed") {
		return fiber.StatusInternalServerError, "An internal server error occurred. Please try again later."
	}

	return fiber.StatusInternalServerError, errMsg
}

func formatValidationErrors(errs validator.ValidationErrors) string {
	var errorMessages []string
	for _, err := range errs {
		var msg string
		field := err.Field()
		param := err.Param()

		switch err.Tag() {
		case "required":
			msg = fmt.Sprintf("'%s' is required", field)
		case "email":
			msg = fmt.Sprintf("'%s' must be a valid email address", field)
		case "min":
			kind := err.Kind().String()
			if kind == "string" {
				msg = fmt.Sprintf("'%s' must be at least %s characters long", field, param)
			} else if kind == "slice" || kind == "array" || kind == "map" {
				msg = fmt.Sprintf("'%s' must have at least %s items", field, param)
			} else {
				msg = fmt.Sprintf("'%s' must be at least %s", field, param)
			}
		case "max":
			kind := err.Kind().String()
			if kind == "string" {
				msg = fmt.Sprintf("'%s' must not exceed %s characters", field, param)
			} else if kind == "slice" || kind == "array" || kind == "map" {
				msg = fmt.Sprintf("'%s' must have at most %s items", field, param)
			} else {
				msg = fmt.Sprintf("'%s' must not exceed %s", field, param)
			}
		case "url":
			msg = fmt.Sprintf("'%s' must be a valid URL", field)
		case "uri":
			msg = fmt.Sprintf("'%s' must be a valid URI", field)
		case "boolean":
			msg = fmt.Sprintf("'%s' must be true or false", field)
		case "len":
			kind := err.Kind().String()
			if kind == "string" {
				msg = fmt.Sprintf("'%s' must be exactly %s characters long", field, param)
			} else {
				msg = fmt.Sprintf("'%s' must contain exactly %s items", field, param)
			}
		case "eqfield":
			msg = fmt.Sprintf("'%s' must match the '%s' field", field, param)
		case "nefield":
			msg = fmt.Sprintf("'%s' must not match the '%s' field", field, param)
		case "gt":
			msg = fmt.Sprintf("'%s' must be greater than %s", field, param)
		case "gte":
			msg = fmt.Sprintf("'%s' must be greater than or equal to %s", field, param)
		case "lt":
			msg = fmt.Sprintf("'%s' must be less than %s", field, param)
		case "lte":
			msg = fmt.Sprintf("'%s' must be less than or equal to %s", field, param)
		case "numeric":
			msg = fmt.Sprintf("'%s' must be a valid number", field)
		case "number":
			msg = fmt.Sprintf("'%s' must be a number", field)
		case "alpha":
			msg = fmt.Sprintf("'%s' must contain only letters", field)
		case "alphanum":
			msg = fmt.Sprintf("'%s' must contain only letters and numbers", field)
		case "alphanumspace":
			msg = fmt.Sprintf("'%s' must contain only letters, numbers, and spaces", field)
		case "uuid", "uuid4", "uuid5":
			msg = fmt.Sprintf("'%s' must be a valid UUID", field)
		case "oneof":
			options := strings.Replace(param, " ", ", ", -1)
			msg = fmt.Sprintf("'%s' must be one of [%s]", field, options)
		case "datetime":
			msg = fmt.Sprintf("'%s' must be a valid date/time in the format %s", field, param)
		case "e164":
			msg = fmt.Sprintf("'%s' must be a valid phone number", field)
		case "json":
			msg = fmt.Sprintf("'%s' must be valid JSON", field)
		case "latitude":
			msg = fmt.Sprintf("'%s' must be a valid latitude", field)
		case "longitude":
			msg = fmt.Sprintf("'%s' must be a valid longitude", field)
		case "hexcolor":
			msg = fmt.Sprintf("'%s' must be a valid hex color", field)
		case "rgb", "rgba":
			msg = fmt.Sprintf("'%s' must be a valid RGB color", field)
		case "isbn", "isbn10", "isbn13":
			msg = fmt.Sprintf("'%s' must be a valid ISBN", field)
		case "ip", "ipv4", "ipv6":
			msg = fmt.Sprintf("'%s' must be a valid IP address", field)
		case "mac":
			msg = fmt.Sprintf("'%s' must be a valid MAC address", field)
		case "base64":
			msg = fmt.Sprintf("'%s' must be valid base64", field)
		case "contains":
			msg = fmt.Sprintf("'%s' must contain '%s'", field, param)
		case "containsany":
			msg = fmt.Sprintf("'%s' must contain at least one of these characters: %s", field, param)
		case "excludes":
			msg = fmt.Sprintf("'%s' must not contain '%s'", field, param)
		case "startswith":
			msg = fmt.Sprintf("'%s' must start with '%s'", field, param)
		case "endswith":
			msg = fmt.Sprintf("'%s' must end with '%s'", field, param)
		default:
			msg = fmt.Sprintf("Field '%s' failed validation rule '%s'", field, err.Tag())
		}
		errorMessages = append(errorMessages, msg)
	}
	return strings.Join(errorMessages, "; ")
}

func parseRawValidationError(errMsg string) string {
	if strings.Contains(errMsg, "failed on the '") {
		parts := strings.Split(errMsg, "failed on the '")
		if len(parts) >= 2 {
			tag := strings.Split(parts[1], "'")[0]

			if strings.Contains(errMsg, "Field validation for '") {
				fieldParts := strings.Split(errMsg, "Field validation for '")
				if len(fieldParts) >= 2 {
					field := strings.Split(fieldParts[1], "'")[0]

					switch tag {
					case "required":
						return fmt.Sprintf("'%s' is required", field)
					case "min":
						return fmt.Sprintf("'%s' does not meet the minimum requirement", field)
					case "max":
						return fmt.Sprintf("'%s' exceeds the maximum allowed length", field)
					case "email":
						return fmt.Sprintf("'%s' must be a valid email address", field)
					case "url":
						return fmt.Sprintf("'%s' must be a valid URL", field)
					default:
						return fmt.Sprintf("Field '%s' failed validation rule '%s'", field, tag)
					}
				}
			}
		}
	}

	return "Validation failed. Please check your input data."
}

func extractDuplicateFieldError(errMsg string) string {
	lowerErr := strings.ToLower(errMsg)

	if strings.Contains(lowerErr, "duplicate key value") {
		if strings.Contains(errMsg, "(") && strings.Contains(errMsg, ")") {
			start := strings.Index(errMsg, "(")
			end := strings.Index(errMsg, ")")
			if start != -1 && end != -1 && end > start {
				field := errMsg[start+1 : end]
				return fmt.Sprintf("A record with this %s already exists", field)
			}
		}
	}

	if strings.Contains(lowerErr, "for key") {
		parts := strings.Split(errMsg, "for key")
		if len(parts) >= 2 {
			key := strings.TrimSpace(parts[1])
			key = strings.Trim(key, "'\"")
			return fmt.Sprintf("A record with this data already exists (key: %s)", key)
		}
	}

	return "A record with this data already exists."
}

func extractNotNullFieldError(errMsg string) string {
	lowerErr := strings.ToLower(errMsg)

	if strings.Contains(lowerErr, "column") && strings.Contains(lowerErr, "cannot be null") {
		parts := strings.Split(lowerErr, "column")
		if len(parts) >= 2 {
			columnPart := strings.TrimSpace(parts[1])
			columnName := strings.Split(columnPart, " ")[0]
			columnName = strings.Trim(columnName, "'\"")
			return fmt.Sprintf("'%s' is required and cannot be null", columnName)
		}
	}

	return "Required field cannot be null."
}
