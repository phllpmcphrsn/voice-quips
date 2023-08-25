package business

import (
	"context"
	"io"
	"os"

	log "golang.org/x/exp/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	errHandler "github.com/phllpmcphrsn/voice-quips/errors"
)

// S3Downloader defines the API for downloading objects to S3 storage
type S3Downloader interface {
    DownloadObject(ctx context.Context, objectName, bucket string) error
}

// AWSDownloader is a struct that implements the S3Downloader interface using AWS's S3 SDK for Go
type AWSDownloader struct {
    // S3 is the service client for Amazon S3
    Client *s3.Client
}

func NewAWSDownloader(client *s3.Client) *AWSDownloader {
	return &AWSDownloader{Client: client}
}

// DownloadObject uploads to an AWS bucket with the given file
func (uploader *AWSDownloader) DownloadObject(ctx context.Context, objectName, bucket string) error {
	var uploadError *errHandler.UploadError

	file, err := os.Open(objectName)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}
	
	defer file.Close()

	response, err := uploader.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key: &objectName,
	})
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	log.Debug("response from object being uploaded", "metadata", response.ResultMetadata)
	return nil
}

// MinioDownloader is a struct that implements the S3Downloader interface using MinIO's S3 SDK for Go
type MinioDownloader struct {
    // S3 is the service client for MinIO
    Client *minio.Client
}

func NewMinioDownloader(client *minio.Client) *MinioDownloader {
	return &MinioDownloader{Client: client}
}

// DownloadObject uploads to an MinIO bucket with the given file
func (uploader *MinioDownloader) DownloadObject(ctx context.Context, objectName, bucket string) error {
	var uploadError *errHandler.UploadError

	file, err := os.Open(objectName)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}
	size := file.Stat()
	defer file.Close()

	response, err := uploader.Client.PutObject(ctx, bucket, objectName, file)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	log.Debug("response from object being uploaded", "metadata", response.ResultMetadata)
	return nil
}
