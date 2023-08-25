package business

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	log "golang.org/x/exp/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	errHandler "github.com/phllpmcphrsn/voice-quips/errors"
	"github.com/phllpmcphrsn/voice-quips/models"
	"github.com/phllpmcphrsn/voice-quips/persistence"
)

// S3UploadService defines the API for uploading objects to S3 storage
type S3UploadService interface {
    UploadObject(ctx context.Context, filename, bucket string) error
}
// MetadataStoreService defines the API for interacting with NO/SQL storage
type MetadataStoreService interface {
	SaveMetadata(context.Context, models.Metadata)
	DeleteMetadata(context.Context, string)
	FindMetadataById(context.Context, string)
	GetAllMetadata(context.Context, []models.Metadata)
}

// AWSUploader is a struct that implements the S3UploadService interface using AWS's S3 SDK for Go
type AWSUploader struct {
    // S3Client is the service client for Amazon S3
    S3Client *s3.Client
	
	// MetadataStore is the SQL client
	MetadataStore persistence.MetadataRepository
}

func NewAWSUploader(client *s3.Client, db persistence.MetadataRepository) *AWSUploader {
	return &AWSUploader{S3Client: client, MetadataStore: db}
}

// UploadObject uploads to an AWS bucket with the given file
func (uploader *AWSUploader) UploadObject(ctx context.Context, filename, bucket string) error {
	var uploadError *errHandler.UploadError

	file, err := os.Open(filename)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	defer file.Close()

	// Want to determine the file extension so that we can set the approriate content-type header
	extension := filepath.Ext(filename)
	contentType := GetContentType(extension)

	response, err := uploader.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key: &filename,
		Body: file,
		ContentType: &contentType,
	})
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	log.Debug("response from object being uploaded", "metadata", response.ResultMetadata)
	return nil
}

func (uploader *AWSUploader) SaveMetadata(ctx context.Context, metadata models.Metadata) (*models.Metadata, error) {
	return uploader.MetadataStore.CreateMetadata(ctx, metadata)
}

// MinioUploader is a struct that implements the S3UploadService interface using MinIO's S3 SDK for Go
type MinioUploader struct {
    // S3Client is the service client for MinIO
    S3Client *minio.Client
	
	// MetadataStore is the SQL client
	MetadataStore *sql.DB
}

func NewMinioUploader(client *minio.Client, db *sql.DB) *MinioUploader {
	return &MinioUploader{S3Client: client, MetadataStore: db}
}

// UploadObject uploads to an MinIO bucket with the given file
func (uploader *MinioUploader) UploadObject(ctx context.Context, filename, bucket string) error {
	var uploadError *errHandler.UploadError

	file, err := os.Open(filename)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}
	
	defer file.Close()

	// Want to determine the file extension so that we can set the approriate content-type header
	extension := filepath.Ext(filename)
	contentType := GetContentType(extension)

	info, err := file.Stat()
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	response, err := uploader.S3Client.PutObject(
		ctx, 
		bucket, 
		filename, 
		file, 
		info.Size(), 
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	log.Debug("response from object being uploaded", "metadata", response)
	return nil
}