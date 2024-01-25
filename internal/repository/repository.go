package repository

import (
	"context"

	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
)

type Repository interface {
	AuthRepository
	PermissionsRepository
	UserInfoRepository
}

type AuthRepository interface {
	Login(ctx context.Context, email, passHash string) (models.User, error)
	CreateUser(ctx context.Context, user *models.User) (int64, error)
}

type PermissionsRepository interface {
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type UserInfoRepository interface {
	GetUserInfo(ctx context.Context, userID int64) (models.User, error)
	UpdateUserInfo(ctx context.Context, userID int64, updatedUser *models.User) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPasswordHash string) error
	AddFamily(ctx context.Context, user *models.User, familyID int64) error
	DeleteFamily(ctx context.Context, user *models.User, familyID int64) error
	DeleteUser(ctx context.Context, userID int64) error
}
