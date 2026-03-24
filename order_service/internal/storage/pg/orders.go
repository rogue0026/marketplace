package pg

import (
	"context"
	"order_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StatusWaitingForPayment = "waiting for payment"
)

const (
	CreateOrderQuery = `
	INSERT INTO orders (user_id, total_price, status)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	AddOrderContentQuery = `
	INSERT INTO order_contents (order_id, product_id, product_quantity, price_per_unit)
	VALUES ($1, $2, $3, $4)
	`
)

type OrdersRepo struct {
	pool *pgxpool.Pool
}

func NewOrdersRepo(pool *pgxpool.Pool) *OrdersRepo {
	r := &OrdersRepo{
		pool: pool,
	}
	return r
}

func (r *OrdersRepo) CreateOrder(ctx context.Context, o *models.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var orderId uint64

	err = tx.QueryRow(
		ctx,
		CreateOrderQuery,
		o.UserId,
		o.TotalPrice,
		StatusWaitingForPayment,
	).Scan(&orderId)
	if err != nil {
		return err
	}

	for _, orderItem := range o.Items {
		_, err := tx.Exec(
			ctx,
			AddOrderContentQuery,
			orderId,
			orderItem.ProductId,
			orderItem.Quantity,
			orderItem.PricePerUnit,
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

func (r *OrdersRepo) GetOrderInfo(ctx context.Context, orderId uint64) (*models.Order, error) {
	return nil, nil
}

func (r *OrdersRepo) ChangeOrderStatus(ctx context.Context, orderId uint64, status string) error {
	return nil
}
