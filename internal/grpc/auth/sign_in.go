package auth

import (
	"context"
	"errors"
	grpc_error "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/badoux/checkmail"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// SignIn authenticates a user based on the provided gRPC request.
// It delegates the user authentication operation to the SignIn method of the AuthService.
func (s *serverAPI) SignIn(
	ctx context.Context,
	req *ssov1.SignInRequest,
) (*ssov1.SignInResponse, error) {
	const op = "server.auth.SignIn"
	log := s.log.With(
		slog.String("op", op),
	)

	if err := validateLogin(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}

	token, err := s.auth.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, grpc_error.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, grpc_error.ErrUserNotFound.Error())
		}
		log.Error("failed to log in user", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.SignInResponse{
		Token: token,
	}, nil
}

func validateLogin(email, password string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return status.Error(codes.InvalidArgument, "email format is invalid")
	}

	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
