package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"os"
	"time"
)

const (
	UserCollection = "user"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	TokenTTL   time.Duration `yaml:"token_ttl"`
	Mongo      MongoConfig   `yaml:"mongo_config"`
	GRPC       GRPCConfig    `yaml:"grpc"`
	HashSalt   string
	SigningKey []byte
}

type MongoConfig struct {
	User             string
	Password         string
	DBName           string            `yaml:"db_name"`
	ConnectionString string            `yaml:"conn_string"`
	Collections      map[string]string `yaml:"collections"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config { //nolint
	var (
		cfg  Config
		path = fetchConfigPath()
	)

	if path == "" {
		panic(fmt.Errorf("config path is empty"))
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Errorf("config file does not exist: %s", path))
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(fmt.Errorf("error due reading config: %w", err))
	}

	if err := parseEnv(&cfg); err != nil {
		panic(fmt.Errorf("error due parsing env: %w", err))
	}

	return &cfg
}

func parseEnv(cfg *Config) error {
	if err := gotenv.Load(); err != nil {
		return fmt.Errorf("env was not uploaded from file: %w", err)
	}

	if err := viper.BindEnv("mongo_user"); err != nil {
		return fmt.Errorf("mongo_user was not set up: %w", err)
	}

	if err := viper.BindEnv("mongo_password"); err != nil {
		return fmt.Errorf("mongo_password was not set up: %w", err)
	}

	if err := viper.BindEnv("hash_salt"); err != nil {
		return fmt.Errorf("mongo_user was not set up: %w", err)
	}

	if err := viper.BindEnv("signing_key"); err != nil {
		return fmt.Errorf("mongo_password was not set up: %w", err)
	}

	cfg.Mongo.User = viper.GetString("mongo_user")
	cfg.Mongo.Password = viper.GetString("mongo_password")
	cfg.HashSalt = viper.GetString("hash_salt")
	cfg.SigningKey = []byte(viper.GetString("signing_key"))

	return nil
}

func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
