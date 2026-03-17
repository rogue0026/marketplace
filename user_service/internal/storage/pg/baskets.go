package pg

import (
	"context"
	"user_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BasketsRepo struct {
	pool *pgxpool.Pool
}

func NewBasketsRepo(pool *pgxpool.Pool) *BasketsRepo {
	repo := &BasketsRepo{
		pool: pool,
	}

	return repo
}

func (r *BasketsRepo) GetProductsByUserId(ctx context.Context, userId uint64) ([]*models.Product, error) {
	sqlQuery := `
	SELECT 
		id,
		user_id,
		product_id,
		quantity,
		price_per_unit
	FROM basket_products
	WHERE user_id = $1;
	`
	rows, err := r.pool.Query(ctx, sqlQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		p := models.Product{}
		err := rows.Scan(
			&p.Id,
			&p.UserId,
			&p.ProductId,
			&p.ProductQuantity,
			&p.PricePerUnit,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *BasketsRepo) AddProductToUserBasket(ctx context.Context, product *models.Product) error {
	sqlQuery := `
	INSERT INTO basket_products (user_id, product_id, quantity, price_per_unit)
	VALUES ($1, $2, $3, $4);
	`

	_, err := r.pool.Exec(
		ctx,
		sqlQuery,
		product.UserId,
		product.ProductId,
		product.ProductQuantity,
		product.PricePerUnit,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *BasketsRepo) DeleteProductFromBasket(ctx context.Context, recordId uint64) error {
	sqlQuery := `
	DELETE FROM basket_products
	WHERE id = $1;
	`

	_, err := r.pool.Exec(
		ctx,
		sqlQuery,
		recordId,
	)
	if err != nil {
		return err
	}

	return nil
}
