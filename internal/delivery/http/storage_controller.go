package http

import (
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

	src, err := file.Open()
	if err != nil {
		c.Log.WithError(err).Error("failed to open file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: "Cannot open file",
		})
	}
	defer src.Close()

	url, err := c.UseCase.UploadFile(
		ctx.Context(),
		file.Filename,
		src,
		file.Header.Get("Content-Type"),
		file.Size,
	)
	if err != nil {
		c.Log.WithError(err).Error("failed to upload file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Errors: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"url":      url,
			"filename": file.Filename,
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
		return err
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

	expiration := ctx.QueryInt("expiration", 3600) // default 1 hour

	url, err := c.UseCase.GetPresignedURL(ctx.Context(), key, expiration)
	if err != nil {
		c.Log.WithError(err).Error("failed to generate presigned URL")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"url": url,
		},
	})
}
