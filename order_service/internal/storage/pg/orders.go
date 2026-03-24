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

	GetOrderInfoQuery = `
	SELECT
	    user_id,
	    total_price
	FROM orders WHERE id = $1;
	`

	GetOrderContentQuery = `
	SELECT 
	    product_id,
	    product_quantity,
	    price_per_unit
	FROM order_contents WHERE order_id = $1
	`

	ChangeOrderStatusQuery = `
	UPDATE orders
	SET status = $2
	WHERE id = $1
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

func (r *OrdersRepo) CreateOrder(ctx context.Context, o *models.Order) (uint64, error) {
	var orderId uint64

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return orderId, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		CreateOrderQuery,
		o.UserId,
		o.TotalPrice,
		StatusWaitingForPayment,
	).Scan(&orderId)
	if err != nil {
		return orderId, err
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
			return orderId, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return orderId, err
	}

	return orderId, nil
}

func (r *OrdersRepo) GetOrderInfo(ctx context.Context, orderId uint64) (*models.Order, error) {
	var userId uint64
	var totalPrice uint64
	err := r.pool.QueryRow(
		ctx,
		GetOrderInfoQuery,
		orderId,
	).Scan(&userId, &totalPrice)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(
		ctx,
		GetOrderContentQuery,
		orderId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderItems := make([]*models.OrderItem, 0)
	for rows.Next() {
		item := models.OrderItem{}

		err := rows.Scan(
			&item.ProductId,
			&item.Quantity,
			&item.PricePerUnit,
		)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, &item)
	}

	o := &models.Order{
		OrderId:    orderId,
		UserId:     userId,
		Items:      orderItems,
		TotalPrice: totalPrice,
	}

	return o, nil
}

func (r *OrdersRepo) ChangeOrderStatus(ctx context.Context, orderId uint64, status string) error {
	_, err := r.pool.Exec(
		ctx,
		ChangeOrderStatusQuery,
		orderId,
		status,
	)
	if err != nil {
		return err
	}

	return nil
}
