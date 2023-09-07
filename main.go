package main

import (
	"errors"
	"flag"
	"os"

	"github.com/phllpmcphrsn/voice-quips/config"
	"github.com/phllpmcphrsn/voice-quips/metadata"
	log "golang.org/x/exp/slog"
)

const (
	defaultAudioDirPath = "./audio_files/"
	audioDirPathUsage = "audio directory path (eg. '/etc/audio/)."
	defaultConfigFilePath = "./config.yml"
	configFilePathUsage = "Config file path (eg. '/etc/api/config.yml'). Config must be named 'config.yml'."
	dbUserUsage = "Username for database. If left empty, the program will look for the DBUSER environment variable"
	dbPasswordUsage = "Password for database. If left empty, the program will look for the DBPASS environment variable"
)

var configFilePath string
var audioDirPath string
var dbUser string
var dbPass string

var errDbUsernameMissing = errors.New("database username not given or found (usage: --dbuser <user> or DBUSER=<user>)")
var errDbPasswordMissing = errors.New("database password not given or found (usage: --dbpass <password> or DBPASS=<password>)")

// ensures all flag bindings occur prior to flag.Parse() being called
func init() {
	flag.StringVar(&configFilePath, "config", defaultConfigFilePath, configFilePathUsage)
	flag.StringVar(&configFilePath, "c", defaultConfigFilePath, configFilePathUsage)
	flag.StringVar(&audioDirPath, "dir", defaultAudioDirPath, audioDirPathUsage)
	flag.StringVar(&audioDirPath, "data", defaultAudioDirPath, audioDirPathUsage)
	flag.StringVar(&dbUser, "dbuser", "", dbUserUsage)
	flag.StringVar(&dbPass, "dbpass", "", dbPasswordUsage)
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

//	@host		localhost:9090
//	@BasePath	/api/v1
func main() {	
	// could place this in init() but it'll cause errors for tests
	// error: "flag provided but not defined"
	flag.Parse()

	var err error

	// if credentials aren't given as args, look for them in the env
	if dbUser == "" && os.Getenv("DBUSER") == "" {
		log.Error(errDbUsernameMissing.Error())
		panic(errDbUsernameMissing)
	}
	if dbPass == "" && os.Getenv("DBPASS") == "" {
		log.Error(errDbPasswordMissing.Error())
		panic(errDbPasswordMissing)
	}

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Error("There was an issue loading the config file", "err", err)
		panic(err)
	}

	setLogger(config.Log.Level)
	
	store, err := metadata.NewPostgresStore(config.Database.MetadataStore)
	if err != nil {
		log.Error("There was an issue reaching the database", "err", err)
		panic(err)
	}
	log.Info("Connected to database...")

	if err := store.CreateTable(); err != nil {
		log.Error("There was an issue creating the database table", "err", err)
		panic(err)
	}
	
	indexedColumns := []string{"file_type", "category"}
	indexName := "file_type_and_category_index"
	if err := store.CreateIndexOn(indexName, indexedColumns)

	// want to check if table has any elements prior to read and populating from csv
	// if it does we'll assume that it's already been populated with data from csv
	count, err := store.Count()
	if err != nil {
		log.Error("An error occured while checking for table's count", "err", err)
		panic(err)
	} else if count == 0 {
		log.Info("Populating cars table...")
		go readCsv(store)
	}

	api := NewAPIServer(store, config.API.Address, config.Env)
	api.StartRouter()
}
