package api

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
