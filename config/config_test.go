package config

import (
	"fmt"
	"os"
	"testing"

	log "log/slog"

	"github.com/stretchr/testify/assert"
)

const testConfigDir string = "./test"
const emptyConfigFileName string = "/empty_config.yml"
const validConfigFileName string = "/valid_config.yml"

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name          string
		configFile    string
		envVars       map[string]string
		expectedError bool
	}{
		{
			name:          "ValidConfigFile",
			configFile:    fmt.Sprintf("%s%s", testConfigDir, validConfigFileName),
			expectedError: false,
		},
		{
			name:          "InvalidConfigFile",
			configFile:    fmt.Sprintf("%s%s", testConfigDir, emptyConfigFileName),
			expectedError: true,
		},
		{
			name:          "NoConfigFileGiven",
			configFile:    "",
			expectedError: false,
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

func TestLoadConfigWithNoFileGiven(t *testing.T) {
	testCases := []struct {
		name          string
		configFile    string
		expectedError bool
	}{
		{
			name:          "NoFileGiven_ConfigYMLExists",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load config
			_, err := LoadConfig("")

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
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
			path:          "./test/valid_config.yml",
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

// func TestHandleCredentials(t *testing.T) {
// 	testCases := []struct {
// 		name         string
// 		store        FileInformationStoreConfig
// 		envVars      map[string]string
// 		expectedUser string
// 		expectedPass []byte
// 	}{
// 		{
// 			name: "CredentialsEnabled",
// 			store: FileInformationStoreConfig{
// 				Credentials: Credentials{
// 					UserVar:     "TEST_USER",
// 					PasswordVar: "TEST_PASS",
// 					User:        "defaultUser",
// 					Password:    []byte("defaultPass"),
// 				},
// 			},
// 			envVars: map[string]string{
// 				"TEST_USER": "testUser",
// 				"TEST_PASS": "testPass",
// 			},
// 			expectedUser: "testUser",
// 			expectedPass: []byte("testPass"),
// 		},
// 		{
// 			name: "CredentialsDisabled",
// 			store: FileInformationStoreConfig{
// 				Credentials: Credentials{
// 					UserVar: "TEST_USER",
// 				},
// 			},
// 			envVars:      map[string]string{},
// 			expectedUser: "defaultUser",
// 			expectedPass: []byte("defaultPass"),
// 		},
// 		// Add more test cases as needed
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Set environment variables for test case
// 			for key, value := range tc.envVars {
// 				os.Setenv(key, value)
// 				defer os.Unsetenv(key)
// 			}

// 			// Handle credentials
// 			handleCredentials(&tc.store)

// 			assert.Equal(t, tc.expectedUser, tc.store.Credentials.UserVar)
// 			assert.Equal(t, tc.expectedPass, tc.store.Credentials.Password)
// 		})
// 	}
// }

