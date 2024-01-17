package auth

import (
	"GRPC_SSO/internal/domain/models"
	"GRPC_SSO/internal/lib/sl"
	"GRPC_SSO/internal/services"
	"context"
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

	if err := validateLogin(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}

	token, err := s.auth.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		s.log.Info("failed to log in user", sl.Err(err))
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
	if err := validateRegister(req.GetEmail(), req.GetPassword(),
		req.GetPhoneNumber(), req.GetName(), req.GetSurname()); err != nil {
		return nil, err
	}

	preUser := &models.User{
		Id:           -1,
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Name:         req.GetName(),
		Surname:      req.GetSurname(),
		PassHash:     req.GetPassword(),
		RegisteredAt: time.Now().UTC(),
	}

	userId, err := s.auth.SignUp(ctx, preUser)
	if err != nil {
		s.log.Error("error due creating user", sl.Err(err))
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

func validateRegister(email, password, phoneNumber, name, surname string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return status.Error(codes.InvalidArgument, "email format is invalid")
	}

	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if phoneNumber == "" {
		return status.Error(codes.InvalidArgument, "phone number is required")
	}

	if name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	if surname == "" {
		return status.Error(codes.InvalidArgument, "surname is required")
	}

	return nil
}
