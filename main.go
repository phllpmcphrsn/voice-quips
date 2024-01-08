package main

import (
	"flag"
	"os"

	log "log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/phllpmcphrsn/voice-quips/api"
	"github.com/phllpmcphrsn/voice-quips/config"
	"github.com/phllpmcphrsn/voice-quips/file"
	"github.com/phllpmcphrsn/voice-quips/s3"
)

const (
	defaultConfigFilePath = "./config.yml"
	configFilePathUsage   = "Config file path (eg. '/etc/api/config.yml'). Config must be named 'config.yml'."
)

var configFilePath string

// ensures all flag bindings occur prior to flag.Parse() being called
func init() {
	flag.StringVar(&configFilePath, "config", defaultConfigFilePath, configFilePathUsage)
	flag.StringVar(&configFilePath, "c", defaultConfigFilePath, configFilePathUsage)
}

func setLogger(level log.Level) {
	logger := log.New(log.NewJSONHandler(os.Stdout, &log.HandlerOptions{Level: level}))
	log.SetDefault(logger)
}

//	@title			Kaggle 2023 Car Models API
//	@version		1.0
//	@description	REST API for Kaggle 2023 Car Models Dataset which can be found here
//	@description	https://www.kaggle.com/datasets/peshimaammuzammil/2023-car-model-dataset-all-data-you-need?resource=download
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://github.com/phllpmcphrsn/KaggleCarAPI/issues
//	@contact.email	phllpmcphrsn@yahoo.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:9090
// @BasePath	/api/v1
func main() {
	// could place this in init() but it'll cause errors for tests
	// error: "flag provided but not defined"
	flag.Parse()

	var err error

	if configFilePath == "" {
		log.Warn("checking default config path since one was not given via a flag")
	}

	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Error("There was an issue loading the config file", "err", err)
		panic(err)
	}
	
	logLevel := config.GetLogLevel(cfg.Log.Level)
	setLogger(logLevel)

	// initialize database and service for file information
	store, err := initDB(cfg.Database.FileInfoConfig)
	if err != nil {
		panic(err)
	}

	fileService := file.NewFileInformationService(store)

	// initialize s3 client and service
	// TODO figure out a better or more extensible way to define a client
	// some layer should be in front of the s3 client such that I'm not coupling
	// a client here. in other words, I need some client abstraction
	client, err := initS3Client(&cfg.Database.S3Config)
	if err != nil {
		panic(err)
	}

	s3Service := s3.NewMinioClient(client)

	server := api.NewAPIServer(cfg.API, s3Service, fileService)
	server.StartRouter()
}

func initDB(cfg config.FileInformationStoreConfig) (*file.PostgresStore, error) {
	store, err := file.NewPostgresStore(cfg)
	if err != nil {
		log.Error("There was an issue reaching the database", "err", err)
		return nil, err
	}
	log.Info("Connected to database...")

	err = store.CreateTable()
	if err != nil {
		log.Error("There was an issue creating the database table", "err", err)
		return nil, err
	}

	indexedColumns := []string{"file_type", "category"}
	indexName := "file_type_and_category_index"
	err = store.CreateIndexOn(indexName, indexedColumns)
	if err != nil {
		log.Error("There was an issue creating the database index. Index will need to be created manually", "err", err)
	}
	
	return store, nil
}

// initialize the s3 client
func initS3Client(cfg *config.S3Config) (*minio.Client, error) {
	accessKey := cfg.Credentials.User
	secretKey := cfg.Credentials.Password
	endpoint := cfg.Endpoint

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, string(secretKey), ""),
		Secure: cfg.SSL.Enabled,
	})
	if err != nil {
		log.Error("There was an issue initializing the MinIO client", "err", err)
		return nil, err
	}

	return minioClient, nil
}
