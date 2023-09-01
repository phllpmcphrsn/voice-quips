package errors

import (
	"errors"
	"fmt"
)

var ErrDbUsernameMissing = errors.New("database username not given or found (usage: --dbuser <user> or DBUSER=<user>)")
var ErrDbPasswordMissing = errors.New("database password not given or found (usage: --dbpass <password> or DBPASS=<password>)")

type APIError struct {
	StatusCode int
	Err        error
}

func NewAPIError(code int, err error) *APIError {
	return &APIError{StatusCode: code, Err: err}
}

func (ae *APIError) Error() string {
	return fmt.Sprintf("status %d: err %v", ae.StatusCode, ae.Err)
}

// TODO determine if we actually need to return a pointer
func InternalServerError(message string) *APIError {
	if message == "" {
		return &APIError{StatusCode: 500, Err: errors.New("an issue occurred server-side")}
	}
	return &APIError{StatusCode: 500, Err: errors.New(message)}
}

type UploadError struct {
	Err error
}

func (ue *UploadError) Error() string {
	return "an error occurred while uploading to storage: " + ue.Err.Error()
}

type DBError struct {
	Err error
}

func NewDBError(err error) *DBError {
	return &DBError{err}
}

func (de *DBError) Error() string {
	return "an error occured while interacting with the database" + de.Err.Error()
}

// Common DB errors
func NoRowsFoundError(message string) *DBError {
	return DBErrorMessage(message)
}

func DuplicateKeyError(message string) *DBError {
	return DBErrorMessage(message)
}

func IndexNotCreatedError(message string) *DBError {
	return DBErrorMessage(message)
}

func DBErrorMessage(message string) *DBError {
	if message == "" {
		return &DBError{errors.New(message)}
	}
	return &DBError{errors.New(message)}
}
