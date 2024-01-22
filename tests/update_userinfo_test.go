package tests

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateUserInfo_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	user := st.SignUpRandomUser(ctx, t)

	updatedUser := suite.CreateRandomUser()

	ctx = st.SignInAndGetContext(user, ctx, t)

	respUpdate, err := st.UserInfoClient.UpdateUserInfo(ctx, &ssov1.UpdateUserInfoRequest{
		NewEmail:       updatedUser.Email,
		NewPhoneNumber: updatedUser.PhoneNumber,
		NewName:        updatedUser.Name,
		NewSurname:     updatedUser.Surname,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respUpdate.GetSucceed())
	assert.True(t, respUpdate.GetSucceed())
}
