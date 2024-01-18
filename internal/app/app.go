package app

import (
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_SSO/internal/app/grpc"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository/mongodb"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/auth"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/permissions"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	cfg *config.Config,
	tokenTTL time.Duration,
) *App {
	repo, err := mongodb.InitMongoRepository(&cfg.Mongo, log, cfg.HashSalt)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository: %w", err))
	}

	authService := auth.New(log, repo, cfg.TokenTTL, cfg.HashSalt, cfg.SigningKey)
	permService := permissions.New(log, repo)

	grpcApp := grpcapp.New(log, grpcPort, authService, permService)

	return &App{
		GRPCApp: grpcApp,
	}
}
