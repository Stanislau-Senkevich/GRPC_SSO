package auth

import (
	"GRPC_SSO/internal/domain/models"
	"GRPC_SSO/internal/lib/sl"
	"GRPC_SSO/internal/repository"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthService struct {
	log        *slog.Logger
	repo       repository.AuthRepository
	tokenTTL   time.Duration
	hashSalt   string
	signingKey []byte
}

func New(
	log *slog.Logger,
	repo repository.AuthRepository,
	tokenTTL time.Duration,
	hashSalt string,
	signingKey []byte) *AuthService {
	return &AuthService{
		log:        log,
		repo:       repo,
		tokenTTL:   tokenTTL,
		hashSalt:   hashSalt,
		signingKey: signingKey,
	}
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	const op = "auth.SignIn"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to log in user")

	userId, err := s.repo.Login(ctx, email, password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return fmt.Sprintf("%d", userId), nil
}

func (s *AuthService) SignUp(ctx context.Context, user *models.User) (int64, error) {
	const op = "auth.SignUp"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.PassHash+s.hashSalt), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	user.PassHash = string(passHash)

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}
