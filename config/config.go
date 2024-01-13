package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	log "log/slog"

	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	AudioDirectory string
	API            APIConfig      `mapstructure:"api"`
	Log            Log            `mapstructure:"log"`
	Database       DatabaseConfig `mapstructure:"database"`
}

// APIConfig holds the API configuration values
type APIConfig struct {
	Address string `mapstructure:"address"`
	Path    string `mapstructure:"path"`
	Env     string `mapstructure:"env"`
}

// Log holds the log configuration values
type Log struct {
	Level string `mapstructure:"level"`
}

// DatabaseConfig holds the database configuration values
type DatabaseConfig struct {
	FileInfoConfig FileInformationStoreConfig `mapstructure:"file"`
	S3Config       S3Config                   `mapstructure:"s3"`
}

type FileInformationStoreConfig struct {
	Host        string      `mapstructure:"host"`
	Port        int         `mapstructure:"port"`
	Name        string      `mapstructure:"name"`
	SSL         SSL         `mapstructure:"ssl"`
	Credentials Credentials `mapstructure:"credentials"`
}

type S3Config struct {
	Bucket      string      `mapstructure:"bucket"`
	Region      string      `mapstructure:"region"`
	Endpoint    string      `mapstructure:"endpoint"`
	Retry       int         `mapstructure:"retry"`
	Timeout     int64       `mapstructure:"timeout"`
	SSL         SSL         `mapstructure:"ssl"`
	Credentials Credentials `mapstructure:"credentials"`
}

// SSL determines if SSL will be enabled for database connections
type SSL struct {
	Enabled bool
}

// Credentials stores persistence layer-related credentials
type Credentials struct {
	GetFromEnv  bool   `mapstructure:"envvar"`
	UserVar     string `mapstructure:"userVar"`
	PasswordVar string `mapstructure:"passwordVar"`

	User     string
	Password []byte
}

// LoadConfig loads the configuration values from the specified file.
func LoadConfig(file string) (*Config, error) {
	// Set the file name and path
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("..")
	}

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("an error occurred while reading in config file", "err", err)
		return nil, err
	}

	// Unmarshal the configuration values into the Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Error("an error occurred while unmarshaling config file", "err", err)
		return nil, err
	}

	// TODO walk the directory or list files to ensure it contains something
	// if path doesn't exist an error will be returned
	err = PathExists(config.AudioDirectory)
	if err != nil {
		log.Error("missing directory for audio files in config", "err", err, "directory", config.AudioDirectory)
		// return nil, errors.New("missing directory for audio files in config")
	}

	// TODO add if-statements checking for fields' values existence
	// TODO API if-statements
	
	// TOOD Database if-statements
	if config.Database.FileInfoConfig.Credentials.GetFromEnv {
		config.Database.FileInfoConfig.Credentials.GetCredentialsFromEnv()
	}

	if config.Database.S3Config.Credentials.GetFromEnv {
		config.Database.S3Config.Credentials.GetCredentialsFromEnv()
	}

	log.Info("unmarshalled config in viper", "config", config)
	return &config, nil
}

// GetCredentialsFromEnv will set the user and password values to envvar values
func (c *Credentials) GetCredentialsFromEnv() {
	user := os.Getenv(c.UserVar)
	if user != "" {
		c.User = user
	}

	password := os.Getenv(c.PasswordVar)
	if password != "" {
		c.Password = []byte(password)
	}
}

func PathExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		// this is where it would be a good point to notify the developer to reload the
		// config with the path
		log.Error("Couldn't find file/directory at the given path.", "err", err, "path", path)
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
