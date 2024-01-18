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
}

type UserInfo interface {
}
