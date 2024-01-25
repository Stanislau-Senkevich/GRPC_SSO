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

func TestDeleteFamily_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(admin, ctx, t)

	famID := rand.Int63() + 1

	respAdd, err := st.UserInfoClient.AddFamily(ctx, &ssov1.AddFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.NoError(t, err)
	require.True(t, respAdd.GetSucceed())

	respDelete, err := st.UserInfoClient.DeleteFamily(ctx, &ssov1.DeleteFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.NoError(t, err)
	require.True(t, respDelete.GetSucceed())
}

func TestDeleteFamily_DuplicateDelete(t *testing.T) {
	ctx, st := suite.New(t)

	admin := models.User{
		Email:    "admin@gmail.com",
		PassHash: "123",
	}

	user := st.SignUpRandomUser(ctx, t)

	ctx = st.SignInAndGetContext(admin, ctx, t)

	famID := rand.Int63() + 1

	respAdd, err := st.UserInfoClient.AddFamily(ctx, &ssov1.AddFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.NoError(t, err)
	require.True(t, respAdd.GetSucceed())

	respDelete, err := st.UserInfoClient.DeleteFamily(ctx, &ssov1.DeleteFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.NoError(t, err)
	require.True(t, respDelete.GetSucceed())

	respDelete, err = st.UserInfoClient.DeleteFamily(ctx, &ssov1.DeleteFamilyRequest{
		UserId:   user.ID,
		FamilyId: famID,
	})
	require.ErrorContains(t, err, grpcerror.ErrUserNotInFamily.Error())
}
