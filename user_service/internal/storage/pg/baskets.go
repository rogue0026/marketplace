package pg

import (
	"context"
	"user_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	AddProductToBasketQuery = `
	INSERT INTO basket_content (user_id, product_id, product_quantity, current_price)
	VALUES ($1, $2, $3, $4)
	`

	DeleteProductFromBasketQuery = `
	DELETE FROM basket_content
	WHERE user_id = $1 AND product_id = $2
	`

	GetUserBasketQuery = `
	SELECT
	    id,
	    user_id,
	    product_id,
	    product_quantity,
	    current_price
	FROM basket_content
	WHERE user_id = $1`
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

func (repo *BasketsRepo) AddProductToBasket(ctx context.Context, p *models.Product) error {
	_, err := repo.pool.Exec(
		ctx,
		AddProductToBasketQuery,
		p.UserId,
		p.ProductId,
		p.ProductQuantity,
		p.PricePerUnit,
	)
	if err != nil {
		return err
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
		return err
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
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	var totalPrice uint64
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(
			p.Id,
			p.UserId,
			p.ProductId,
			p.ProductQuantity,
			p.PricePerUnit,
		); err != nil {
			return nil, err
		}
		totalPrice += p.PricePerUnit * p.ProductQuantity
		products = append(products, p)
	}

	basket := &models.UserBasketInfo{
		Products:   products,
		TotalPrice: totalPrice,
	}

	return basket, nil
}
