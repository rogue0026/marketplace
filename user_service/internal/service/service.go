package service

import (
	"context"
	"fmt"
	"user_service/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	users   UsersRepository
	wallets WalletsRepository
	baskets BasketsRepository
}

func NewUserService(users UsersRepository, wallets WalletsRepository, baskets BasketsRepository) *UserService {
	us := &UserService{
		users:   users,
		wallets: wallets,
		baskets: baskets,
	}

	return us
}

func (s *UserService) CreateNewUser(ctx context.Context, user *models.User) (uint64, error) {
	plainPassword := user.PasswordHash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password for user with login=%s: %w", user.Username, err)
	}

	user.PasswordHash = string(hashedPassword)

	userId, err := s.users.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId uint64) error {
	err := s.users.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ShowProductsFromBasket(ctx context.Context, userId uint64) (*models.UserBasketInfo, error) {
	basketInfo, err := s.baskets.GetUserBasket(ctx, userId)
	if err != nil {
		return nil, err
	}

	return basketInfo, nil
}

func (s *UserService) AddProductToBasket(ctx context.Context, userId uint64, product *models.Product) error {
	err := s.baskets.AddProductToBasket(ctx, userId, product)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteProductFromBasket(ctx context.Context, userId uint64, productId uint64) error {
	err := s.baskets.DeleteProductFromBasket(ctx, userId, productId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ClearUserBasket(ctx context.Context, userId uint64) error {
	err := s.baskets.ClearUserBasket(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ToUpBalance(ctx context.Context, userId uint64, amount uint64) error {
	err := s.wallets.AddMoney(ctx, userId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ToDownBalance(ctx context.Context, userId uint64, amount uint64) error {
	err := s.wallets.WriteOffMoney(ctx, userId, amount)
	if err != nil {
		return err
	}

	return nil
}
