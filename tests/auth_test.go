package tests

import (
	"fmt"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSignInSignUp_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	fmt.Println(st.Cfg)

	user := st.SignUpRandomUser(ctx, t)

	token := st.SignInAndGetToken(user, ctx, t)

	loginTime := time.Now()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return st.Cfg.SigningKey, nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, user.ID, int64(claims["user_id"].(float64)))
	assert.Equal(t, user.Email, claims["email"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestSignUp_DuplicatedSignUp(t *testing.T) {
	ctx, st := suite.New(t)

	user := st.SignUpRandomUser(ctx, t)

	respSignUp, err := st.AuthClient.SignUp(ctx, suite.SignUpRequestFromUser(user))
	require.Error(t, err)
	assert.Empty(t, respSignUp)
	assert.ErrorContains(t, err, grpcerror.ErrUserExists.Error())
}

func TestSignUp_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	table := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Sign up with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Sign up with empty email",
			email:       "",
			password:    suite.RandomFakePassword(),
			expectedErr: "email format is invalid",
		},
		{
			name:        "Sign up with both empty",
			email:       "",
			password:    "",
			expectedErr: "email format is invalid",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
