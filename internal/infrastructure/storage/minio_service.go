package storage

import (
	"bytes"
	"context"
	"fmt"
	"service-test-runner/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

func NewMinioService(cfg *config.MinIOConfig) (*MinioService, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	// Check if bucket exists, create if it doesn't
	exists, err := client.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	return &MinioService{
		client:     client,
		bucketName: cfg.BucketName,
		endpoint:   cfg.Endpoint,
	}, nil
}

func (s *MinioService) UploadPDF(objectName string, fileBytes []byte) error {
	reader := bytes.NewReader(fileBytes)
	_, err := s.client.PutObject(
		context.Background(),
		s.bucketName,
		objectName,
		reader,
		int64(len(fileBytes)),
		minio.PutObjectOptions{ContentType: "application/pdf"},
	)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

func (s *MinioService) GetFileURL(objectName string) string {
	// Return direct URL to the object in MinIO
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucketName, objectName)
}
