package mongodb

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_SSO/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

// Login authenticates a user by verifying the provided email and password against
// the stored user data in MongoDB. If the authentication is successful, it returns
// the authenticated user; otherwise, it returns an error indicating the failure.
func (m *MongoRepository) Login(ctx context.Context, email, passwordSalted string) (models.User, error) {
	const op = "auth.mongo.Login"

	var user models.User

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	filter := bson.D{
		{"email", email},
	}

	res := coll.FindOne(ctx, filter)
	if res.Err() != nil {
		return models.User{}, grpcerror.ErrUserNotFound
	}

	if err := res.Decode(&user); err != nil {
		log.Error("failed to decode user", sl.Err(err))
		return models.User{}, fmt.Errorf("failed to decode user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(passwordSalted)); err != nil {
		return models.User{}, grpcerror.ErrUserNotFound
	}

	return user, nil
}

// CreateUser creates a new user in the MongoDB database. It first checks if a user
// with the provided email already exists, and if not, it assigns a new user ID,
// sets the user's role, and inserts the user into the database.
func (m *MongoRepository) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	const op = "auth.mongo.CreateUser"
	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	filter := bson.D{
		{"email", user.Email},
	}

	curEmail, err := coll.Find(ctx, filter)
	if err != nil {
		log.Error("failed to search in db", sl.Err(err))
		return -1, fmt.Errorf("failed to search in db: %w", err)
	}
	defer func() {
		_ = curEmail.Close(context.Background())
	}()

	if curEmail.RemainingBatchLength() > 0 {
		return -1, grpcerror.ErrUserExists
	}

	id, err := m.getNewUserId()
	if err != nil {
		return -1, err
	}

	user.ID = id
	user.Role = models.UserRole

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		log.Error("failed to insert user", sl.Err(err))
		return -1, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

// getNewUserId generates a new unique user ID
func (m *MongoRepository) getNewUserId() (int64, error) {
	var seq models.Sequence

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.SequenceCollection])

	filter := bson.D{
		{"collection_name", config.UserCollection},
	}

	update := bson.D{
		{"$inc", bson.D{
			{"counter", 1},
		},
		},
	}

	res := coll.FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() != nil {
		return -1, fmt.Errorf("failed to get id: %w", res.Err())
	}

	err := res.Decode(&seq)
	if err != nil {
		return -1, fmt.Errorf("failed to decode sequence: %w", err)
	}

	return seq.Counter, nil
}
