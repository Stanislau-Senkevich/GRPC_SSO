package app

import (
	"context"
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_SSO/internal/app/grpc"
	grpcclient "github.com/Stanislau-Senkevich/GRPC_SSO/internal/client/family/grpc"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	jwtmanager "github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository/mongodb"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/auth"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/family"
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
	cfg *config.Config,
	tokenTTL time.Duration,
) *App {
	repo, err := mongodb.InitMongoRepository(&cfg.Mongo, log)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository: %w", err))
	}

	jwtManager := jwtmanager.New(cfg.SigningKey, tokenTTL)
	log.Info("jwt-manager initialized")

	familyClient, err := grpcclient.New(
		context.Background(), log,
		cfg.ClientsConfig.Family.Address,
		cfg.ClientsConfig.Family.Timeout,
		cfg.ClientsConfig.Family.RetriesCount)
	if err != nil {
		panic(fmt.Errorf("failed to initialize client SSO: %w", err))
	}
	log.Info("family client initialized")

	authService := auth.New(log, repo, jwtManager, cfg.HashSalt)
	log.Info("auth service initialized")

	permService := permissions.New(log, repo)
	log.Info("permissions service initialized")

	userInfoService := userinfo.New(log, repo, jwtManager, cfg.HashSalt)
	log.Info("userinfo service initialized")

	familyService := family.New(
		familyClient,
		jwtManager,
		cfg.ClientsConfig.AdminEmail,
		cfg.ClientsConfig.AdminPassword)
	log.Info("family service initialized")

	accessibleRoles := map[string][]string{
		"/permissions.Permissions/IsAdmin":   {"admin"},
		"/userinfo.UserInfo/GetUserInfo":     {"user", "admin"},
		"/userinfo.UserInfo/UpdateUserInfo":  {"user", "admin"},
		"/userinfo.UserInfo/ChangePassword":  {"user", "admin"},
		"/userinfo.UserInfo/GetUserInfoByID": {"admin"},
		"/userinfo.UserInfo/AddFamily":       {"admin"},
		"/userinfo.UserInfo/DeleteFamily":    {"admin"},
		"/userinfo.UserInfo/DeleteUser":      {"admin"},
	}

	grpcApp := grpcapp.New(
		log, &cfg.GRPC,
		authService, permService,
		userInfoService, familyService,
		accessibleRoles, jwtManager,
	)

	return &App{
		GRPCApp: grpcApp,
	}
}
