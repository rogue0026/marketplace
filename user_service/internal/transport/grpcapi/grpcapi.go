package grpcapi

import (
	"context"
	"log/slog"
	"user_service/internal/models"
	"user_service/internal/service"
	"user_service/pkg/logger"

	pb "github.com/rogue0026/marketplace-proto/users"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceHandler struct {
	pb.UnimplementedUserServiceServer
	s *service.UserService
}

func NewHandler(s *service.UserService) *UserServiceHandler {
	h := &UserServiceHandler{
		s: s,
	}

	return h
}

func (s *UserServiceHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	appLogger := logger.Extract(ctx)

	usr := &models.User{
		Username:     req.Login,
		PasswordHash: req.Password,
	}
	userId, err := s.s.CreateNewUser(ctx, usr)
	if err != nil {
		appLogger.Error(
			"failed to create user",
			slog.String("login", req.Login),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}
	resp := &pb.CreateUserResponse{
		UserId: userId,
	}

	return resp, nil
}

func (s *UserServiceHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	appLogger := logger.Extract(ctx)

	err := s.s.DeleteUser(ctx, req.UserId)
	if err != nil {
		appLogger.Error(
			"failed to delete user",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) AddMoney(ctx context.Context, req *pb.AddMoneyRequest) (*emptypb.Empty, error) {
	appLogger := logger.Extract(ctx)

	err := s.s.ToUpBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		appLogger.Error(
			"failed to add money for user",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) WriteOffMoney(ctx context.Context, req *pb.WriteOffMoneyRequest) (*emptypb.Empty, error) {
	appLogger := logger.Extract(ctx)

	err := s.s.ToDownBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		appLogger.Error(
			"failed to write off money for user",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) GetProductsFromBasket(ctx context.Context, req *pb.GetProductsFromBasketRequest) (*pb.GetProductsFromBasketResponse, error) {
	appLogger := logger.Extract(ctx)

	basketInfo, err := s.s.ShowProductsFromBasket(ctx, req.UserId)
	if err != nil {
		appLogger.Error(
			"failed to get products from basket",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	basketProducts := make([]*pb.BasketProduct, 0)
	for _, elem := range basketInfo.Products {
		p := &pb.BasketProduct{
			Id:       elem.Id,
			Quantity: elem.Quantity,
		}
		basketProducts = append(basketProducts, p)
	}

	resp := &pb.GetProductsFromBasketResponse{
		BasketItems: basketProducts,
	}

	return resp, nil
}

func (s *UserServiceHandler) AddProductToBasket(ctx context.Context, req *pb.AddProductsToBasketRequest) (*emptypb.Empty, error) {
	appLogger := logger.Extract(ctx)

	p := &models.Product{
		Id:       req.ProductId,
		Quantity: req.ProductQuantity,
	}

	err := s.s.AddProductToBasket(ctx, req.UserId, p)
	if err != nil {
		appLogger.Error(
			"failed to add product to basket",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) ClearUserBasket(ctx context.Context, req *pb.ClearUserBasketRequest) (*emptypb.Empty, error) {
	appLogger := logger.Extract(ctx)

	err := s.s.ClearUserBasket(ctx, req.UserId)
	if err != nil {
		appLogger.Error(
			"failed to clear user basket",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	return nil, nil
}
