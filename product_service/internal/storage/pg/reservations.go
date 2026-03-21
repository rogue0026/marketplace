package pg

import (
	"context"
	"product_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ADD_RESERVATIONS_QUERY = `
	INSERT INTO product_reservations (order_id, product_id, quantity, status, expires_at) 
	VALUES ($1, $2, $3, $4, $5);`

	REDUCE_PRODUCTS_QUANTITY_QUERY = `
	UPDATE products
	SET total_quantity = total_quantity - $2
	WHERE id = $1;
	`

	FIND_EXPIRED_RESERVATIONS_QUERY = `
	SELECT 
    	id, 
    	order_id,
    	product_id,
    	quantity
	FROM product_reservations WHERE expires_at < current_timestamp AND status = $1`

	CANCEL_RESERVATION_QUERY = `
	UPDATE products
	SET total_quantity = total_quantity + $2
	WHERE id = $1
	`

	DELETE_EXPIRED_RESERVATIONS_QUERY = `
	DELETE FROM product_reservations
    WHERE id = $1`
)

type ReservationsRepo struct {
	pool *pgxpool.Pool
}

func NewReservationsRepo(pool *pgxpool.Pool) *ReservationsRepo {
	reservations := &ReservationsRepo{
		pool: pool,
	}

	return reservations
}

func (r *ReservationsRepo) ReserveProducts(ctx context.Context, reservations []*models.Reservation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// добавляем записи в таблицу резервированияи и уменьшаем количество доступного товара на складе
	for _, elem := range reservations {
		// todo сделать обработку ситуации когда на складе нет нужного количества товаров
		_, err = tx.Exec(
			ctx,
			REDUCE_PRODUCTS_QUANTITY_QUERY,
			elem.ProductId,
			elem.Quantity,
		)
		if err != nil {
			return err
		}

		_, err := tx.Exec(
			ctx,
			ADD_RESERVATIONS_QUERY,
			elem.OrderId,
			elem.ProductId,
			elem.Quantity,
			elem.Status,
			elem.ExpiresAt,
		)
		if err != nil {
			return err
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReservationsRepo) CleanExpiredReservations(ctx context.Context, status int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// ищем резервирования, которые не были вовремя оплачены
	rows, err := tx.Query(ctx, FIND_EXPIRED_RESERVATIONS_QUERY, status)
	if err != nil {
		return err
	}
	defer rows.Close()

	expiredReservations := make([]*models.Reservation, 0)
	for rows.Next() {
		er := &models.Reservation{}
		if err := rows.Scan(er.Id, er.OrderId, er.ProductId, er.Quantity); err != nil {
			return err
		}
		expiredReservations = append(expiredReservations, er)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// возвращаем неоплаченные товары обратно на склад
	for _, elem := range expiredReservations {
		_, err := tx.Exec(ctx, CANCEL_RESERVATION_QUERY, elem.ProductId, elem.Quantity)
		if err != nil {
			return err
		}
	}

	// удаляем записи из таблицы резервирования товаров
	for _, elem := range expiredReservations {
		_, err := tx.Exec(ctx, DELETE_EXPIRED_RESERVATIONS_QUERY, elem.Id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
