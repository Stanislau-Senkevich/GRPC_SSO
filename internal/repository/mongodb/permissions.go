package mongodb

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

// IsAdmin checks if the user with the provided user ID has admin privileges.
// It queries the MongoDB database to retrieve the user information and determines
// if the user has the admin role.
func (m *MongoRepository) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	var user models.User
	const op = "permissions.mongo.IsAdmin"

	log := m.log.With(
		slog.String("op", op),
	)

	filter := bson.D{{"user_id", userId}}

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	res := coll.FindOne(ctx, filter)
	if res.Err() != nil {
		return false, grpcerror.ErrUserNotFound
	}

	if err := res.Decode(&user); err != nil {
		log.Error("failed to decode user", sl.Err(err))
		return false, fmt.Errorf("failed to decode user: %w", err)
	}

	return string(user.Role) == string(models.AdminRole), nil
}
