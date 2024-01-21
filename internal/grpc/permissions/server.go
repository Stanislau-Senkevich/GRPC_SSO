package permissions

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc"
	"log/slog"
)

type serverAPI struct {
	ssov1.UnimplementedPermissionsServer
	log  *slog.Logger
	perm services.Permissions
}

// Register sets up the gRPC server to handle Permissions service requests.
func Register(gRPC *grpc.Server, log *slog.Logger, perm services.Permissions) {
	ssov1.RegisterPermissionsServer(gRPC, &serverAPI{
		log:  log,
		perm: perm,
	})
}
