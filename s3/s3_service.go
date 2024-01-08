package s3

import (
	"context"
	"os"
	"path/filepath"

	log "log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
)

// Uploader defines the API for uploading objects to S3 storage
type Uploader interface {
	UploadObject(ctx context.Context, filename, bucket string) error
}

// Downloader defines the API for downloading objects to S3 storage
type Downloader interface {
	DownloadObject(ctx context.Context, objectName, bucket string) ([]byte, error)
}

type DownloadUploader interface {
	Uploader
	Downloader
}

// S3Client is a struct that implements the DownloadUploader interface using AWS's S3 SDK for Go
type S3Client struct {
	// S3Client is the service client for Amazon S3
	S3Client *s3.Client
}

func New(client *s3.Client) *S3Client {
	return &S3Client{client}
}

// UploadObject uploads to an AWS bucket with the given file
func (a *S3Client) UploadObject(ctx context.Context, filename, bucket string) error {
	var uploadError *UploadError

	file, err := os.Open(filename)
	if err != nil {
		uploadError.Err = err
		return uploadError
	}

	defer file.Close()

	// want to determine the file extension so that we can set the approriate content-type header
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
func (a *S3Client) DownloadObject(ctx context.Context, objectName, bucket string) error {
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

// MinioClient is a struct that implements the DownloadUploader interface using MinIO's S3 SDK for Go
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

// DownloadObject downloads an object from a MinIO bucket with the given object name
// TODO this needs to return the file/data received
func (m *MinioClient) DownloadObject(ctx context.Context, objectName, bucket string) ([]byte, error) {
	var downloadError *DownloadError
	data := make([]byte, 1)
	file, err := os.Open(objectName)
	if err != nil {
		downloadError.Err = err
		return data, downloadError
	}
	defer file.Close()


	opts := minio.GetObjectOptions{}
	response, err := m.S3Client.GetObject(ctx, bucket, objectName, opts)
	if err != nil {
		downloadError.Err = err
		return data, downloadError
	}
	
	defer response.Close()
	log.Debug("completed download reqeust", "bucket", bucket, "file", objectName)
	// TODO does minio have easily accessible metadata on downloads
	return data, nil
}
