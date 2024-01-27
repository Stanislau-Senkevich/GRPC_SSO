package grpcerror

import "errors"

var (
	ErrInternalError   = errors.New("internal error")
	ErrUserNotInFamily = errors.New("user already not in the family")
	ErrUserInFamily    = errors.New("user already in the family")
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNoToken         = errors.New("authorization token was not provided")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenClaims     = errors.New("failed to get token claims")
	ErrForbidden       = errors.New("forbidden")
)
