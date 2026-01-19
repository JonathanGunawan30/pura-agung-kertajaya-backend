package usecase

import (
	"bytes"
	"context"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	_ "golang.org/x/image/webp"
)

type StorageUsecase interface {
	UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (map[string]string, error)
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

func (u *storageUsecase) UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (map[string]string, error) {
	uploadedKeys := make(map[string]string)
	var mu sync.Mutex

	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Unix()

	onProcessed := func(presetName string, data []byte) error {
		key := fmt.Sprintf("uploads/%s_%d_%s.webp", nameWithoutExt, timestamp, presetName)
		pSize := int64(len(data))

		_, err := u.storageRepo.Upload(ctx, key, bytes.NewReader(data), "image/webp", pSize)
		if err != nil {
			u.log.WithError(err).Errorf("failed to upload variant: %s", presetName)
			return err
		}

		mu.Lock()
		uploadedKeys[presetName] = key
		mu.Unlock()

		return nil
	}

	err := util.ProcessAndHandleImage(file, util.AllPresets, onProcessed)

	if err != nil {
		u.log.WithError(err).Error("failed to process image")

		wrappedErr := fmt.Errorf("invalid image format: %w", err)

		cleanupCtx := context.Background()
		for _, key := range uploadedKeys {
			errDelete := u.storageRepo.Delete(cleanupCtx, key)
			if errDelete != nil {
				u.log.WithError(errDelete).Errorf("failed to rollback file: %s", key)
			}
		}

		return nil, wrappedErr
	}

	u.log.WithField("count", len(uploadedKeys)).Info("all variants uploaded successfully")
	return uploadedKeys, nil
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
