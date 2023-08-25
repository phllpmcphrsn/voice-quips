package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

type MockS3Uploader struct {
	UploadObjectFn func(ctx context.Context, filename, bucket string) error
}

func (m *MockS3Uploader) UploadObject(ctx context.Context, filename, bucket string) error {
	return m.UploadObjectFn(ctx, filename, bucket)
}

func TestAWSUploader_UploadObject(t *testing.T) {
	testCases := []struct {
		name           string
		mockUploader   *MockS3Uploader
		filename       string
		bucket         string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "SuccessfulUpload",
			mockUploader: &MockS3Uploader{
				UploadObjectFn: func(ctx context.Context, filename, bucket string) error {
					return nil
				},
			},
			filename:      "example.txt",
			bucket:        "my-bucket",
			expectedError: false,
		},
		{
			name: "UploadError",
			mockUploader: &MockS3Uploader{
				UploadObjectFn: func(ctx context.Context, filename, bucket string) error {
					return errors.New("upload failed")
				},
			},
			filename:       "example.txt",
			bucket:         "my-bucket",
			expectedError:  true,
			expectedErrMsg: "upload failed",
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uploader := &AWSUploader{S3Client: nil, MetadataStore: nil}
			uploader.S3Uploader = tc.mockUploader

			err := uploader.UploadObject(context.Background(), tc.filename, tc.bucket)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Write similar tests for MinioUploader

