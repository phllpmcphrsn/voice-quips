package file

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/phllpmcphrsn/voice-quips/config"
	log "golang.org/x/exp/slog"
)

type FileInformationRepository interface {
	FindById(context.Context, string) (*AudioFile, error)
	FindAll(context.Context) ([]*AudioFile, error)
	Create(context.Context, AudioFile) (*AudioFile, error)
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

// CreateIndexOn creates an index on the list of columns given
func (p *PostgresStore) CreateIndexOn(name string, columns []string) error {
	numberOfColumns := len(columns)
	parameter := 2 // starting at 2 since the SQL statement will already be using $1
	stmtParameters := ""

	// this is how we'll handle mulitple columns being given
	// placeholders are needed for each column so, here, we're counting
	// all of those columns and creating parameters for each
	i := 0
	for i < numberOfColumns {
		// checking if we're at the last number so that we can adjust the end of the
		// string to not have a space - purely for formatting
		if i == numberOfColumns-1 {
			stmtParameters = stmtParameters + fmt.Sprintf("$%d", parameter)
		} else {
			stmtParameters = stmtParameters + fmt.Sprintf("$%d ", parameter)
		}
		parameter++
		i++
	}
	stmt := "CREATE INDEX IF NOT EXISTS $1 ON AudioFile(" + stmtParameters + ")"
	_, err := p.db.Exec(stmt, columns)
	if err != nil {
		log.Error("Index for column(s) could not be created", "err", err)
		return IndexNotCreatedError("")
	}
	return nil
}

func (p *PostgresStore) FindAudioFileById(ctx context.Context, id string) (*AudioFile, error) {
	log.Debug("Retrieving a AudioFile record from the DB", "id", id)
	var audioFile AudioFile

	// Query for a single row
	selectStmt := "SELECT * FROM file_info WHERE id = $1"
	err := p.db.QueryRowContext(ctx, selectStmt, id).Scan(
		&audioFile.ID,
		&audioFile.Filename,
		&audioFile.FileType,
		&audioFile.FileType,
		&audioFile.Category,
		&audioFile.UploadDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NoRowsFoundError("")
		}
		return nil, NewDBError(err)
	}
	return &audioFile, nil
}

func (p *PostgresStore) FindAllAudioFile(ctx context.Context) ([]*AudioFile, error) {
	var audioFiles []*AudioFile

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

	var audioFile *AudioFile
	for rows.Next() {
		err = rows.Scan(
			&audioFile.ID,
			&audioFile.Filename,
			&audioFile.FileType,
			&audioFile.FileType,
			&audioFile.Category,
			&audioFile.UploadDate,
		)
		if err != nil {
			return nil, NewDBError(err)
		}
		audioFiles = append(audioFiles, audioFile)
	}

	// According to go.dev, one reason to check for an error is that if the results are incomplete
	// due to the overall query failing then we'll need to check for that error after the loop
	err = rows.Err()
	if err != nil {
		return nil, NewDBError(err)
	}
	return audioFiles, nil
}

func (p *PostgresStore) Create(ctx context.Context, audioFile AudioFile) (*AudioFile, error) {
	log.Debug("Inserting a AudioFile record into the DB", "record", audioFile)
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

	var savedAudioFile *AudioFile
	err := p.db.QueryRowContext(
		ctx,
		insertStmt,
		&audioFile.Filename,
		&audioFile.FileType,
		&audioFile.S3Link,
		&audioFile.Category,
		&audioFile.UploadDate,
	).Scan(&savedAudioFile)

	if err != nil {
		log.Error("An error occurred while inserting to db", "err", err)
		return nil, err
	}

	log.Debug("Successfully inserted row", "record", savedAudioFile)
	return savedAudioFile, nil
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
