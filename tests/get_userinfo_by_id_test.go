package tests

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserInfoByID_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(admin, ctx, t)

	respGetInfo, err := st.UserInfoClient.GetUserInfoByID(ctx, &ssov1.GetUserInfoByIDRequest{
		UserId: user.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respGetInfo)

	assert.Equal(t, user.Email, respGetInfo.GetEmail())
	assert.Equal(t, user.Name, respGetInfo.GetName())
	assert.Equal(t, user.PhoneNumber, respGetInfo.GetPhoneNumber())
	assert.Equal(t, user.Surname, respGetInfo.GetSurname())
}

func TestGetUserInfoById_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	table := []struct {
		name        string
		userID      int64
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Forbidden",
			userID:      1,
			email:       "notadmin@gmail.com",
			password:    "123",
			expectedErr: grpcerror.ErrForbidden.Error(),
		},
		{
			name:        "User not found",
			userID:      -1,
			email:       "admin@gmail.com",
			password:    "123",
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

			respGetInfo, err := st.UserInfoClient.GetUserInfoByID(ctx, &ssov1.GetUserInfoByIDRequest{
				UserId: tt.userID,
			})
			require.Error(t, err)
			require.Empty(t, respGetInfo)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
