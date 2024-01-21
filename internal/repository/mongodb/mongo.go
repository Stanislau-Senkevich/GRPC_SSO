package mongodb

import (
	"context"
	"fmt"
	"github.com/Stanislau-Senkevich/GRPC_SSO/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type MongoRepository struct {
	Db     *mongo.Client
	Config *config.MongoConfig
	log    *slog.Logger
}

// InitMongoRepository initializes a new MongoRepository instance with the provided
// configuration, logger, and hash salt. It establishes a connection to the MongoDB
// server, performs a ping to ensure connectivity, and returns the initialized
// MongoRepository instance.
func InitMongoRepository(cfg *config.MongoConfig, logger *slog.Logger) (
	*MongoRepository, error) {
	const op = "mongo.InitMongoRepository"

	log := logger.With(
		slog.String("op", op),
	)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	conn := fmt.Sprintf(cfg.ConnectionString, cfg.User, cfg.Password)
	//conn := cfg.ConnectionString
	opts := options.Client().ApplyURI(conn).SetServerAPIOptions(serverAPI)

	log.Info("trying to connect to mongodb")

	db, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	log.Info("connected successfully")
	log.Info("trying to ping mongodb")

	if err = db.Database(cfg.DBName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}
	log.Info("pinged successfully")

	return &MongoRepository{
		Db:     db,
		Config: cfg,
		log:    logger,
	}, nil
}
