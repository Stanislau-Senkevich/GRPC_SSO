package grpcapp

import (
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/grpc/auth"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/grpc/permissions"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/grpc/userinfo"
	jwtmanager "github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	gRPCConfig *config.GRPCConfig
}

// New creates a new instance of the application with the specified dependencies and configurations.
func New(
	log *slog.Logger,
	gRPCConfig *config.GRPCConfig,

	authService services.Auth,
	permService services.Permissions,
	userInfoService services.UserInfo,
	accessibleRoles map[string][]string,
	jwtManager *jwtmanager.Manager,
) *App {
	interceptor := NewJWTInterceptor(jwtManager, accessibleRoles)

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
		grpc.ConnectionTimeout(gRPCConfig.Timeout),
	)

	auth.Register(gRPCServer, log, authService)
	permissions.Register(gRPCServer, log, permService)
	userinfo.Register(gRPCServer, log, userInfoService)

	return &App{log, gRPCServer, gRPCConfig}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run starts the gRPC server and listens for incoming requests on the specified port.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.gRPCConfig.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop gracefully stops the running gRPC server, allowing it to finish processing existing requests.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping grpc server", slog.Int("port", a.gRPCConfig.Port))

	a.gRPCServer.GracefulStop()
}
