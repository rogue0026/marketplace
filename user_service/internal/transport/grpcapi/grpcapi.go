package grpcapi

import (
	"context"
	"user_service/internal/models"
	"user_service/internal/service"

	pb "github.com/rogue0026/marketplace-proto/users"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceHandler struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
}

func (s *UserServiceHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	usr := &models.User{
		Username:     req.Login,
		PasswordHash: req.Password,
	}
	userId, err := s.service.CreateNewUser(ctx, usr)
	if err != nil {
		return nil, err
	}
	resp := &pb.CreateUserResponse{
		UserId: userId,
	}

	return resp, nil
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

	basketProducts := make([]*pb.BasketProduct, 0)
	for _, elem := range basketInfo.Products {
		p := &pb.BasketProduct{
			Id:              elem.Id,
			UserId:          elem.UserId,
			ProductId:       elem.ProductId,
			ProductQuantity: elem.ProductQuantity,
			PricePerUnit:    elem.PricePerUnit,
		}
		basketProducts = append(basketProducts, p)
	}

	resp := &pb.GetProductsFromBasketResponse{
		UserId:      req.UserId,
		BasketItems: basketProducts,
		TotalPrice:  basketInfo.TotalPrice,
	}

	return resp, nil
}

func (s *UserServiceHandler) AddProductToBasket(ctx context.Context, req *pb.AddProductsToBasketRequest) (*emptypb.Empty, error) {
	p := &models.Product{
		Id:              req.UserId,
		UserId:          req.UserId,
		ProductId:       req.ProductId,
		ProductQuantity: req.ProductQuantity,
		PricePerUnit:    req.PricePerUnit,
	}

	err := s.service.AddProductToBasket(ctx, p)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *UserServiceHandler) ClearUserBasket(ctx context.Context, req *pb.ClearUserBasketRequest) (*emptypb.Empty, error) {
	err := s.service.ClearUserBasket(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
