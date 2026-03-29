package pg

import (
	"context"
	"fmt"
	"user_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	AddProductToBasketQuery = `
	INSERT INTO basket_content (user_id, product_id, product_quantity)
	VALUES ($1, $2, $3)
	`

	DeleteProductFromBasketQuery = `
	DELETE FROM basket_content
	WHERE user_id = $1 AND product_id = $2
	`

	GetUserBasketQuery = `
	SELECT
	    product_id,
	    product_quantity
	FROM basket_content
	WHERE user_id = $1`

	ClearUserBasketQuery = `
	DELETE FROM basket_content
	WHERE user_id = $1
	`
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

func (repo *BasketsRepo) AddProductToBasket(ctx context.Context, userId uint64, p *models.Product) error {
	_, err := repo.pool.Exec(
		ctx,
		AddProductToBasketQuery,
		userId,
		p.Id,
		p.Quantity,
	)
	if err != nil {
		return fmt.Errorf("failed to add product to basket for user with user_id=%d: %w", userId, err)
	}

	return nil
}

func (repo *BasketsRepo) DeleteProductFromBasket(ctx context.Context, userId uint64, productId uint64) error {
	_, err := repo.pool.Exec(
		ctx,
		DeleteProductFromBasketQuery,
		userId,
		productId,
	)
	if err != nil {
		return fmt.Errorf("failed to delete product from basket for user with user_id=%d: %w", userId, err)
	}

	return nil
}

func (repo *BasketsRepo) GetUserBasket(ctx context.Context, userId uint64) (*models.UserBasketInfo, error) {
	rows, err := repo.pool.Query(
		ctx,
		GetUserBasketQuery,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rows with basket info for user with user_id=%d: %w", userId, err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		p := models.Product{}
		if err := rows.Scan(
			&p.Id,
			&p.Quantity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row while getting user basket for user_id=%d: %w", userId, err)
		}
		products = append(products, &p)
	}

	basket := &models.UserBasketInfo{
		Products: products,
	}

	return basket, nil
}

func (repo *BasketsRepo) ClearUserBasket(ctx context.Context, userId uint64) error {
	_, err := repo.pool.Exec(ctx, ClearUserBasketQuery, userId)
	if err != nil {
		return fmt.Errorf("failed to clear user basket for user_id=%d: %w", userId, err)
	}

	return nil
}
