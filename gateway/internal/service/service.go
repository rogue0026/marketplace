package service

import (
	"context"
	"gateway/internal/clients"
	"gateway/internal/models"

	"github.com/rogue0026/marketplace-proto/products"
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
			Id:             elem.ProductId,
			Name:           elem.Name,
			Price:          elem.CurrentPrice,
			StockRemaining: elem.RemainingStock,
		}
		result = append(result, &p)
	}

	return result, nil
}
