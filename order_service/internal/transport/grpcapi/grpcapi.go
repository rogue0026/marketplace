package grpcapi

import (
	"context"
	"log/slog"
	"order_service/internal/service"
	"order_service/pkg/logger"

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
	appLogger := logger.Extract(ctx)

	orderId, err := h.s.CreateOrder(ctx, req.UserId)
	if err != nil {
		appLogger.Error(
			"failed to create user",
			slog.Uint64("user_id", req.UserId),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	resp := &pb.CreateNewOrderResponse{
		OrderId: orderId,
	}

	return resp, nil
}

func (h *OrdersHandler) PayForOrder(ctx context.Context, req *pb.PayForOrderRequest) (*emptypb.Empty, error) {
	appLogger := logger.Extract(ctx)

	err := h.s.PayForOrder(ctx, req.OrderId)
	if err != nil {
		appLogger.Error(
			"failed to pay for order",
			slog.Uint64("order_id", req.OrderId),
		)
		return nil, err
	}

	return nil, nil
}
