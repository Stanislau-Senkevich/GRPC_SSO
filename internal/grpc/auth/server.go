package auth

import (
	"context"
	"errors"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/services"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/badoux/checkmail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	log  *slog.Logger
	auth services.Auth
}

func Register(gRPC *grpc.Server, log *slog.Logger, auth services.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		log:  log,
		auth: auth,
	})
}

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
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, repository.ErrUserNotFound.Error())
		}
		log.Error("failed to log in user", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.SignInResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) SignUp(
	ctx context.Context,
	req *ssov1.SignUpRequest,
) (*ssov1.SignUpResponse, error) {
	const op = "server.auth.SignIn"
	log := s.log.With(
		slog.String("op", op),
	)

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
		return nil, err
	}

	userId, err := s.auth.SignUp(ctx, preUser)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, repository.ErrUserExists.Error())
		}
		log.Error("failed to create user", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal error: %s")
	}

	s.log.Info("user was created", slog.Int64("user_id", userId))

	return &ssov1.SignUpResponse{
		UserId: userId,
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
