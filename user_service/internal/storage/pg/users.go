package pg

import (
	"context"
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
