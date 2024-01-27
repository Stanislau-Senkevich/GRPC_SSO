package permissions

import (
	"context"
	"errors"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// IsAdmin checks if the user identified by the provided user ID has admin privileges.
// It delegates the admin check operation to the IsAdmin method of the PermissionsService.
func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	const op = "perm.grpc.IsAdmin"

	log := s.log.With(slog.String("op", op))

	log.Info("checking if user is admin", slog.Int64("user_id", req.UserId))

	isAdmin, err := s.perm.IsAdmin(ctx, req.UserId)
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to check if user is admin", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("successfully checked if user is admin")

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
