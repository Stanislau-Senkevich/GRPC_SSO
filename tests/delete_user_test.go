package tests

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteUser_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(admin, ctx, t)

	resp, err := st.UserInfoClient.DeleteUser(ctx, &ssov1.DeleteUserRequest{
		UserId: user.ID,
	})
	require.NoError(t, err)
	require.True(t, resp.GetSucceed())
}
