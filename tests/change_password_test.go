package tests

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangePassword_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	user := st.SignUpRandomUser(ctx, t)

	newPass := suite.RandomFakePassword()

	ctx = st.SignInAndGetContext(user, ctx, t)

	respChange, err := st.UserInfoClient.ChangePassword(ctx, &ssov1.ChangePasswordRequest{
		OldPassword: user.PassHash,
		NewPassword: newPass,
	})
	require.NoError(t, err)
	require.True(t, respChange.GetSucceed())

	respSignIn, err := st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Email:    user.Email,
		Password: newPass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respSignIn.GetToken())
}
