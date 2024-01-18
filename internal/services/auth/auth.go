package auth

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository"
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

	log.Info("user successfully logged in")

	user := models.User{ID: userId, Email: email}

	token, err := jwt.NewToken(user, s.tokenTTL, s.signingKey)
	if err != nil {
		s.log.Error("failed to generate jwt-token", sl.Err(err))
		return "", err
	}

	log.Info("token successfully generated")

	return token, nil
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
