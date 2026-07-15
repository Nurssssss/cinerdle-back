package errors

import "errors"

var ErrUserNotFound = errors.New("User not found")
var ErrInvalidCredentials = errors.New("Incorrect login information")
var ErrEmailAlreadyExists = errors.New("The user already exists")
var ErrNickNameAlreadyExists = errors.New("Nickname already exists")
var ErrCodeVerififactionExpired = errors.New("Code verification expired")
var ErrInvalidInput = errors.New("Invalid request from the client")
var ErrInternalServer = errors.New("Internal server error")
