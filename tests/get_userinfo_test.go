package tests

import (
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestGetUserInfo_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(user, ctx, t)

	respGetInfo, err := st.UserInfoClient.GetUserInfo(ctx, &ssov1.GetUserInfoRequest{})
	require.NoError(t, err)
	assert.Equal(t, user.Email, respGetInfo.GetEmail())
	assert.Equal(t, user.Name, respGetInfo.GetName())
	assert.Equal(t, user.PhoneNumber, respGetInfo.GetPhoneNumber())
	assert.Equal(t, user.Surname, respGetInfo.GetSurname())
}

func TestGetUserInfo_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	table := []struct {
		name        string
		token       string
		errExpected string
	}{
		{
			name:        "Empty token",
			token:       "",
			errExpected: grpcerror.ErrInvalidToken.Error(),
		},
		{
			name:        "Invalid token",
			token:       "sgfdss",
			errExpected: grpcerror.ErrInvalidToken.Error(),
		},
		{
			name:        "Expired token",
			token:       "Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1hY3lrdWhsbWFuQGJlcmdzdHJvbS5uZXQiLCJleHAiOjE3MDU4NTI5MDIsInJvbGUiOiJ1c2VyIiwidXNlcl9pZCI6NX0.yp81A5-wnVkdch_oZuMX3WRzBbyvnqfuEW_GBGAbEmiypiCVrO3Pkl97H_u3Awdd",
			errExpected: grpcerror.ErrInvalidToken.Error(),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.AppendToOutgoingContext(ctx, "authorization", tt.token)
			resp, err := st.UserInfoClient.GetUserInfo(ctx, &ssov1.GetUserInfoRequest{})
			require.Error(t, err)
			require.Empty(t, resp)
			assert.ErrorContains(t, err, tt.errExpected)
		})
	}
}
