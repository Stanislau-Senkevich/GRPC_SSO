package userinfo

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

// DeleteUser deletes a user based on the provided gRPC request containing the user ID.
// It delegates the deletion operation to the DeleteUser method of the UserInfoService.
func (s *serverAPI) DeleteUser(
	ctx context.Context,
	req *ssov1.DeleteUserRequest) (
	*ssov1.DeleteUserResponse, error) {
	const op = "userinfo.grpc.DeleteUser"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("deleting user", slog.Int64("user_id", req.GetUserId()))

	err := s.userInfo.DeleteUser(ctx, req.GetUserId())
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to delete user", sl.Err(err),
			slog.Int64("user_id", req.GetUserId()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Info("user successfully deleted", slog.Int64("user_id", req.GetUserId()))

	return &ssov1.DeleteUserResponse{
		Succeed: true,
	}, nil
}
