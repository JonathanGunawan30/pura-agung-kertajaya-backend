package usecase

import (
	"bytes"
	"context"
	"fmt"
	_ "image/gif"  // Penting untuk decode image
	_ "image/jpeg" // Penting untuk decode image
	_ "image/png"  // Penting untuk decode image
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"

	_ "golang.org/x/image/webp" // Penting untuk decode/encode webp
)

type StorageUsecase interface {
	UploadFile(ctx context.Context, filename string, file io.Reader, contentType string, fileSize int64) (map[string]string, error)
	DownloadFile(ctx context.Context, key string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, expiration int) (string, error)
}

type storageUsecase struct {
	storageRepo repository.StorageRepository
}

func NewStorageUsecase(storageRepo repository.StorageRepository) StorageUsecase {
	return &storageUsecase{
		storageRepo: storageRepo,
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
			return fmt.Errorf("failed to upload variant %s: %w", presetName, err)
		}

		mu.Lock()
		uploadedKeys[presetName] = key
		mu.Unlock()

		return nil
	}

	err := util.ProcessAndHandleImage(file, util.AllPresets, onProcessed)
	if err != nil {
		wrappedErr := fmt.Errorf("image processing failed: %w", err)

		cleanupCtx := context.Background()
		for _, key := range uploadedKeys {
			_ = u.storageRepo.Delete(cleanupCtx, key)
		}

		return nil, wrappedErr
	}

	return uploadedKeys, nil
}

func (u *storageUsecase) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	return u.storageRepo.Download(ctx, key)
}

func (u *storageUsecase) DeleteFile(ctx context.Context, key string) error {
	return u.storageRepo.Delete(ctx, key)
}

func (u *storageUsecase) GetPresignedURL(ctx context.Context, key string, expiration int) (string, error) {
	return u.storageRepo.GetPresignedURL(ctx, key, expiration)
}
