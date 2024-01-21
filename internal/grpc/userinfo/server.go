package userinfo

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc"
	"log/slog"
)

type serverAPI struct {
	ssov1.UnimplementedUserInfoServer
	log      *slog.Logger
	userInfo services.UserInfo
}

// Register registers the UserInfo gRPC service implementation with the provided gRPC server.
func Register(gRPC *grpc.Server, log *slog.Logger, userInfo services.UserInfo) {
	ssov1.RegisterUserInfoServer(gRPC, &serverAPI{
		log:      log,
		userInfo: userInfo,
	})
}
