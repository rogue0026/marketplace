package grpcapi

import (
	"context"
	"order_service/internal/service"

	pb "github.com/rogue0026/marketplace-proto/orders"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrdersHandler struct {
	pb.UnimplementedOrderServiceServer
	s *service.OrderService
}

func NewOrdersHandler(s *service.OrderService) *OrdersHandler {
	h := &OrdersHandler{
		s: s,
	}

	return h
}

func (h *OrdersHandler) CreateNewOrder(ctx context.Context, req *pb.CreateNewOrderRequest) (*pb.CreateNewOrderResponse, error) {
	orderId, err := h.s.CreateOrder(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateNewOrderResponse{
		OrderId: orderId,
	}

	return resp, nil
}

func (h *OrdersHandler) PayForOrder(ctx context.Context, req *pb.PayForOrderRequest) (*emptypb.Empty, error) {
	err := h.s.PayForOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
