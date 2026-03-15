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
	ToUpBalance(ctx context.Context)
	ToDownBalance(ctx context.Context)
}

type BasketsRepository interface {
	GetProductsFromBasketByUserId(ctx context.Context, userId uint64) ([]*models.Product, error)
	DeleteProductFromBasket(ctx context.Context)
}

type UserService struct {
	usersRepo   UsersRepository
	walletsRepo WalletsRepository
	basketsRepo BasketsRepository
}

func NewUserService() *UserService {
	return nil
}

func (s *UserService) CreateNewUser(ctx context.Context, user *models.User) (uint64, error) {
	userId, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId uint64) error {
	err := s.usersRepo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ShowProductsFromBasket(ctx context.Context, userId uint64) ([]*models.Product, error) {
	userProducts, err := s.basketsRepo.GetProductsFromBasketByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return userProducts, nil
}

func (s *UserService) AddProductToBasket(ctx context.Context) {}

func (s *UserService) DeleteProductFromBasket(ctx context.Context) {}

func (s *UserService) ToUpBalance(ctx context.Context) {}

func (s *UserService) ToDownBalance(ctx context.Context) {}
