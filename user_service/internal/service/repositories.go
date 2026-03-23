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
	CreateWallet(ctx context.Context, userId uint64) error
	AddMoney(ctx context.Context, userId uint64, amount uint64) error
	WriteOffMoney(ctx context.Context, userId uint64, amount uint64) error
}

type BasketsRepository interface {
	AddProductToBasket(ctx context.Context, p *models.Product) error
	DeleteProductFromBasket(ctx context.Context, userId uint64, productId uint64) error
	GetUserBasket(ctx context.Context, userId uint64) (*models.UserBasketInfo, error)
	ClearUserBasket(ctx context.Context, userId uint64) error
}
