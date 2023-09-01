package metadata

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/phllpmcphrsn/voice-quips/config"
	errHandler "github.com/phllpmcphrsn/voice-quips/errors"
	log "golang.org/x/exp/slog"
)

type MetadataRepository interface {
	FindById(context.Context, string) (*Metadata, error)
	FindAll(context.Context) ([]*Metadata, error)
	Create(context.Context, Metadata) (*Metadata, error)
	Delete(context.Context, string) error
}

type PostgresStore struct {
	// will handle Postgres DB instance
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

// CreateIndexOn creates an index on the list of columns given
func (p *PostgresStore) CreateIndexOn(name string, columns []string) error {
	numberOfColumns := len(columns)
	parameter := 2   // starting at 2 since the SQL statement will already be using $1
	stmtParameters := ""
	
	// this is how we'll handle mulitple columns being given
	// placeholders are needed for each column so, here, we're counting
	// all of those columns and creating parameters for each
	i := 0
	for i < numberOfColumns {
		// checking if we're at the last number so that we can adjust the end of the
		// string to not have a space - purely for formatting
		if i == numberOfColumns - 1 {
			stmtParameters = stmtParameters + fmt.Sprintf("$%d", parameter)
		} else {
			stmtParameters = stmtParameters + fmt.Sprintf("$%d ", parameter)
		}
		parameter++
		i++
	}
	stmt := "CREATE INDEX IF NOT EXISTS $1 ON metadata(" + stmtParameters + ")"
	_, err := p.db.Exec(stmt, columns)
	if err != nil {
		log.Error("Index for column(s) could not be created", "err", err)
		return errHandler.IndexNotCreatedError("")
	}
	return nil
}

func (p *PostgresStore) FindMetadataById(ctx context.Context, id string) (*Metadata, error) {
	log.Debug("Retrieving a metadata record from the DB", "id", id)
	var metadata Metadata

	// Query for a single row
	selectStmt := "SELECT * FROM metadata WHERE id = $1"
	err := p.db.QueryRowContext(ctx, selectStmt, id).Scan(
		&metadata.ID,
		&metadata.Filename,
		&metadata.FileType,
		&metadata.FileType,
		&metadata.Category,
		&metadata.UploadDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errHandler.NoRowsFoundError("")
		}
		return nil, errHandler.NewDBError(err)
	}
	return &metadata, nil
}

func (p *PostgresStore) FindAllMetadata(ctx context.Context) ([]*Metadata, error) {
	var metadatum []*Metadata

	// Query for all rows
	selectStmt := "SELECT * FROM metadata"
	rows, err := p.db.QueryContext(ctx, selectStmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errHandler.NoRowsFoundError("")
		}
		return nil, errHandler.NewDBError(err)
	}
	defer rows.Close()

	var metadata *Metadata
	for rows.Next() {
		err = rows.Scan(
			&metadata.ID,
			&metadata.Filename,
			&metadata.FileType,
			&metadata.FileType,
			&metadata.Category,
			&metadata.UploadDate,
		)
		if err != nil {
			return nil, errHandler.NewDBError(err)
		}
		metadatum = append(metadatum, metadata)
	}

	// According to go.dev, one reason to check for an error is that if the results are incomplete 
	// due to the overall query failing then we'll need to check for that error after the loop
	err = rows.Err()
	if err != nil {
		return nil, errHandler.NewDBError(err)
	}
	return metadatum, nil
}

func (p *PostgresStore) Create(ctx context.Context, metadata Metadata) (*Metadata, error) {
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

	var savedMetadata *Metadata
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

func (p *PostgresStore) Delete(ctx context.Context, id string) error {
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