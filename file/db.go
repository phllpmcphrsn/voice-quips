package file

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	log "log/slog"

	_ "github.com/lib/pq"
	"github.com/phllpmcphrsn/voice-quips/config"
)

type FileInformationRepository interface {
	FindById(context.Context, string) (*FileInformation, error)
	FindAll(context.Context) ([]*FileInformation, error)
	Create(context.Context, FileInformation) (*FileInformation, error)
	Delete(context.Context, string) error
}

type PostgresStore struct {
	// will handle Postgres DB instance
	db *sql.DB
}

func NewPostgresStore(config config.FileInformationStoreConfig) (*PostgresStore, error) {
	// TODO think about moving to the config
	var ssl string
	if config.SSL.Enabled {
		ssl = "require"
	} else {
		ssl = "disable"
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
	stmt := `CREATE TABLE IF NOT EXISTS file_info (
		id serial primary key,
		filename varchar(50),
		file_type varchar(6),
		s3_link varchar(200),
		category varchar(50),
		title varchar(100),
		artist varchar(75),
		album varchar(100),
		year smallint,
		upload_date timestamp
	)`

	_, err := p.db.Exec(stmt)
	if err != nil {
		log.Error("An error occured while creating the file_info table", "err", err)
		return err
	}
	return nil
}

// CreateIndexOn creates an index on the list of columns given. The name provided will
// be 
func (p *PostgresStore) CreateIndexOn(name string, columns []string) error {
	stmtParameters := strings.Join(columns, ",")

	stmt := "CREATE INDEX IF NOT EXISTS " + name + " ON file_info(" + stmtParameters + ")"
	_, err := p.db.Exec(stmt)
	if err != nil {
		log.Error("Index for column(s) could not be created", "err", err)
		return IndexNotCreatedError("")
	}
	return nil
}

func (p *PostgresStore) FindById(ctx context.Context, id string) (*FileInformation, error) {
	log.Debug("Retrieving a file_info record from the DB", "id", id)
	var fileInformation FileInformation

	// Query for a single row
	selectStmt := "SELECT * FROM file_info WHERE id = $1"
	err := p.db.QueryRowContext(ctx, selectStmt, id).Scan(
		&fileInformation.ID,
		&fileInformation.Filename,
		&fileInformation.FileType,
		&fileInformation.FileType,
		&fileInformation.Category,
		&fileInformation.UploadDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NoRowsFoundError("")
		}
		return nil, NewDBError(err)
	}
	return &fileInformation, nil
}

func (p *PostgresStore) FindAll(ctx context.Context) ([]*FileInformation, error) {
	var fileInformations []*FileInformation

	// Query for all rows
	selectStmt := "SELECT * FROM file_info"
	rows, err := p.db.QueryContext(ctx, selectStmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NoRowsFoundError("")
		}
		return nil, NewDBError(err)
	}
	defer rows.Close()

	var fileInformation *FileInformation
	for rows.Next() {
		err = rows.Scan(
			&fileInformation.ID,
			&fileInformation.Filename,
			&fileInformation.FileType,
			&fileInformation.FileType,
			&fileInformation.Category,
			&fileInformation.UploadDate,
		)
		if err != nil {
			return nil, NewDBError(err)
		}
		fileInformations = append(fileInformations, fileInformation)
	}

	// According to go.dev, one reason to check for an error is that if the results are incomplete
	// due to the overall query failing then we'll need to check for that error after the loop
	err = rows.Err()
	if err != nil {
		return nil, NewDBError(err)
	}
	return fileInformations, nil
}

func (p *PostgresStore) Create(ctx context.Context, fileInformation FileInformation) (*FileInformation, error) {
	log.Debug("Inserting a file_info record into the DB", "record", fileInformation)
	insertStmt := `
	INSERT INTO file_info (
		filename,
		file_type,
		s3_link,
		category,
		upload_date
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING *`

	var savedfileInformation *FileInformation
	err := p.db.QueryRowContext(
		ctx,
		insertStmt,
		&fileInformation.Filename,
		&fileInformation.FileType,
		&fileInformation.S3Link,
		&fileInformation.Category,
		&fileInformation.UploadDate,
	).Scan(&savedfileInformation)

	if err != nil {
		log.Error("An error occurred while inserting to db", "err", err)
		return nil, err
	}

	log.Debug("Successfully inserted row", "record", savedfileInformation)
	return savedfileInformation, nil
}

func (p *PostgresStore) Delete(ctx context.Context, id string) error {
	log.Debug("Deleting audio file record from the DB", "id", id)
	deleteStmt := `DELETE FROM file_info WHERE id=$1`

	_, err := p.db.ExecContext(ctx, deleteStmt, id)

	if err != nil {
		log.Error("An error occurred while deleting from db", "err", err, "id", id)
		return err
	}

	log.Debug("Successfully deleted row", "id", id)
	return nil
}
