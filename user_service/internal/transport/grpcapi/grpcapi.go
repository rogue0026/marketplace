package grpcapi

import (
	"context"
	pb "github.com/rogue0026/marketplace-proto/users"
	"google.golang.org/protobuf/types/known/emptypb"
	"user_service/internal/models"
	"user_service/internal/service"
)

type UserServiceHandler struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
}

func (s *UserServiceHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	usr := &models.User{
		Username:     req.Login,
		PasswordHash: req.Password,
	}
	_, err := s.service.CreateNewUser(ctx, usr)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) AddMoney(ctx context.Context, req *pb.AddMoneyRequest) (*emptypb.Empty, error) {
	err := s.service.ToUpBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) WriteOffMoney(ctx context.Context, req *pb.WriteOffMoneyRequest) (*emptypb.Empty, error) {
	err := s.service.ToDownBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) GetProductsFromBasket(ctx context.Context, req *pb.GetProductsFromBasketRequest) (*pb.GetProductsFromBasketResponse, error) {
	basketInfo, err := s.service.ShowProductsFromBasket(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetProductsFromBasketResponse{UserId: 1000} // todo
	return resp, nil
}

func (s *UserServiceHandler) AddProductToBasket(ctx context.Context, req *pb.AddProductsToBasketRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *UserServiceHandler) ClearUserBasket(ctx context.Context, req *pb.ClearUserBasketRequest) (*emptypb.Empty, error) {
	return nil, nil
}
