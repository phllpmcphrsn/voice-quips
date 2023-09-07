package file

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileInformationRepository struct {
	mock.Mock
}

func (m *MockFileInformationRepository) Create(ctx context.Context, audioFile AudioFile) (*AudioFile, error) {
	args := m.Called(ctx, audioFile)
	return args.Get(0).(*AudioFile), args.Error(1)
}

func (m *MockFileInformationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileInformationRepository) FindById(ctx context.Context, id string) (*AudioFile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*AudioFile), args.Error(1)
}

func (m *MockFileInformationRepository) FindAll(ctx context.Context) ([]*AudioFile, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*AudioFile), args.Error(1)
}

func TestAudioFileService_SaveAudioFile(t *testing.T) {
	testCases := []struct {
		name              string
		mockRepository    *MockFileInformationRepository
		inputAudioFile    AudioFile
		expectedAudioFile *AudioFile
		expectedError     bool
		returnedError     error
	}{
		{
			name:              "SuccessfulSave",
			mockRepository:    new(MockFileInformationRepository),
			inputAudioFile:    AudioFile{ID: 123, Filename: "TestAudioFile"},
			expectedAudioFile: &AudioFile{ID: 123, Filename: "TestAudioFile"},
			expectedError:     false,
			returnedError:     nil,
		},
		{
			name:              "SaveError",
			mockRepository:    new(MockFileInformationRepository),
			inputAudioFile:    AudioFile{ID: 123, Filename: "TestAudioFile"},
			expectedAudioFile: nil,
			expectedError:     true,
			returnedError:     &DBError{},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			service := AudioFileService{repo: tc.mockRepository}
			tc.mockRepository.On("Create", ctx, tc.inputAudioFile).Return(tc.expectedAudioFile, tc.returnedError)

			result, err := service.Save(ctx, tc.inputAudioFile)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAudioFile, result)
			}
		})
	}
}

// Write similar tests for DeleteAudioFile, FindAudioFileById, and GetAllAudioFile
