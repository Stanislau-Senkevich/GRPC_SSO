package repository

import (
	"context"
	"errors"

	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	AuthRepository
	PermissionsRepository
	UserInfoRepository
}

type AuthRepository interface {
	Login(ctx context.Context, email, passHash string) (int64, error)
	CreateUser(ctx context.Context, user *models.User) (int64, error)
}

type PermissionsRepository interface {
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type UserInfoRepository interface {
}
