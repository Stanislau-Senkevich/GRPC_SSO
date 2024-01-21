package permissions

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository"
	"log/slog"
)

type PermService struct {
	log  *slog.Logger
	repo repository.PermissionsRepository
}

// New creates and returns a new instance of the PermService with the provided
// dependencies and configurations.
func New(log *slog.Logger, repo repository.PermissionsRepository) *PermService {
	return &PermService{log: log, repo: repo}
}

func (s *PermService) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	return s.repo.IsAdmin(ctx, userId)
}
