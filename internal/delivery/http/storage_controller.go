package http

import (
	"fmt"
	"strings"

	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
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

func (c *StorageController) getLogger(ctx *fiber.Ctx) *logrus.Entry {
	user := middleware.GetUser(ctx)

	userID := "guest"
	userRole := "unknown"

	if user != nil {
		userID = fmt.Sprintf("%v", user.ID)
		userRole = user.Role
	}

	return c.Log.WithFields(logrus.Fields{
		"user_id":   userID,
		"user_role": userRole,
		"ip":        ctx.IP(),
		"req_id":    ctx.Get("X-Request-ID"),
	})
}

func (c *StorageController) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.getLogger(ctx).Warnf("failed to get file from form: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "No file uploaded",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		c.getLogger(ctx).WithField("content_type", contentType).Warn("invalid file type upload attempt")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "Only image files are allowed (JPEG, PNG, WEBP)",
		})
	}

	if file.Size > 10*1024*1024 { // 10MB Limit
		c.getLogger(ctx).WithField("size", file.Size).Warn("file too large upload attempt")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "File size must not exceed 10MB",
		})
	}

	src, err := file.Open()
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to open file stream")
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
		if strings.Contains(err.Error(), "invalid image") || strings.Contains(err.Error(), "format") {
			c.getLogger(ctx).WithError(err).Warn("image processing failed")
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Errors: "Invalid image format or corrupted file",
			})
		}

		c.getLogger(ctx).WithError(err).Error("failed to upload/process file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: "Internal server error during file upload",
		})
	}

	c.getLogger(ctx).WithFields(logrus.Fields{
		"filename": file.Filename,
		"variants": len(variants),
	}).Info("File uploaded and processed successfully")

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"message":  "File uploaded and processed successfully",
			"filename": file.Filename,
			"variants": variants,
		},
	})
}

func (c *StorageController) UploadSingle(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.getLogger(ctx).Warnf("failed to get file from form: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "No file uploaded",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		c.getLogger(ctx).WithField("content_type", contentType).Warn("invalid file type upload attempt")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "Only image files are allowed (JPEG, PNG, WEBP)",
		})
	}

	if file.Size > 1*1024*1024 {
		c.getLogger(ctx).WithField("size", file.Size).Warn("file too large upload attempt")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Errors: "File size must not exceed 10MB",
		})
	}

	src, err := file.Open()
	if err != nil {
		c.getLogger(ctx).WithError(err).Error("failed to open file stream")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: "Cannot open file",
		})
	}

	defer src.Close()

	uploadedKey, uploadedURL, err := c.UseCase.UploadSingleFile(
		ctx.Context(),
		file.Filename,
		src,
		contentType,
		file.Size,
	)

	if err != nil {
		if strings.Contains(err.Error(), "invalid image") || strings.Contains(err.Error(), "format") {
			c.getLogger(ctx).WithError(err).Warn("image processing failed")
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Errors: "Invalid image format or corrupted file",
			})
		}

		c.getLogger(ctx).WithError(err).Error("failed to upload/process single file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: "Internal server error during file upload",
		})
	}

	c.getLogger(ctx).WithFields(logrus.Fields{
		"filename": file.Filename,
		"key":      uploadedKey,
	}).Info("Single file uploaded successfully")

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"message":  "File uploaded successfully",
			"filename": file.Filename,
			"key":      uploadedKey,
			"url":      uploadedURL,
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
		c.getLogger(ctx).WithField("key", key).WithError(err).Error("failed to delete file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: err.Error(),
		})
	}

	c.getLogger(ctx).WithField("key", key).Info("File deleted successfully")
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
		c.getLogger(ctx).WithField("key", key).WithError(err).Error("failed to generate presigned URL")
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
