package app

import (
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_SSO/internal/app/grpc"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	jwtmanager "github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository/mongodb"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/auth"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/permissions"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/userinfo"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

// New creates a new instance of the application with the provided configuration and dependencies.
func New(
	log *slog.Logger,
	grpcPort int,
	cfg *config.Config,
	tokenTTL time.Duration,
) *App {
	repo, err := mongodb.InitMongoRepository(&cfg.Mongo, log)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository: %w", err))
	}

	jwtManager := jwtmanager.New(cfg.SigningKey, tokenTTL)

	authService := auth.New(log, repo, jwtManager, cfg.HashSalt)
	permService := permissions.New(log, repo)
	userInfoService := userinfo.New(log, repo, jwtManager, cfg.HashSalt)

	accessibleRoles := map[string][]string{
		"/userinfo.UserInfo/GetUserInfo":     {"user", "admin"},
		"/userinfo.UserInfo/UpdateUserInfo":  {"user", "admin"},
		"/userinfo.UserInfo/ChangePassword":  {"user", "admin"},
		"/userinfo.UserInfo/GetUserInfoByID": {"admin"},
		"/userinfo.UserInfo/DeleteUser":      {"admin"},
		"/permissions.Permissions/IsAdmin":   {"admin"},
	}

	grpcApp := grpcapp.New(
		log, grpcPort,
		authService, permService, userInfoService,
		accessibleRoles, jwtManager,
	)

	return &App{
		GRPCApp: grpcApp,
	}
}
