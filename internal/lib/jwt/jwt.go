package jwt

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	signingKey []byte
	tokenTTL   time.Duration
}

// New creates and returns a new instance of the JWTManager with the provided
// signing key and tokenTTL.
func New(signingKey []byte, tokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		signingKey: signingKey,
		tokenTTL:   tokenTTL,
	}
}

// NewToken generates a new JWT token for the provided user with the configured
// TTL and signing key. The token includes user-specific claims such as
// user ID, email, role, and expiration time.
func (m *JWTManager) NewToken(user models.User) (string, error) {
	claims := jwt.MapClaims{}

	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(m.tokenTTL).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	tokenString, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseToken parses the provided JWT token string and validates its signature
// using the configured signing key. It returns the claims embedded in the token
// if the signature is valid.
func (m *JWTManager) ParseToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(tkn *jwt.Token) (i interface{}, err error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"])
		}
		return m.signingKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to get claims from token")
	}

	return claims, nil
}

// GetClaims extracts and returns the JWT claims from the authorization token
// in the provided context. It relies on the ParseToken method to parse and
// validate the token's signature.
func (m *JWTManager) GetClaims(ctx context.Context) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get token claims")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, fmt.Errorf("authorization token was not provided")
	}

	accessToken := strings.Fields(values[0])[1]

	claims, err := m.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
