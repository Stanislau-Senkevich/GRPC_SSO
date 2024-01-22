package tests

import (
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := config.BindEnv("../.env"); err != nil {
		panic(fmt.Errorf("failed to bind env: %w", err))
	}

	code := m.Run()

	os.Exit(code)
}
