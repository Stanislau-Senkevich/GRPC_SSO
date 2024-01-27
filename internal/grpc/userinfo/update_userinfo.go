package userinfo

import (
	"context"
	"errors"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/badoux/checkmail"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// UpdateUserInfo updates user information based on the provided gRPC request.
// It validates the new email format and delegates
// the update operation to the UpdateUserInfo method of the UserInfoService.
func (s *serverAPI) UpdateUserInfo(
	ctx context.Context,
	req *ssov1.UpdateUserInfoRequest) (
	*ssov1.UpdateUserInfoResponse, error) {

	const op = "userinfo.grpc.UpdateUserInfo"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("updating user info")

	if err := checkmail.ValidateFormat(req.GetNewEmail()); req.GetNewEmail() != "" && err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email was provided")
	}

	updateInfo := &models.User{
		Email:       req.GetNewEmail(),
		PhoneNumber: req.GetNewPhoneNumber(),
		Name:        req.GetNewName(),
		Surname:     req.GetNewSurname(),
	}

	err := s.userInfo.UpdateUserInfo(ctx, updateInfo)
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Info(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to update user", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("info successfully updated")

	return &ssov1.UpdateUserInfoResponse{
		Succeed: true,
	}, nil
}
