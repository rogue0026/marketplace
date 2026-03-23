package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CreateWalletQuery = `
	INSERT INTO wallets (user_id) 
	VALUES ($1)`

	AddMoneyQuery = `
	UPDATE wallets 
	SET balance = balance + $2
	WHERE user_id = $1
	`

	WriteOffMoneyQuery = `
	UPDATE wallets
	SET balance = balance - $2
	WHERE user_id = $1
	`
)

type WalletsRepo struct {
	pool *pgxpool.Pool
}

func NewWalletsRepo(pool *pgxpool.Pool) *WalletsRepo {
	r := &WalletsRepo{
		pool: pool,
	}

	return r
}

func (r *WalletsRepo) CreateWallet(ctx context.Context, userId uint64) error {
	_, err := r.pool.Exec(ctx, CreateWalletQuery, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletsRepo) AddMoney(ctx context.Context, userId uint64, amount uint64) error {
	_, err := r.pool.Exec(ctx, AddMoneyQuery, userId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletsRepo) WriteOffMoney(ctx context.Context, userId uint64, amount uint64) error {
	_, err := r.pool.Exec(ctx, WriteOffMoneyQuery, userId, amount)
	if err != nil {
		return err
	}

	return nil
}
