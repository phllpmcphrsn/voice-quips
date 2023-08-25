package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	log "golang.org/x/exp/slog"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name          string
		configFile    string
		envVars       map[string]string
		expectedError bool
	}{
		{
			name:          "ValidConfigFile",
			configFile:    "./test/valid_config.yml",
			expectedError: false,
		},
		{
			name:          "InvalidConfigFile",
			configFile:    "./test/empty_config.yml",
			expectedError: true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables for test case
			for key, value := range tc.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Load config
			_, err := LoadConfig(tc.configFile)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleCredentials(t *testing.T) {
	testCases := []struct {
		name         string
		store        StorageConfig
		envVars      map[string]string
		expectedUser string
		expectedPass []byte
	}{
		{
			name: "CredentialsEnabled",
			store: StorageConfig{
				Credentials: Credentials{
					Enabled:     true,
					UserVar:     "TEST_USER",
					PasswordVar: "TEST_PASS",
					User:        "defaultUser",
					Password:    []byte("defaultPass"),
				},
			},
			envVars: map[string]string{
				"TEST_USER": "testUser",
				"TEST_PASS": "testPass",
			},
			expectedUser: "testUser",
			expectedPass: []byte("testPass"),
		},
		{
			name: "CredentialsDisabled",
			store: StorageConfig{
				Credentials: Credentials{
					Enabled: false,
					UserVar: "TEST_USER",
				},
			},
			envVars:      map[string]string{},
			expectedUser: "defaultUser",
			expectedPass: []byte("defaultPass"),
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables for test case
			for key, value := range tc.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Handle credentials
			handleCredentials(&tc.store)

			assert.Equal(t, tc.expectedUser, tc.store.Credentials.UserVar)
			assert.Equal(t, tc.expectedPass, tc.store.Credentials.Password)
		})
	}
}

func TestPathExists(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		expectedError bool
	}{
		{
			name:          "ExistingPath",
			path:          "testdata",
			expectedError: false,
		},
		{
			name:          "NonExistingPath",
			path:          "nonexistent",
			expectedError: true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := PathExists(tc.path)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	testCases := []struct {
		name           string
		logLevel       string
		expectedOutput log.Level
	}{
		{
			name:           "DebugLogLevel",
			logLevel:       "debug",
			expectedOutput: log.LevelDebug,
		},
		{
			name:           "InfoLogLevel",
			logLevel:       "info",
			expectedOutput: log.LevelInfo,
		},
		{
			name:           "WarnLogLevel",
			logLevel:       "warn",
			expectedOutput: log.LevelWarn,
		},
		{
			name:           "ErrorLogLevel",
			logLevel:       "error",
			expectedOutput: log.LevelError,
		},
		{
			name:           "InvalidLogLevel",
			logLevel:       "invalid",
			expectedOutput: log.LevelInfo,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := GetLogLevel(tc.logLevel)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
