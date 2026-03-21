package pg

import (
	"context"
	"product_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	AddProductsQuery = `
	INSERT INTO products (name, price_per_unit, total_quantity)
    VALUES ($1, $2, $3)
    `

	GetProductsQuery = `
    SELECT 
        id,
        name,
        price_per_unit,
        total_quantity
    FROM products
    ORDER BY id
    LIMIT $1 OFFSET $2
    `

	ToUpProductsQuantityQuery = `
	UPDATE products
	SET total_quantity = total_quantity + $1
	WHERE id = $2
	`

	ToDownProductsQuantityQuery = `
	UPDATE products
	SET total_quantity = total_quantity - $1
	WHERE id = $2
	`

	DeleteProductQuery = `
		DELETE FROM products
		WHERE id = $1;
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

func (r *ProductsRepo) AddProducts(ctx context.Context, products ...*models.Product) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, p := range products {
		_, err := tx.Exec(
			ctx,
			AddProductsQuery,
			p.Name,
			p.PricePerUnit,
			p.TotalQuantity,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepo) GetProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
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
			&p.TotalQuantity,
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

func (r *ProductsRepo) ToUpProductQuantity(ctx context.Context, productId uint64, quantity uint64) error {
	if _, err := r.pool.Exec(ctx, ToUpProductsQuantityQuery, quantity, productId); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepo) ToDownProductQuantity(ctx context.Context, productId uint64, quantity uint64) error {
	if _, err := r.pool.Exec(ctx, ToDownProductsQuantityQuery, quantity, productId); err != nil {
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
