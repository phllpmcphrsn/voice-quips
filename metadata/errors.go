package metadata

import "errors"

var ErrDbUsernameMissing = errors.New("database username not given or found (usage: --dbuser <user> or DBUSER=<user>)")
var ErrDbPasswordMissing = errors.New("database password not given or found (usage: --dbpass <password> or DBPASS=<password>)")

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
