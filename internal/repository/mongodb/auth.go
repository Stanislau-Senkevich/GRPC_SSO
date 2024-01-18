package mongodb

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/domain/models"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/lib/sl"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (m *MongoRepository) Login(ctx context.Context, email, password string) (int64, error) {
	const op = "mongo.auth.Login"

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
		return -1, repository.ErrUserNotFound
	}

	err := res.Decode(&user)
	if err != nil {
		log.Error("failed to decode user", sl.Err(err))
		return -1, repository.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password+m.hashSalt))
	if err != nil {
		return -1, repository.ErrUserNotFound
	}

	return user.ID, nil
}

func (m *MongoRepository) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	const op = "mongo.auth.CreateUser"
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
		return -1, repository.ErrUserExists
	}

	id, err := m.getNewUserId()
	if err != nil {
		return -1, err
	}

	user.ID = id

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		log.Error("failed to insert user", sl.Err(err))
		return -1, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

func (m *MongoRepository) getNewUserId() (int64, error) {
	const op = "mongo.auth.getNewUserId"

	log := m.log.With(
		slog.String("op", op),
	)

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	id, err := coll.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Error("failed to generate new id", sl.Err(err))
		return -1, fmt.Errorf("failed to generate new id: %w", err)
	}
	id++
	return id, nil
}
