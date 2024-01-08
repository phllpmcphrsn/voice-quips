package s3

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func TestAWSClient_UploadObject(t *testing.T) {
	testCases := []struct {
		name          string
		mockS3Client  *MockS3Client
		filename      string
		bucket        string
		expectedError bool
	}{
		{
			name:          "SuccessfulUpload",
			mockS3Client:  new(MockS3Client),
			filename:      "example.txt",
			bucket:        "my-bucket",
			expectedError: false,
		},
		{
			name:          "UploadError",
			mockS3Client:  new(MockS3Client),
			filename:      "example.txt",
			bucket:        "my-bucket",
			expectedError: true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &S3Client{S3Client: tc.mockS3Client}

			err := client.UploadObject(context.Background(), tc.filename, tc.bucket)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAWSClient_DownloadObject(t *testing.T) {
	testCases := []struct {
		name          string
		mockS3Client  *MockS3Client
		objectName    string
		bucket        string
		expectedError bool
	}{
		{
			name: "SuccessfulDownload",
			mockS3Client: &MockS3Client{
				On("GetObject", mock.Anything, mock.AnythingOfType("*s3.GetObjectInput")).Return(&s3.GetObjectOutput{}, nil),
			},
			objectName:    "example.txt",
			bucket:        "my-bucket",
			expectedError: false,
		},
		{
			name: "DownloadError",
			mockS3Client: &MockS3Client{
				On("GetObject", mock.Anything, mock.AnythingOfType("*s3.GetObjectInput")).Return(nil, errors.New("download failed")),
			},
			objectName:    "example.txt",
			bucket:        "my-bucket",
			expectedError: true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &S3Client{S3Client: tc.mockS3Client}

			err := client.DownloadObject(context.Background(), tc.objectName, tc.bucket)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
