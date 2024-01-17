package mongodb

import (
	"GRPC_SSO/internal/config"
	"GRPC_SSO/internal/domain/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (m *MongoRepository) Login(ctx context.Context, email, password string) (int64, error) {
	var user models.User

	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	filter := bson.D{
		{"email", email},
	}

	res := coll.FindOne(ctx, filter)

	err := res.Decode(&user)
	if err != nil {
		return -1, fmt.Errorf("failed to decode user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password+m.hashSalt))
	if err != nil {
		return -1, fmt.Errorf("failed to log in user: %w", err)
	}

	return user.Id, nil
}

func (m *MongoRepository) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	filter := bson.D{
		{"email", user.Email},
	}

	curEmail, err := coll.Find(ctx, filter)
	if err != nil {
		return -1, fmt.Errorf("error due searching in db: %w", err)
	}
	defer func() {
		_ = curEmail.Close(context.Background())
	}()

	if curEmail.RemainingBatchLength() > 0 {
		return -1, fmt.Errorf("email is already taken")
	}

	id, err := m.getNewUserId()
	if err != nil {
		return -1, err
	}

	user.Id = id

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		return -1, fmt.Errorf("error due inserting user: %w", err)
	}

	return id, nil
}

func (m *MongoRepository) getNewUserId() (int64, error) {
	coll := m.Db.Database(m.Config.DBName).Collection(
		m.Config.Collections[config.UserCollection])

	id, err := coll.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return -1, fmt.Errorf("error due generating new id: %w", err)
	}
	id++
	return id, nil
}
