package auth

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
	"time"
)

// SignUp registers a new user based on the provided gRPC request.
// It delegates the user registration operation to the SignUp method of the AuthService.
func (s *serverAPI) SignUp(
	ctx context.Context,
	req *ssov1.SignUpRequest,
) (*ssov1.SignUpResponse, error) {
	const op = "auth.grpc.SignUp"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to sign-up user")

	preUser := &models.User{
		ID:           -1,
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Name:         req.GetName(),
		Surname:      req.GetSurname(),
		PassHash:     req.GetPassword(),
		RegisteredAt: time.Now().UTC(),
	}

	if err := validateRegister(preUser); err != nil {
		log.Info("invalid input", sl.Err(err))
		return nil, err
	}

	userId, err := s.auth.SignUp(ctx, preUser)
	if err != nil {
		if errors.Is(err, grpcerror.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserExists.Error())
		}
		log.Error("failed to create user", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error: %s")
	}

	log.Info("user was created", slog.Int64("user_id", userId))

	return &ssov1.SignUpResponse{
		UserId: userId,
	}, nil
}

func validateRegister(user *models.User) error {
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return status.Error(codes.InvalidArgument, "email format is invalid")
	}

	if user.PassHash == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if user.PhoneNumber == "" {
		return status.Error(codes.InvalidArgument, "phone number is required")
	}

	if user.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	if user.Surname == "" {
		return status.Error(codes.InvalidArgument, "surname is required")
	}

	return nil
}
