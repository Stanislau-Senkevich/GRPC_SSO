package tests

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
)

func Test_IsAdmin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	ctx = st.SignInAndGetContext(admin, ctx, t)

	resp, err := st.PermissionsClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: 1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp)
}

func Test_IsAdmin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	table := []struct {
		name        string
		email       string
		password    string
		userID      int64
		expectedErr string
	}{
		{
			name:        "No rights to use",
			email:       "notadmin@gmail.com",
			password:    "123",
			userID:      1,
			expectedErr: grpcerror.ErrForbidden.Error(),
		},
		{
			name:        "User not found",
			email:       "admin@gmail.com",
			password:    "123",
			userID:      -1,
			expectedErr: grpcerror.ErrUserNotFound.Error(),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			user := models.User{
				Email:    tt.email,
				PassHash: tt.password,
			}

			ctx := st.SignInAndGetContext(user, ctx, t)

			resp, err := st.PermissionsClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
				UserId: tt.userID,
			})
			require.Error(t, err)
			require.Empty(t, resp)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestIsAdmin_BadToken(t *testing.T) {
	ctx, st := suite.New(t)

	table := []struct {
		name        string
		token       string
		userID      int64
		expectedErr string
	}{
		{
			name:        "Empty token",
			token:       "",
			userID:      1,
			expectedErr: grpcerror.ErrInvalidToken.Error(),
		},
		{
			name:        "Invalid token",
			token:       "Bearer adas",
			userID:      1,
			expectedErr: grpcerror.ErrInvalidToken.Error(),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.AppendToOutgoingContext(ctx, "authorization", tt.token)

			resp, err := st.PermissionsClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
				UserId: tt.userID,
			})

			require.Error(t, err)
			require.Empty(t, resp)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
