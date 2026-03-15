package pg

import (
	"context"
	"product_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductsRepository struct {
	pool *pgxpool.Pool
}

func NewProductsRepository(pool *pgxpool.Pool) *ProductsRepository {
	pr := &ProductsRepository{
		pool: pool,
	}

	return pr
}

func (r *ProductsRepository) AddProducts(ctx context.Context, products ...*models.Product) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
        INSERT INTO products (name, price_per_unit, total_quantity)
        VALUES ($1, $2, $3)
    `

	for _, p := range products {
		_, err := tx.Exec(
			ctx,
			query,
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

func (r *ProductsRepository) GetProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
	sqlQuery := `
        SELECT id, name, price_per_unit, total_quantity
        FROM products
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	offset := (pageNumber - 1) * itemsPerPage
	rows, err := r.pool.Query(ctx, sqlQuery, itemsPerPage, offset)
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

func (r *ProductsRepository) ToUpProductQuantity(ctx context.Context, productId uint64, quantity uint64) error {
	sqlQuery := `
		UPDATE products
		SET total_quantity = total_quantity + $1
		WHERE id = $2
	`

	if _, err := r.pool.Exec(ctx, sqlQuery, quantity, productId); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepository) ToDownProductQuantity(ctx context.Context, productId uint64, quantity uint64) error {
	sqlQuery := `
		UPDATE products
		SET total_quantity = total_quantity - $1
		WHERE id = $2
	`

	if _, err := r.pool.Exec(ctx, sqlQuery, quantity, productId); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepository) DeleteProduct(ctx context.Context, productId uint64) error {
	sqlQuery := `
		DELETE FROM products
		WHERE id = $1;
	`

	if _, err := r.pool.Exec(ctx, sqlQuery, productId); err != nil {
		return err
	}

	return nil
}
