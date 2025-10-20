package repository

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type StorageRepository interface {
	Upload(ctx context.Context, key string, file io.Reader, contentType string, fileSize int64) (string, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, expiration int) (string, error)
}

type storageRepository struct {
	client    *s3.Client
	bucket    string
	publicURL string
	log       *logrus.Logger
}

func NewStorageRepository(client *s3.Client, cfg *viper.Viper, log *logrus.Logger) StorageRepository {
	return &storageRepository{
		client:    client,
		bucket:    cfg.GetString("r2.bucket_name"),
		publicURL: cfg.GetString("r2.public_url"),
		log:       log,
	}
}

func (r *storageRepository) Upload(ctx context.Context, key string, file io.Reader, contentType string, fileSize int64) (string, error) {
	r.log.WithFields(logrus.Fields{
		"key":          key,
		"content_type": contentType,
		"file_size":    fileSize,
	}).Info("uploading file to R2")

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
	})
	if err != nil {
		r.log.WithError(err).Error("failed to upload file to R2")
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("%s/%s", r.publicURL, key)
	r.log.WithField("url", url).Info("file uploaded successfully")

	return url, nil
}

func (r *storageRepository) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	r.log.WithField("key", key).Info("downloading file from R2")

	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		r.log.WithError(err).Error("failed to download file from R2")
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return result.Body, nil
}

func (r *storageRepository) Delete(ctx context.Context, key string) error {
	r.log.WithField("key", key).Info("deleting file from R2")

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		r.log.WithError(err).Error("failed to delete file from R2")
		return fmt.Errorf("failed to delete file: %w", err)
	}

	r.log.Info("file deleted successfully")
	return nil
}

func (r *storageRepository) GetPresignedURL(ctx context.Context, key string, expiration int) (string, error) {
	r.log.WithFields(logrus.Fields{
		"key":        key,
		"expiration": expiration,
	}).Info("generating presigned URL")

	presignClient := s3.NewPresignClient(r.client)

	presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiration) * time.Second
	})

	if err != nil {
		r.log.WithError(err).Error("failed to generate presigned URL")
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}
