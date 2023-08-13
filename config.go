package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	log "golang.org/x/exp/slog"

	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	Env      string
	AudioDirectory string
	API      APIConfig
	Log      LogLevel
	Database DatabaseConfig
}

// APIConfig holds the API configuration values
type APIConfig struct {
	Address string
}

// LogLevel holds the log configuration values
type LogLevel struct {
	LevelStr string
	Level log.Level
}

// DatabaseConfig holds the database configuration values
type DatabaseConfig struct {
	MetadataStore	StorageConfig
	BlobStore	StorageConfig
}

type StorageConfig struct {
	Host string
	Port int
	Name string
	SSL SSL
	Credentials Credentials
}

// SSL determines if SSL will be enabled for database connections
type SSL struct {
	Enabled bool
}

// Credentials stores persistence layer-related credentials
type Credentials struct {
	Enabled bool
	UserVar string
	PasswordVar string

	User string
	Password []byte
}

// LoadConfig loads the configuration values from the specified file.
func LoadConfig(file string) (*Config, error) {
	// Set the file name and path
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		// If no file specified, look for the default file in the current directory
		dir, err := os.Getwd()
		if err != nil {
			log.Error("Failed to get current directory", "err", err, "")
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(dir)
	}

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Unmarshal the configuration values into the Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	// if path doesn't exist an error will be returned
	err = PathExists(config.AudioDirectory)
	if err != nil {
		return nil, err
	}
	
	if config.Database.BlobStore.Credentials.Enabled {
		handleCredentials(&config.Database.BlobStore)
	}

	if config.Database.MetadataStore.Credentials.Enabled {
		handleCredentials(&config.Database.MetadataStore)
	}

	// Get a valid slog log level
	config.Log.Level = GetLogLevel(config.Log.LevelStr)

	return &config, nil
}

// handleCredentials will set the database config's 
func handleCredentials(store *StorageConfig) {
	user := os.Getenv(store.Credentials.UserVar)
	if user != "" {
		store.Credentials.UserVar = user
	}
	
	password := os.Getenv(store.Credentials.PasswordVar)
	if password != "" {
		store.Credentials.Password = []byte(password)
	}
}

func PathExists(path string) error {
	if _, err := os.Stat(path); err == nil {	
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		// this is where it would be a good point to notify the developer to reload the
		// config with the path
		log.Error("Couldn't find file/directory at the given path.", "err", err)
		// may need to see about running a watch on the config file so that the user can
		// do a hot reload
		return err
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Error("Unknown error occurred while checking for path's existence", "err", err) 
		return err
	}
}

// GetConfigFilePath returns the absolute path of the config file based on the current directory.
func GetConfigFilePath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current directory", "err", err)
	}
	return filepath.Join(dir, "config.yaml")
}


// GetLogLevel returns the slog log level based on a string representation of the log level.
// INFO is used as the default in case a level isn't given or is unexpected
func GetLogLevel(logLevel string) log.Level {
	level := strings.ToLower(logLevel)
	print(level)
	log.Info(level)
	switch level {
	case "debug":
		return log.LevelDebug
	case "info":
		return log.LevelInfo
	case "warn":
		return log.LevelWarn
	case "error":
		return log.LevelError
	default:
		log.Warn("Supported log levels are: debug, info, warn, and error")
		return log.LevelInfo
	}
}
