package userinfo

import (
	"context"
	"errors"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

// GetUserInfo retrieves the current user's information based on the provided gRPC request.
// It delegates the retrieval operation to the GetUserInfo method of the UserInfoService.
func (s *serverAPI) GetUserInfo(
	ctx context.Context,
	_ *ssov1.GetUserInfoRequest) (
	*ssov1.GetUserInfoResponse, error) {
	const op = "userinfo.grpc.GetUserInfo"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get user info")

	user, err := s.userInfo.GetUserInfo(ctx)
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to get user info", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user info successfully retrieved")

	return &ssov1.GetUserInfoResponse{
		UserId:       user.ID,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Name:         user.Name,
		Surname:      user.Surname,
		RegisteredAt: timestamppb.New(user.RegisteredAt.UTC()),
	}, nil
}
