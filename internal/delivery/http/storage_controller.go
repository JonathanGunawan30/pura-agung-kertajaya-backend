package http

import (
	"strings"

	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type StorageController struct {
	UseCase usecase.StorageUsecase
	Log     *logrus.Logger
}

func NewStorageController(usecase usecase.StorageUsecase, log *logrus.Logger) *StorageController {
	return &StorageController{
		UseCase: usecase,
		Log:     log,
	}
}

func (c *StorageController) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("failed to get file from form")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "No file uploaded",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		c.Log.WithField("content_type", contentType).Error("invalid file type")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "Only image files are allowed (JPEG, PNG, WEBP)",
		})
	}

	if file.Size > 10*1024*1024 {
		c.Log.WithField("size", file.Size).Error("file too large")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "File size must not exceed 10MB",
		})
	}

	src, err := file.Open()
	if err != nil {
		c.Log.WithError(err).Error("failed to open file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: "Cannot open file",
		})
	}
	defer src.Close()

	variants, err := c.UseCase.UploadFile(
		ctx.Context(),
		file.Filename,
		src,
		contentType,
		file.Size,
	)
	if err != nil {
		c.Log.WithError(err).Error("failed to upload file")

		statusCode := fiber.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid image") {
			statusCode = fiber.StatusBadRequest
		}

		return ctx.Status(statusCode).JSON(model.WebResponse[any]{
			Errors: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"message":  "File uploaded and processed successfully",
			"filename": file.Filename,
			"variants": variants,
		},
	})
}

func (c *StorageController) Delete(ctx *fiber.Ctx) error {
	key := ctx.Query("key")
	if key == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "Key parameter is required",
		})
	}

	err := c.UseCase.DeleteFile(ctx.Context(), key)
	if err != nil {
		c.Log.WithError(err).Error("failed to delete file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[string]{
		Data: "File deleted successfully",
	})
}

func (c *StorageController) GetPresignedURL(ctx *fiber.Ctx) error {
	key := ctx.Query("key")
	if key == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "Key parameter is required",
		})
	}

	expiration := ctx.QueryInt("expiration", 3600)

	url, err := c.UseCase.GetPresignedURL(ctx.Context(), key, expiration)
	if err != nil {
		c.Log.WithError(err).Error("failed to generate presigned URL")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"url": url,
		},
	})
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
	}

	contentType = strings.ToLower(contentType)
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}
