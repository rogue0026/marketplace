/*
создание заказа
1. получили запрос от grpc endpoint'a
2. отправили событие в кафку о том что создается заказ
3. user service отправляет на это событие
*/

package service

import (
	"context"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context, userId uint64) error
	GetUserOrders(ctx context.Context)
}

type EventPublisher interface {
	Publish(ctx context.Context, e Event) error
}

type Event interface {
	GetEventName() string
	GetTopicName() string
}

type OrderService struct {
	repo         OrdersRepository
	eventWriters EventPublisher
}

func (s *OrderService) CreateOrder(ctx context.Context, userId uint64) {
}

func (s *OrderService) GetUserOrders(ctx context.Context, userId uint64) {}
