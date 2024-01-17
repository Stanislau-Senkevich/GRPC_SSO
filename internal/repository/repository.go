package repository

import (
	"GRPC_SSO/internal/domain/models"
	"context"
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
}

type UserInfoRepository interface {
}
