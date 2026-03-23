package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

const ()

type ReservationsRepo struct {
	pool *pgxpool.Pool
}

func NewReservationsRepo(pool *pgxpool.Pool) *ReservationsRepo {
	reservations := &ReservationsRepo{
		pool: pool,
	}

	return reservations
}
