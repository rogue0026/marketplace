package pg

import (
	"context"
	"fmt"
	"user_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CreateUserQuery = `
	INSERT INTO users (login, password_hash)
	VALUES ($1, $2)
	RETURNING id
	`

	DeleteUserQuery = `
	DELETE FROM users
	WHERE id = $1
	`
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

	err := r.pool.QueryRow(
		ctx,
		CreateUserQuery,
		user.Username,
		user.PasswordHash,
	).Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("failed to create user with login=%s: %w", user.Username, err)
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
		return fmt.Errorf("failed to delete user with user_id=%d: %w", userId, err)
	}

	return nil
}
