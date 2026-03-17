package pg

import (
	"context"
	"user_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepo struct {
	pool *pgxpool.Pool
}

func NewUsersRepo(pool *pgxpool.Pool) *UsersRepo {
	r := &UsersRepo{
		pool: pool,
	}

	return r
}

func (r *UsersRepo) CreateUser(ctx context.Context, user *models.User) (uint64, error) {
	var userId uint64
	sqlQuery := `
	INSERT INTO users (username, password_hash)
	VALUES ($1, $2) RETURNING id;
	`
	_, err := r.pool.Exec(ctx, sqlQuery, user.Username, user.PasswordHash)
	if err != nil {
		return userId, err
	}

	err = r.pool.QueryRow(ctx, sqlQuery, user.Username, user.PasswordHash).Scan(&userId)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func (r *UsersRepo) DeleteUser(ctx context.Context, userId uint64) error {
	sqlQuery := `
	DELETE FROM users
	WHERE id = $1;
	`
	_, err := r.pool.Exec(ctx, sqlQuery, userId)
	if err != nil {
		return err
	}

	return nil
}
