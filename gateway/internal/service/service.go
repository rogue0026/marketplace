package service

import (
	"context"
	"gateway/internal/clients"
	"gateway/internal/models"

	"github.com/rogue0026/marketplace-proto/products"
	"github.com/rogue0026/marketplace-proto/users"
)

type GatewayService struct {
	c *clients.Clients
}

func NewGatewayService(c *clients.Clients) *GatewayService {
	s := &GatewayService{
		c: c,
	}
	return s
}

func (s *GatewayService) ProductCatalog(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
	req := &products.ShowProductsRequest{
		PageNumber:   pageNumber,
		ItemsPerPage: itemsPerPage,
	}

	resp, err := s.c.ProductsClient.ShowProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Product, 0, len(resp.Products))
	for _, elem := range resp.Products {
		p := models.Product{
			Id:       elem.ProductId,
			Name:     elem.Name,
			Price:    elem.CurrentPrice,
			Quantity: elem.RemainingStock,
		}
		result = append(result, &p)
	}

	return result, nil
}

func (s *GatewayService) NewUser(ctx context.Context, login string, password string) (uint64, error) {
	var userId uint64

	resp, err := s.c.UsersClient.CreateUser(ctx, &users.CreateUserRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return userId, err
	}

	userId = resp.UserId

	return userId, nil
}

func (s *GatewayService) DeleteUser(ctx context.Context, userId uint64) error {
	_, err := s.c.UsersClient.DeleteUser(ctx, &users.DeleteUserRequest{UserId: userId})
	if err != nil {
		return err
	}

	return nil
}

func (s *GatewayService) BasketInfo(ctx context.Context, userId uint64) ([]*models.Product, error) {

	basketData, err := s.c.UsersClient.GetProductsFromBasket(ctx, &users.GetProductsFromBasketRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(basketData.BasketItems))
	productIdToQuantity := make(map[uint64]uint64)
	for _, elem := range basketData.BasketItems {
		ids = append(ids, elem.Id)
		productIdToQuantity[elem.Id] = elem.Quantity
	}

	productDetails, err := s.c.ProductsClient.ShowProductsByIds(ctx, &products.ShowProductsByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	productsInBasket := make([]*models.Product, 0)
	for _, elem := range productDetails.Products {
		p := &models.Product{
			Id:       elem.ProductId,
			Name:     elem.Name,
			Price:    elem.CurrentPrice,
			Quantity: productIdToQuantity[elem.ProductId],
		}
		productsInBasket = append(productsInBasket, p)
	}
	return productsInBasket, nil
}

func (s *GatewayService) AddProductToBasket(ctx context.Context, userId uint64, productId uint64, quantity uint64) error {
	_, err := s.c.UsersClient.AddProductToBasket(ctx, &users.AddProductsToBasketRequest{
		UserId:          userId,
		ProductId:       productId,
		ProductQuantity: quantity,
	})
	if err != nil {
		return err
	}

	return nil
}
