package service

import (
	"context"
	"user_service/internal/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *models.User) (uint64, error)
	DeleteUser(ctx context.Context, userId uint64) error
}

type WalletsRepository interface {
	ToUpBalance(ctx context.Context, userId uint64, amount uint64) error
	ToDownBalance(ctx context.Context, userId uint64, amount uint64) error
}

type BasketsRepository interface {
	GetProductsByUserId(ctx context.Context, userId uint64) ([]*models.Product, error)
	AddProductToUserBasket(ctx context.Context, product *models.Product) error
	DeleteProductFromBasket(ctx context.Context, productId uint64) error
}
