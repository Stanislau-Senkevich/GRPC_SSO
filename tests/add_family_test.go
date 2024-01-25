package tests

import (
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/tests/suite"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"testing"
)

func TestAddFamily_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(admin, ctx, t)

	famID := rand.Int63() + 1

	resp, err := st.UserInfoClient.AddFamily(ctx, &ssov1.AddFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.NoError(t, err)
	require.True(t, resp.GetSucceed())
}

func TestAddFamily_DuplicateAdd(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(admin, ctx, t)

	famID := rand.Int63() + 1

	resp, err := st.UserInfoClient.AddFamily(ctx, &ssov1.AddFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.NoError(t, err)
	require.True(t, resp.GetSucceed())

	resp, err = st.UserInfoClient.AddFamily(ctx, &ssov1.AddFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.ErrorContains(t, err, grpcerror.ErrUserInFamily.Error())
}
