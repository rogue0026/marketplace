package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *WalletsRepo) ToUpBalance(ctx context.Context, userId uint64, amount uint64) error {
	sqlQuery := `
	UPDATE wallets
	SET balance = balance + $1
	WHERE user_id = $2;
	`
	_, err := r.pool.Exec(ctx, sqlQuery, amount, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *WalletsRepo) ToDownBalance(ctx context.Context, userId uint64, amount uint64) error {
	sqlQuery := `
	UPDATE wallets
	SET balance = balance - $1
	WHERE user_id = $2;
	`

	_, err := r.pool.Exec(ctx, sqlQuery, amount, userId)
	if err != nil {
		return err
	}

	return nil
}
