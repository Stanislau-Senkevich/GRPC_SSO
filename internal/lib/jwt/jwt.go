package jwt

import (
	"fmt"
	"time"

	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"github.com/golang-jwt/jwt"
)

func NewToken(user models.User, duration time.Duration, signingKey []byte) (string, error) {
	claims := jwt.MapClaims{}

	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
