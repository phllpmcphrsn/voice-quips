package main

import "errors"

var errDbUsernameMissing = errors.New("database username not given or found (usage: --dbuser <user> or DBUSER=<user>)")
var errDbPasswordMissing = errors.New("database password not given or found (usage: --dbpass <password> or DBPASS=<password>)")
// var errDbCarNotFound = errors.New()
type APIError struct {
	ErrorCode    int
	ErrorMessage string
}