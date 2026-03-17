package service

import (
	"context"
	"user_service/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	usersRepo   UsersRepository
	walletsRepo WalletsRepository
	basketsRepo BasketsRepository
}

func NewUserService(users UsersRepository, wallets WalletsRepository, baskets BasketsRepository) *UserService {
	us := &UserService{
		usersRepo:   users,
		walletsRepo: wallets,
		basketsRepo: baskets,
	}

	return us
}

func (s *UserService) CreateNewUser(ctx context.Context, user *models.User) (uint64, error) {
	plainPassword := user.PasswordHash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.PasswordHash = string(hashedPassword)

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
	userProducts, err := s.basketsRepo.GetProductsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return userProducts, nil
}

func (s *UserService) AddProductToBasket(ctx context.Context, product *models.Product) error {
	err := s.basketsRepo.AddProductToUserBasket(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteProductFromBasket(ctx context.Context, productId uint64) error {
	err := s.basketsRepo.DeleteProductFromBasket(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ToUpBalance(ctx context.Context, userId uint64, amount uint64) error {
	err := s.walletsRepo.ToUpBalance(ctx, userId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ToDownBalance(ctx context.Context, userId uint64, amount uint64) error {
	err := s.walletsRepo.ToDownBalance(ctx, userId, amount)
	if err != nil {
		return err
	}

	return nil
}
