package permissions

import (
	"context"
	"errors"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type serverAPI struct {
	ssov1.UnimplementedPermissionsServer
	log  *slog.Logger
	perm services.Permissions
}

func Register(gRPC *grpc.Server, log *slog.Logger, perm services.Permissions) {
	ssov1.RegisterPermissionsServer(gRPC, &serverAPI{
		log:  log,
		perm: perm,
	})
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	const op = "perm.server.IsAdmin"

	log := s.log.With(slog.String("op", op))

	isAdmin, err := s.perm.IsAdmin(ctx, req.UserId)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, status.Error(codes.InvalidArgument, repository.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to check if user is admin", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Info("successfully checked if user is admin")

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
