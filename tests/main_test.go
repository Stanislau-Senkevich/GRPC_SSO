package tests

import (
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/subosito/gotenv"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := gotenv.Load("../.env")

	if err = config.BindEnv(); err != nil {
		panic(fmt.Errorf("failed to bind env: %w", err))
	}

	code := m.Run()

	os.Exit(code)
}
