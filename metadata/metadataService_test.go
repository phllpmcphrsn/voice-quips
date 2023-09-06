package metadata

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetadataRepository struct {
	mock.Mock
}

func (m *MockMetadataRepository) Create(ctx context.Context, metadata Metadata) (*Metadata, error) {
	args := m.Called(ctx, metadata)
	return args.Get(0).(*Metadata), args.Error(1)
}

func (m *MockMetadataRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMetadataRepository) FindById(ctx context.Context, id string) (*Metadata, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Metadata), args.Error(1)
}

func (m *MockMetadataRepository) FindAll(ctx context.Context) ([]*Metadata, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Metadata), args.Error(1)
}

func TestMetadataService_SaveMetadata(t *testing.T) {
	testCases := []struct {
		name             string
		mockRepository   *MockMetadataRepository
		inputMetadata    Metadata
		expectedMetadata *Metadata
		expectedError    bool
		returnedError	 error
	}{
		{
			name: "SuccessfulSave",
			mockRepository: new(MockMetadataRepository),
			inputMetadata:    Metadata{ID: 123, Filename: "TestMetadata"},
			expectedMetadata: &Metadata{ID: 123, Filename: "TestMetadata"},
			expectedError:    false,
			returnedError: nil,
		},
		{
			name: "SaveError",
			mockRepository: new(MockMetadataRepository),
			inputMetadata:    Metadata{ID: 123, Filename: "TestMetadata"},
			expectedMetadata: nil,
			expectedError:    true,
			returnedError: &DBError{},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			service := MetadataService{repo: tc.mockRepository}
			tc.mockRepository.On("Create", ctx, tc.inputMetadata).Return(tc.expectedMetadata, tc.returnedError)
			
			result, err := service.Save(ctx, tc.inputMetadata)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMetadata, result)
			}
		})
	}
}

// Write similar tests for DeleteMetadata, FindMetadataById, and GetAllMetadata
