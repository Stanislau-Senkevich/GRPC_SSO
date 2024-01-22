package suite

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	ssov1 "github.com/Stanislau-Senkevich/protocols/gen/go/sso"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net"
	"strconv"
	"testing"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg               *config.Config
	AuthClient        ssov1.AuthClient
	PermissionsClient ssov1.PermissionsClient
	UserInfoClient    ssov1.UserInfoClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_tests.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(&cfg.GRPC),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:                 t,
		Cfg:               cfg,
		AuthClient:        ssov1.NewAuthClient(cc),
		PermissionsClient: ssov1.NewPermissionsClient(cc),
		UserInfoClient:    ssov1.NewUserInfoClient(cc),
	}
}

func (s *Suite) SignUpRandomUser(ctx context.Context, t *testing.T) models.User {
	user := CreateRandomUser()

	respSignUp, err := s.AuthClient.SignUp(ctx, SignUpRequestFromUser(user))
	require.NoError(t, err)
	require.NotEmpty(t, respSignUp.GetUserId())
	user.ID = respSignUp.GetUserId()

	return user
}

func (s *Suite) SignInAndGetToken(user models.User, ctx context.Context, t *testing.T) string {
	respSignIn, err := s.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Email:    user.Email,
		Password: user.PassHash,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respSignIn.GetToken())

	return respSignIn.GetToken()
}

func (s *Suite) SignInAndGetContext(user models.User, ctx context.Context, t *testing.T) context.Context {
	token := s.SignInAndGetToken(user, ctx, t)

	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func grpcAddress(cfg *config.GRPCConfig) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.Port))
}

func RandomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, rand.Intn(20)+1)
}

func CreateRandomUser() models.User {
	return models.User{
		Email:       gofakeit.Email(),
		PhoneNumber: gofakeit.Phone(),
		Name:        gofakeit.Name(),
		Surname:     gofakeit.Name(),
		PassHash:    RandomFakePassword(),
	}
}

func SignUpRequestFromUser(user models.User) *ssov1.SignUpRequest {
	return &ssov1.SignUpRequest{
		Email:       user.Email,
		Password:    user.PassHash,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		Surname:     user.Surname,
	}
}
