package mongodb

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

func (m *MongoRepository) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	var user models.User
	const op = "mongo.permissions.IsAdmin"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("checking if user is admin", slog.Int64("user_id", userId))

	filter := bson.D{{"user_id", userId}}

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	res := coll.FindOne(ctx, filter)
	if res.Err() != nil {
		return false, repository.ErrUserNotFound
	}

	err := res.Decode(&user)
	if err != nil {
		log.Error("failed to decode user", sl.Err(err))
		return false, fmt.Errorf("failed to decode user: %w", err)
	}

	return user.IsAdmin, nil
}
