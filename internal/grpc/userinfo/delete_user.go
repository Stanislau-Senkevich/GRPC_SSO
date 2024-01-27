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

	log.Info("deleting user from families", slog.Int64("user_id", req.GetUserId()))

	user, err := s.userInfo.GetUserInfoByID(ctx, req.GetUserId())
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to get user's family list",
			sl.Err(err), slog.Int64("user_id", req.GetUserId()))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user info successfully retrieved",
		slog.Any("user", user))

	err = s.family.DeleteUserFromFamilies(ctx, req.GetUserId(), user.FamilyIDs)
	if err != nil {
		log.Error("failed to delete user from family", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user was deleted from families, trying to delete user's invites",
		slog.Int64("user_id", req.GetUserId()))

	err = s.family.DeleteUserInvites(ctx, req.GetUserId())
	if err != nil {
		log.Error("failed to delete user invites", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user's invites deleted, trying to delete user",
		slog.Int64("user_id", req.GetUserId()))

	err = s.userInfo.DeleteUser(ctx, req.GetUserId())
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to delete user", sl.Err(err),
			slog.Int64("user_id", req.GetUserId()))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user successfully deleted", slog.Int64("user_id", req.GetUserId()))

	return &ssov1.DeleteUserResponse{
		Succeed: true,
	}, nil
}
