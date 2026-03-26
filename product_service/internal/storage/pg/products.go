package pg

import (
	"context"
	"product_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	AddProductsQuery = `
	INSERT INTO products (name, stock_remaining, current_price)
    VALUES ($1, $2, $3)
    `

	GetProductsQuery = `
	SELECT
		id,
		name,
		current_price,
		stock_remaining
	FROM products
	ORDER BY id
	LIMIT $1 OFFSET $2
	`

	ProductsByIdQuery = `
	SELECT
		id,
		name,
		current_price
	FROM products
	WHERE id = ANY($1)
	`

	ToUpProductsQuantityQuery = `
	UPDATE products
	SET stock_remaining = stock_remaining + $1
	WHERE id = $2
	`

	ToDownProductsQuantityQuery = `
	UPDATE products
	SET stock_remaining = stock_remaining - $1
	WHERE id = $2
	`

	DeleteProductQuery = `
		DELETE FROM products
		WHERE id = $1;
	`

	ReserveProductQuery = `
	INSERT INTO product_reservations (order_id, product_id, quantity)
	VALUES ($1, $2, $3)`

	DeleteOrderReservationsQuery = `
	DELETE FROM product_reservations
	WHERE order_id = $1
	`
)

type ProductsRepo struct {
	pool *pgxpool.Pool
}

func NewProductsRepo(pool *pgxpool.Pool) *ProductsRepo {
	r := &ProductsRepo{
		pool: pool,
	}

	return r
}

func (r *ProductsRepo) AddProduct(ctx context.Context, p *models.Product) error {
	_, err := r.pool.Exec(ctx, AddProductsQuery, p.Name, p.StockRemaining, p.PricePerUnit)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductsRepo) ProductList(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
	offset := (pageNumber - 1) * itemsPerPage

	rows, err := r.pool.Query(ctx, GetProductsQuery, itemsPerPage, offset)
	if err != nil {
		return nil, err
	}

	products := make([]*models.Product, 0)
	for rows.Next() {
		p := &models.Product{}
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.PricePerUnit,
			&p.StockRemaining,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductsRepo) ProductsById(ctx context.Context, ids []uint64) ([]*models.Product, error) {
	rows, err := r.pool.Query(
		ctx,
		ProductsByIdQuery,
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.Product, 0)
	for rows.Next() {
		p := models.Product{}
		err = rows.Scan(
			&p.Id,
			&p.Name,
			&p.PricePerUnit,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &p)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r *ProductsRepo) ToUpProductQuantity(ctx context.Context, amount uint64, productId uint64) error {
	if _, err := r.pool.Exec(ctx, ToUpProductsQuantityQuery, amount, productId); err != nil {
		return err
	}
	return nil
}

func (r *ProductsRepo) ToDownProductQuantity(ctx context.Context, amount uint64, productId uint64) error {
	if _, err := r.pool.Exec(ctx, ToDownProductsQuantityQuery, amount, productId); err != nil {
		return err
	}
	return nil
}

func (r *ProductsRepo) DeleteProduct(ctx context.Context, productId uint64) error {
	if _, err := r.pool.Exec(ctx, DeleteProductQuery, productId); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepo) ReserveProducts(ctx context.Context, reservations []*models.Reservation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	for _, r := range reservations {
		_, err = tx.Exec(ctx, ReserveProductQuery, r.OrderId, r.ProductId, r.Quantity)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, ToDownProductsQuantityQuery, r.Quantity, r.ProductId)
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

func (r *ProductsRepo) CancelReservationForOrder(ctx context.Context, orderId uint64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(
		ctx,
		`SELECT product_id, quantity FROM product_reservations WHERE order_id = $1`,
		orderId,
	)
	if err != nil {
		return err
	}
	for rows.Next() {
		var productId uint64
		var quantity uint64
		err := rows.Scan(&productId, &quantity)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, ToUpProductsQuantityQuery, quantity, productId)
		if err != nil {
			return err
		}

	}

	err = rows.Err()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, DeleteOrderReservationsQuery, orderId)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepo) DeleteReservationsForOrder(ctx context.Context, orderId uint64) error {
	_, err := r.pool.Exec(
		ctx,
		DeleteOrderReservationsQuery,
		orderId,
	)
	if err != nil {
		return err
	}

	return nil
}
