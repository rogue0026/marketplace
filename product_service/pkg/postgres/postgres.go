package postgres

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmptyConnString = errors.New("empty connection string")
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if len(url) == 0 {
		return nil, ErrEmptyConnString
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	
	return pool, nil
}
