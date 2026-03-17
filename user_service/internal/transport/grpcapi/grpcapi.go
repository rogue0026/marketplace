package grpcapi

import (
	"context"
	"user_service/api/pb"
	"user_service/internal/models"
	"user_service/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	s *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	h := &Handler{
		s: userService,
	}

	return h
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserId, error) {
	u := &models.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}
	userId, err := h.s.CreateNewUser(ctx, u)
	if err != nil {
		return nil, err
	}

	resp := &pb.UserId{
		UserId: userId,
	}

	return resp, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.UserId) (*emptypb.Empty, error) {
	err := h.s.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) ToUpBalance(ctx context.Context, req *pb.WalletMessage) (*emptypb.Empty, error) {
	err := h.s.ToUpBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) ToDownBalance(ctx context.Context, req *pb.WalletMessage) (*emptypb.Empty, error) {
	err := h.s.ToDownBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) GetProductsFromBasket(ctx context.Context, req *pb.UserId) (*pb.ProductsInBasket, error) {
	products, err := h.s.ShowProductsFromBasket(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, elem := range products {
		p := &pb.Product{
			Id:              elem.Id,
			UserId:          elem.UserId,
			ProductId:       elem.ProductId,
			ProductQuantity: elem.ProductQuantity,
			PricePerUnit:    elem.PricePerUnit,
		}
		pbProducts[i] = p
	}

	resp := &pb.ProductsInBasket{
		Products: pbProducts,
	}

	return resp, nil
}

func (h *Handler) AddProductToBasket(ctx context.Context, p *pb.Product) (*emptypb.Empty, error) {
	product := &models.Product{
		UserId:          p.UserId,
		ProductId:       p.ProductId,
		ProductQuantity: p.ProductQuantity,
		PricePerUnit:    p.PricePerUnit,
	}

	err := h.s.AddProductToBasket(ctx, product)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) DeleteProductFromBasket(ctx context.Context, req *pb.ProductId) (*emptypb.Empty, error) {
	err := h.s.DeleteProductFromBasket(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
