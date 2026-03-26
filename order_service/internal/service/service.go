package service

import (
	"context"
	"order_service/internal/models"

	pbproducts "github.com/rogue0026/marketplace-proto/products"
	pbusers "github.com/rogue0026/marketplace-proto/users"
)

const (
	StatusOrderWaitingForPayment = "waiting for payment"
	StatusOrderInProcessing      = "in processing"
	StatusOrderPayedSuccessfully = "payed successfully"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context, o *models.Order) (uint64, error)
	GetOrderInfo(ctx context.Context, orderId uint64) (*models.Order, error)
	ChangeOrderStatus(ctx context.Context, orderId uint64, status string) error
}

type OrderService struct {
	orders         OrdersRepository
	usersClient    pbusers.UserServiceClient
	productsClient pbproducts.ProductServiceClient
}

func NewOrderService(
	orders OrdersRepository,
	usersClient pbusers.UserServiceClient,
	productsClient pbproducts.ProductServiceClient,
) (*OrderService, error) {

	s := &OrderService{
		orders:         orders,
		usersClient:    usersClient,
		productsClient: productsClient,
	}

	return s, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, userId uint64) (uint64, error) {
	var orderId uint64
	basket, err := s.usersClient.GetProductsFromBasket(ctx, &pbusers.GetProductsFromBasketRequest{
		UserId: userId,
	})
	if err != nil {
		return orderId, err
	}

	ids := make([]uint64, 0, len(basket.BasketItems))
	productIdToQuantity := make(map[uint64]uint64)
	for _, item := range basket.BasketItems {
		ids = append(ids, item.Id)
		productIdToQuantity[item.Id] = item.Quantity
	}

	productsInfo, err := s.productsClient.ShowProductsByIds(ctx, &pbproducts.ShowProductsByIdsRequest{Ids: ids})
	if err != nil {
		return orderId, err
	}

	orderItems := make([]*models.OrderItem, 0, len(productsInfo.Products))
	var totalPrice uint64
	for _, p := range productsInfo.Products {
		totalPrice = totalPrice + p.CurrentPrice*productIdToQuantity[p.ProductId]
		item := &models.OrderItem{
			ProductId:    p.ProductId,
			Quantity:     productIdToQuantity[p.ProductId],
			PricePerUnit: p.CurrentPrice,
		}
		orderItems = append(orderItems, item)
	}

	o := &models.Order{
		UserId:     userId,
		Items:      orderItems,
		TotalPrice: totalPrice,
	}

	orderId, err = s.orders.CreateOrder(ctx, o)
	if err != nil {
		return orderId, err
	}

	_, err = s.usersClient.ClearUserBasket(ctx, &pbusers.ClearUserBasketRequest{UserId: userId})
	if err != nil {
		return orderId, err
	}

	return orderId, nil
}

func (s *OrderService) PayForOrder(ctx context.Context, orderId uint64) error {
	// изменяем статус заказа
	if err := s.orders.ChangeOrderStatus(ctx, orderId, StatusOrderInProcessing); err != nil {
		return err
	}

	// получаем содержимое заказа
	orderInfo, err := s.orders.GetOrderInfo(ctx, orderId)
	if err != nil {
		return err
	}

	// готовим запрос на резервирование товаров
	productReservations := make([]*pbproducts.Reservation, 0)
	for _, item := range orderInfo.Items {
		r := &pbproducts.Reservation{
			ProductId: item.ProductId,
			Amount:    item.Quantity,
			OrderId:   orderId,
		}
		productReservations = append(productReservations, r)
	}

	reserveProductsReq := &pbproducts.ReserveProductsRequest{
		Reservations: productReservations,
	}

	// резервируем товары для заказа
	_, err = s.productsClient.ReserveProducts(ctx, reserveProductsReq)
	if err != nil {
		return err
	}

	// пробуем списать деньги с кошелька пользователя
	writeOffMoneyReq := &pbusers.WriteOffMoneyRequest{
		UserId: orderInfo.UserId,
		Amount: orderInfo.TotalPrice,
	}
	_, err = s.usersClient.WriteOffMoney(ctx, writeOffMoneyReq)
	if err != nil {
		// если списать деньги не удалось, то возвращаем товары обратно на склад
		_, err = s.productsClient.CancelReservation(ctx, &pbproducts.CancelReservationRequest{
			OrderId: orderId,
		})

		if err != nil {
			return err
		}
	}

	// изменяем статус заказа на оплаченный и удаляем данные о зарезервированных товарах
	err = s.orders.ChangeOrderStatus(ctx, orderId, StatusOrderPayedSuccessfully)
	if err != nil {
		return err
	}

	// удаляем записи зарезервированных товаров
	_, err = s.productsClient.DeleteReservation(
		ctx,
		&pbproducts.DeleteReservationRequest{
			OrderId: orderId,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
