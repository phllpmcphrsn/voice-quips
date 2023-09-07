package s3

import (
	"context"
	"os"
	"path/filepath"

	log "golang.org/x/exp/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
)

// S3Uploader defines the API for uploading objects to S3 storage
type S3Uploader interface {
	UploadObject(ctx context.Context, filename, bucket string) error
}

// S3Downloader defines the API for downloading objects to S3 storage
type S3Downloader interface {
	DownloadObject(ctx context.Context, objectName, bucket string) error
}

type S3DownloadUploader interface {
	S3Uploader
	S3Downloader
}

// AWSClient is a struct that implements the S3DownloadUploader interface using AWS's S3 SDK for Go
type AWSClient struct {
	// S3Client is the service client for Amazon S3
	S3Client *s3.Client
}

func NewAWSClient(client *s3.Client) *AWSClient {
	return &AWSClient{client}
}

// UploadObject uploads to an AWS bucket with the given file
func (a *AWSClient) UploadObject(ctx context.Context, filename, bucket string) error {
	var uploadError *UploadError

	file, err := os.Open(filename)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	defer file.Close()

	// Want to determine the file extension so that we can set the approriate content-type header
	extension := filepath.Ext(filename)
	contentType := GetContentType(extension)

	response, err := a.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &filename,
		Body:        file,
		ContentType: &contentType,
	})
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	log.Debug("response from object being uploaded", "metadata", response.ResultMetadata)
	return nil
}

// DownloadObject downloads from the given file an AWS bucket
func (a *AWSClient) DownloadObject(ctx context.Context, objectName, bucket string) error {
	var downloadError *DownloadError

	file, err := os.Open(objectName)
	if err != nil {
		downloadError.Err = err
		return downloadError
	}

	defer file.Close()

	response, err := a.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objectName,
	})
	if err != nil {
		downloadError.Err = err
		return downloadError
	}

	log.Debug("response from object being uploaded", "metadata", response.ResultMetadata)
	return nil
}

// MinioClient is a struct that implements the S3DownloadUploader interface using MinIO's S3 SDK for Go
type MinioClient struct {
	// S3Client is the service client for MinIO
	S3Client *minio.Client
}

func NewMinioClient(client *minio.Client) *MinioClient {
	return &MinioClient{S3Client: client}
}

// UploadObject uploads to an MinIO bucket with the given file
func (uploader *MinioClient) UploadObject(ctx context.Context, filename, bucket string) error {
	var uploadError *UploadError

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

// DownloadObject uploads to an MinIO bucket with the given file
// TODO this needs to return the file/data received
func (m *MinioClient) DownloadObject(ctx context.Context, objectName, bucket string) error {
	var downloadError *DownloadError

	file, err := os.Open(objectName)
	if err != nil {
		downloadError.Err = err
		return downloadError
	}
	defer file.Close()

	size, err := file.Stat()
	if err != nil {
		downloadError.Err = err
		return downloadError
	}

	response, err := m.S3Client.PutObject(ctx, bucket, objectName, file, size.Size(), minio.PutObjectOptions{})
	if err != nil {
		downloadError.Err = err
		return downloadError
	}

	log.Debug("response from object being uploaded", "metadata", response)
	return nil
}
