package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/phllpmcphrsn/voice-quips/config"
	"github.com/phllpmcphrsn/voice-quips/models"
	log "golang.org/x/exp/slog"
)

type MetadataRepository interface {
	GetMetadataById(context.Context, string) (*models.Metadata, error)
	GetAllMetadata(context.Context) ([]*models.Metadata, error)
	CreateMetadata(context.Context, models.Metadata) (*models.Metadata, error)
	DeleteMetadata(context.Context, string) error
}

type PostgresStore struct {
	// will handle Postgres DB isntance
	db *sql.DB
}

func NewPostgresStore(config config.StorageConfig) (*PostgresStore, error) {
	// TODO think about moving to the config
	var ssl string
	if config.SSL.Enabled {
		ssl = "enabled"
	} else {
		ssl = "disabled"
	}

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=%s", config.Host, config.Name, config.Credentials.User, string(config.Credentials.Password), config.Port, ssl)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// check for conenction to the db
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) CreateTable() error {
	stmt := `CREATE TABLE IF NOT EXISTS metadata (
		id serial primary key,
		filename varchar(50),
		file_type varchar(6),
		s3_link varchar(200),
		category varchar(50),
		upload_date timestamp
	)`

	_, err := p.db.Exec(stmt)
	if err != nil {
		log.Error("An error occured while creating the metadata table", "err", err)
		return err
	}
	return nil
}

func (p *PostgresStore) GetMetadataById(context.Context, string) (*models.Metadata, error) {
	
}

func (p *PostgresStore) CreateMetadata(ctx context.Context, metadata models.Metadata) (*models.Metadata, error) {
	log.Debug("Inserting a metadata record into the DB", "record", metadata)
	insertStmt := `
	INSERT INTO metadata (
		filename,
		file_type,
		s3_link,
		category,
		upload_date
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING *`

	var savedMetadata *models.Metadata
	err := p.db.QueryRowContext(
		ctx,
		insertStmt,
		&metadata.Filename,
		&metadata.FileType,
		&metadata.S3Link,
		&metadata.Category,
		&metadata.UploadDate,
	).Scan(&savedMetadata)

	if err != nil {
		log.Error("An error occurred while inserting to db", "err", err)
		return nil, err
	}

	log.Debug("Successfully inserted row", "record", savedMetadata)
	return savedMetadata, nil
}

func (p *PostgresStore) DeleteMetadata(ctx context.Context, id string) error {
	log.Debug("Deleting metadata record from the DB", "id", id)
	deleteStmt := `	DELETE FROM metadata WHERE id=$1`

	_, err := p.db.ExecContext(ctx, deleteStmt, id)

	if err != nil {
		log.Error("An error occurred while deleting from db", "err", err, "id", id)
		return err
	}

	log.Debug("Successfully deleted row", "id", id)
	return nil
}