package service

import (
	"context"
	"order_service/internal/models"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context, orderData *models.OrderInfo) error
	GetUserOrders(ctx context.Context, userId uint64) ([]*models.OrderInfo, error)
}

type OrderService struct {
	orders OrdersRepository
}

func (s *OrderService) CreateOrder(ctx context.Context, info *models.OrderInfo) error {
	/*
		1. просто создает запись в базе данных с информацией о заказе (на этом этапе заказ содержит всю необходимую для оплаты информацию)
		2. далее нужно оплатить заказ (оплата подразумевает резервирование товара на складе и резервирование денег на кошельке)
		3. если резервирование прошло успешно то вызываем метод для списания товара со склада и признания выручки
	*/
	// создаем заказ в базе, здесь он пока что не оплачен
	err := s.orders.CreateOrder(ctx, info)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) PayForTheOrder(ctx context.Context, orderId uint64) error {
	
	// тут пытаемся зарезервировать товар на складе
	// если товар есть то пытаемся зарезервировать деньги в кошельке
	// если деньги удалось зарезервировать то производим списание денег с баланса пользователя и
	return nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userId uint64) {}
