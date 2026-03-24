package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrEmptyDatabaseURL = errors.New("empty database url")
	ErrEmptyGRPCAddr    = errors.New("empty grpc address")
)

type AppConfig struct {
	DatabaseURL  string
	GRPCAddr     string
	UsersAddr    string
	ProductsAddr string
}

func Load() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if len(databaseURL) == 0 {
		return nil, ErrEmptyDatabaseURL
	}

	grpcAddr := os.Getenv("GRPC_SERVER_ADDR")
	if len(grpcAddr) == 0 {
		return nil, ErrEmptyGRPCAddr
	}

	cfg := &AppConfig{
		DatabaseURL: databaseURL,
		GRPCAddr:    grpcAddr,
	}

	return cfg, nil
}
