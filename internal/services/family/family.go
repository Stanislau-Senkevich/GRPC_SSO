package family

import (
	"context"
	"errors"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/client/family/grpc"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type FamilyService struct {
	client        *grpc.Client
	manager       *jwt.Manager
	adminEmail    string
	adminPassword string
}

func New(
	client *grpc.Client,
	manager *jwt.Manager,
	adminEmail string,
	adminPassword string,
) *FamilyService {
	return &FamilyService{
		client:        client,
		manager:       manager,
		adminEmail:    adminEmail,
		adminPassword: adminPassword,
	}
}

func (s *FamilyService) DeleteUserFromFamilies(ctx context.Context, userID int64, familyIDs []int64) error {
	const op = "family.service.DeleteUserFromFamilies"

	log := s.client.Log.With(
		slog.String("op", op),
	)

	if familyIDs == nil {
		return errors.New("nil familyIDs were provided") //nolint
	}

	for _, fID := range familyIDs {
		_, err := s.client.FamilyLeader.RemoveUser(ctx, &famv1.RemoveUserRequest{
			UserId:   userID,
			FamilyId: fID,
		})
		if errors.Is(err, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())) {
			log.Error("failed to delete user from family", sl.Err(err),
				slog.Int64("family_id", fID), slog.Int64("user_id", userID))
			return fmt.Errorf("%s: %w", op, err)
		}
		if err != nil {
			log.Warn("failed to delete user from family", sl.Err(err),
				slog.Int64("family_id", fID), slog.Int64("user_id", userID))
		}
	}

	return nil
}

func (s *FamilyService) DeleteUserInvites(ctx context.Context, userID int64) error {
	const op = "family.service.DeleteUserInvites"

	log := s.client.Log.With(
		slog.String("op", op),
	)

	_, err := s.client.Invite.DeleteUserInvites(ctx, &famv1.DeleteUserInvitesRequest{
		UserId: userID,
	})
	if err != nil {
		log.Error("failed to delete user's invites",
			sl.Err(err), slog.Int64("user_id", userID))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
