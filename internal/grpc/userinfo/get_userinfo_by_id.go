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

// GetUserInfoByID retrieves user information based on the provided gRPC request
// containing the user ID. It delegates the retrieval operation to the GetUserInfoByID
// method of the UserInfoService.
func (s *serverAPI) GetUserInfoByID(
	ctx context.Context,
	req *ssov1.GetUserInfoByIDRequest) (
	*ssov1.GetUserInfoByIDResponse, error) {
	const op = "userinfo.grpc.GetUserInfoByID"

	log := s.log.With(
		slog.String("op", op),
	)

	user, err := s.userInfo.GetUserInfoByID(ctx, req.GetUserId())
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to get user info", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Info("user info successfully retrieved")
	return &ssov1.GetUserInfoByIDResponse{
		UserId:       user.ID,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Name:         user.Name,
		Surname:      user.Surname,
		RegisteredAt: timestamppb.New(user.RegisteredAt.UTC()),
	}, nil
}
