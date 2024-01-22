package grpcerror

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNoToken         = errors.New("authorization token was not provided")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenClaims     = errors.New("failed to get token claims")
	ErrForbidden       = errors.New("forbidden")
)
