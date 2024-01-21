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

// ChangePassword changes the user's password based on the provided gRPC request.
// It delegates the password change operation to the ChangePassword method of the UserInfoService.
func (s *serverAPI) ChangePassword(
	ctx context.Context,
	req *ssov1.ChangePasswordRequest) (
	*ssov1.ChangePasswordResponse, error) {
	const op = "userinfo.grpc.ChangePassword"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("changing user's password")

	if req.GetNewPassword() == "" {
		log.Info("no password was provided")
		return nil, status.Error(codes.InvalidArgument, "new password is required")
	}

	if len(req.GetNewPassword()) >= 72 {
		log.Info("password is too long")
		return nil, status.Error(codes.InvalidArgument, "password is too long")
	}

	err := s.userInfo.ChangePassword(ctx, req.OldPassword, req.GetNewPassword())
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if errors.Is(err, grpcerror.ErrInvalidPassword) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrInvalidPassword.Error())
	}
	if err != nil {
		log.Error("failed to change user's password", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Info("password successfully changed")

	return &ssov1.ChangePasswordResponse{
		Succeed: true,
	}, nil
}
