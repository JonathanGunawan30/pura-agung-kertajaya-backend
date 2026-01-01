package usecase

import (
	"bytes"
	"context"
	"fmt"
	"image"
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

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fileBytes := buf.Bytes()

	_, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		u.log.WithError(err).Error("invalid image format")
		return nil, fmt.Errorf("invalid image format: %w", err)
	}

	processedMap, err := util.ProcessImage(bytes.NewReader(fileBytes), util.AllPresets)
	if err != nil {
		u.log.WithError(err).Error("failed to process image variants")
		return nil, fmt.Errorf("failed to process image variants: %w", err)
	}

	var wg sync.WaitGroup
	var errOnce sync.Once
	var uploadErr error

	for presetName, data := range processedMap {
		wg.Add(1)

		pName := presetName
		pData := data

		go func() {
			defer wg.Done()

			key := fmt.Sprintf("uploads/%s_%d_%s.webp", nameWithoutExt, timestamp, pName)
			fileReader := bytes.NewReader(pData)
			pSize := int64(len(pData))

			_, err = u.storageRepo.Upload(ctx, key, fileReader, "image/webp", pSize)
			if err != nil {
				u.log.WithError(err).Errorf("failed to upload variant: %s", pName)
				errOnce.Do(func() {
					uploadErr = err
				})
				return
			}

			mu.Lock()
			uploadedKeys[pName] = key
			mu.Unlock()
		}()
	}

	wg.Wait()

	if uploadErr != nil {
		u.log.WithError(uploadErr).Warn("one or more uploads failed, rolling back...")

		cleanupCtx := context.Background()
		for _, key := range uploadedKeys {
			errDelete := u.storageRepo.Delete(cleanupCtx, key)
			if errDelete != nil {
				u.log.WithError(errDelete).Errorf("CRITICAL: failed to rollback file: %s", key)
			} else {
				u.log.Infof("rollback successful for: %s", key)
			}
		}

		return nil, uploadErr
	}

	u.log.WithField("keys", uploadedKeys).Info("files uploaded successfully")
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
