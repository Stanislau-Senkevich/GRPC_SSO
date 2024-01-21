package services

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
)

type Services interface {
	Auth
	Permissions
	UserInfo
}

type Auth interface {
	SignIn(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, user *models.User) (int64, error)
}

type Permissions interface {
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type UserInfo interface {
	GetUserInfo(ctx context.Context) (models.User, error)
	GetUserInfoByID(ctx context.Context, userID int64) (models.User, error)
	UpdateUserInfo(ctx context.Context, updatedUser *models.User) error
	ChangePassword(ctx context.Context, oldPassword, newPasswordHash string) error
	DeleteUser(ctx context.Context, userID int64) error
}
