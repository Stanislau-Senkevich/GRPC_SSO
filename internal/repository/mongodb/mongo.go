package mongodb

import (
	"GRPC_SSO/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type MongoRepository struct {
	Db       *mongo.Client
	Config   *config.MongoConfig
	log      *slog.Logger
	hashSalt string
}

func InitMongoRepository(cfg *config.MongoConfig, log *slog.Logger, hashSalt string) (
	*MongoRepository, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	//conn := fmt.Sprintf(cfg.ConnectionString, cfg.User, cfg.Password)
	conn := cfg.ConnectionString
	opts := options.Client().ApplyURI(conn).SetServerAPIOptions(serverAPI)

	log.Info("trying to connect to mongodb")

	db, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("error due connecting to mongo: %w", err)
	}

	log.Info("connected successfully")
	log.Info("trying to ping mongodb")

	if err = db.Database(cfg.DBName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, fmt.Errorf("error due pinging mongodb: %w", err)
	}
	log.Info("pinged successfully")

	return &MongoRepository{
		Db:       db,
		Config:   cfg,
		log:      log,
		hashSalt: hashSalt,
	}, nil
}
