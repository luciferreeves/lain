package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lain/config"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitMinIO() error {
	client, err := minio.New(config.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIO.AccessKey, config.MinIO.SecretKey, ""),
		Secure: config.MinIO.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.MinIO.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, config.MinIO.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	minioClient = client
	return nil
}

func UploadAttachment(userEmail string, emailID uint, filename string, data []byte, contentType string) (string, error) {
	if minioClient == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	path := fmt.Sprintf("attachments/%s/%d/%s", userEmail, emailID, filename)

	ctx := context.Background()

	_, err := minioClient.PutObject(
		ctx,
		config.MinIO.BucketName,
		path,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to upload attachment: %w", err)
	}

	return path, nil
}

func DownloadAttachment(path string) ([]byte, error) {
	if minioClient == nil {
		return nil, fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()

	object, err := minioClient.GetObject(ctx, config.MinIO.BucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment object: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read attachment data: %w", err)
	}

	return data, nil
}

func DeleteAttachment(path string) error {
	if minioClient == nil {
		return fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()

	err := minioClient.RemoveObject(ctx, config.MinIO.BucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	return nil
}

func DeleteAttachmentsByEmail(userEmail string, emailID uint) error {
	if minioClient == nil {
		return fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()
	prefix := fmt.Sprintf("attachments/%s/%d/", userEmail, emailID)

	objectCh := minioClient.ListObjects(ctx, config.MinIO.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return fmt.Errorf("failed to list attachments: %w", object.Err)
		}

		err := minioClient.RemoveObject(ctx, config.MinIO.BucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete attachment %s: %w", object.Key, err)
		}
	}

	return nil
}

func GetAttachmentURL(path string, expiryDuration time.Duration) (string, error) {
	if minioClient == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()

	url, err := minioClient.PresignedGetObject(ctx, config.MinIO.BucketName, path, expiryDuration, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return url.String(), nil
}

func GetAttachmentFilename(path string) string {
	return filepath.Base(path)
}

func AttachmentExists(path string) (bool, error) {
	if minioClient == nil {
		return false, fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()

	_, err := minioClient.StatObject(ctx, config.MinIO.BucketName, path, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check attachment existence: %w", err)
	}

	return true, nil
}
