package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	UserCollection     = "user"
	SequenceCollection = "sequence"
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

// MustLoad reads and loads the configuration from a file specified by the fetched config path.
// It panics if any error occurs during the process, ensuring that the application cannot proceed
// without a valid configuration.
func MustLoad() *Config { //nolint

	path := fetchConfigPath()

	if path == "" {
		panic(fmt.Errorf("config path is empty"))
	}

	if err := BindEnv(); err != nil {
		panic(fmt.Errorf("failed to bind the enviroment: %w", err))
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	var cfg Config

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Errorf("config file does not exist: %s", path))
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	if err := parseEnv(&cfg); err != nil {
		panic(fmt.Errorf("failed to parse the enviroment: %w", err))
	}

	return &cfg
}

// parseEnv sets up configuration parameters by binding them to environment variables using the
// Viper library. It ensures that necessary environment variables are available and assigns their
// values to corresponding fields in the provided Config struct. If any binding operation fails,
// it returns an error indicating the specific failure.
func parseEnv(cfg *Config) error {

	cfg.Mongo.User = viper.GetString("mongo_user")
	cfg.Mongo.Password = viper.GetString("mongo_password")
	cfg.HashSalt = viper.GetString("hash_salt")
	cfg.SigningKey = []byte(viper.GetString("signing_key"))

	return nil
}

func BindEnv() error {
	//if err := gotenv.Load(); err != nil {
	//	return fmt.Errorf("failed to parse .env file: %w", err)
	//}

	if err := viper.BindEnv("mongo_user"); err != nil {
		return fmt.Errorf("failed to set up mongo_user: %w", err)
	}

	if err := viper.BindEnv("mongo_password"); err != nil {
		return fmt.Errorf("failed to set up mongo_password: %w", err)
	}

	if err := viper.BindEnv("hash_salt"); err != nil {
		return fmt.Errorf("failed to set up hash_salt: %w", err)
	}

	if err := viper.BindEnv("signing_key"); err != nil {
		return fmt.Errorf("failed to set up signing_key: %w", err)
	}

	return nil
}

func fetchConfigPath() string {
	var res string

	//--config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
