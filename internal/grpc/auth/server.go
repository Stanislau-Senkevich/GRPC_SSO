package auth

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc"
	"log/slog"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	log  *slog.Logger
	auth services.Auth
}

// Register associates the gRPC implementation of the Auth service with the provided gRPC server.
func Register(gRPC *grpc.Server, log *slog.Logger, auth services.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		log:  log,
		auth: auth,
	})
}
