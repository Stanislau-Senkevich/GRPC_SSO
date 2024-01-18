package app

import (
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_SSO/internal/app/grpc"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository/mongodb"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services/auth"
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
		panic(fmt.Errorf("error due initializing repository: %w", err))
	}

	authService := auth.New(log, repo, cfg.TokenTTL, cfg.HashSalt, cfg.SigningKey)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCApp: grpcApp,
	}
}
