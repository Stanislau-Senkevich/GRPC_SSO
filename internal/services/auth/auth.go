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
)

type AuthService struct {
	log      *slog.Logger
	repo     repository.AuthRepository
	hashSalt string
	manager  *jwt.JWTManager
}

// New creates and returns a new instance of the AuthService
func New(
	log *slog.Logger,
	repo repository.AuthRepository,
	manager *jwt.JWTManager,
	hashSalt string,
) *AuthService {
	return &AuthService{
		log:      log,
		repo:     repo,
		manager:  manager,
		hashSalt: hashSalt,
	}
}

// SignIn authenticates a user with the provided email and password by first validating
// the credentials against the authentication repository. If successful, it generates
// a JWT token for the user and returns it.
func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	const op = "auth.SignIn"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to log in user")

	passSalted := password + s.hashSalt
	user, err := s.repo.Login(ctx, email, passSalted)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully logged in")

	token, err := s.manager.NewToken(user)
	if err != nil {
		s.log.Error("failed to generate jwt-token", sl.Err(err))
		return "", err
	}

	log.Info("token successfully generated")

	return token, nil
}

// SignUp registers a new user by first generating a password hash, and then creating
// a new user entry in the authentication repository. It returns the assigned user ID
// upon successful registration.
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
