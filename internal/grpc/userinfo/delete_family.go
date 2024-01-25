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

func (s *serverAPI) DeleteFamily(
	ctx context.Context,
	req *ssov1.DeleteFamilyRequest) (
	*ssov1.DeleteFamilyResponse, error) {
	const op = "userinfo.grpc.AddFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.userInfo.DeleteFamily(ctx, req.GetFamilyId(), req.GetUserId())
	if errors.Is(err, grpcerror.ErrUserNotInFamily) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotInFamily.Error())
	}
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to delete family",
			sl.Err(err), slog.Int64("family_id", req.GetFamilyId()),
			slog.Int64("user_id", req.GetUserId()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.DeleteFamilyResponse{
		Succeed: true,
	}, nil
}
