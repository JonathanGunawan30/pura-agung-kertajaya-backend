package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"pura-agung-kertajaya-backend/internal/util"
	"strings"
	"time"

	"pura-agung-kertajaya-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type StorageUsecase interface {
	UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (string, error)
	DownloadFile(ctx context.Context, key string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, expiration int) (string, error)
}

type storageUsecase struct {
	storageRepo repository.StorageRepository
	log         *logrus.Logger
	validate    *validator.Validate
}

func NewStorageUsecase(
	storageRepo repository.StorageRepository,
	log *logrus.Logger,
	validate *validator.Validate,
) StorageUsecase {
	return &storageUsecase{
		storageRepo: storageRepo,
		log:         log,
		validate:    validate,
	}
}

func (u *storageUsecase) UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (string, error) {
	finalFile := file
	finalSize := fileSize
	finalContentType := contentType

	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	if strings.HasPrefix(contentType, "image/") {
		presets := []util.ImagePreset{util.PresetDesktop}

		processedMap, err := util.ProcessImage(file, presets)
		if err == nil {
			if data, ok := processedMap["desktop"]; ok {
				finalFile = bytes.NewReader(data)
				finalSize = int64(len(data))
				finalContentType = "image/webp"
				ext = ".webp"
			}
		} else {
			u.log.WithError(err).Warn("Failed to process image, uploading original file instead")
		}
	}
	key := fmt.Sprintf("uploads/%s_%d%s", nameWithoutExt, time.Now().Unix(), ext)

	u.log.WithFields(logrus.Fields{
		"filename":     filename,
		"key":          key,
		"content_type": finalContentType,
		"file_size":    finalSize,
	}).Info("uploading file")

	return u.storageRepo.Upload(ctx, key, finalFile, finalContentType, finalSize)
}
func (u *storageUsecase) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	u.log.WithField("key", key).Info("downloading file")
	return u.storageRepo.Download(ctx, key)
}

func (u *storageUsecase) DeleteFile(ctx context.Context, key string) error {
	u.log.WithField("key", key).Info("deleting file")
	return u.storageRepo.Delete(ctx, key)
}

func (u *storageUsecase) GetPresignedURL(ctx context.Context, key string, expiration int) (string, error) {
	u.log.WithFields(logrus.Fields{
		"key":        key,
		"expiration": expiration,
	}).Info("generating presigned URL")
	return u.storageRepo.GetPresignedURL(ctx, key, expiration)
}
