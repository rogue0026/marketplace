package application

import (
	"errors"
	"os"
)

var (
	ErrEmptyDatabaseURL = errors.New("empty database url")
	ErrEmptyGRPCAddress = errors.New("empty grpc address")
)

type Config struct {
	DatabaseURL    string
	GRPCServerAddr string
}

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if len(dbURL) == 0 {
		return nil, ErrEmptyDatabaseURL
	}

	grpcAddr := os.Getenv("GRPC_SERVER_ADDRESS")
	if len(grpcAddr) == 0 {
		return nil, ErrEmptyGRPCAddress
	}

	cfg := &Config{
		DatabaseURL:    dbURL,
		GRPCServerAddr: grpcAddr,
	}

	return cfg, nil
}
