package service

import (
	"context"

	pb "github.com/rogue0026/marketplace-proto/users"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context) (uint64, error)
	DeleteOrder(ctx context.Context, orderId uint64) error
}

type OrderService struct {
	orders OrdersRepository
	users  pb.UserServiceClient
}

func NewOrderService(orders OrdersRepository) *OrderService {
	s := &OrderService{
		orders: orders,
	}

	return s
}

func (s *OrderService) CreateOrder(ctx context.Context, orderData interface{}) error {

	return nil
}
