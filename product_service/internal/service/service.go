package service

import (
	"context"
	"product_service/internal/models"
)

type ProductsRepository interface {
	AddProduct(ctx context.Context, p *models.Product) error
	ProductList(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error)
	ProductsById(ctx context.Context, ids []uint64) ([]*models.Product, error)
	ToUpProductQuantity(ctx context.Context, amount uint64, productId uint64) error
	ToDownProductQuantity(ctx context.Context, amount uint64, productId uint64) error
	DeleteProduct(ctx context.Context, productId uint64) error
	ReserveProducts(ctx context.Context, reservations []*models.Reservation) error
	CancelReservationForOrder(ctx context.Context, orderId uint64) error
	DeleteReservationsForOrder(ctx context.Context, orderId uint64) error
}

type ProductService struct {
	products ProductsRepository
}

func NewProductService(products ProductsRepository) *ProductService {
	ps := &ProductService{
		products: products,
	}

	return ps
}

func (ps *ProductService) ShowProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
	products, err := ps.products.ProductList(ctx, pageNumber, itemsPerPage)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) ProductsByIds(ctx context.Context, ids []uint64) ([]*models.Product, error) {
	products, err := ps.products.ProductsById(ctx, ids)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) AddProduct(ctx context.Context, p *models.Product) error {
	err := ps.products.AddProduct(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, productId uint64) error {
	err := ps.products.DeleteProduct(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ToUpProductQuantity(ctx context.Context, productId uint64, quantity uint64) error {
	err := ps.products.ToUpProductQuantity(ctx, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ToDownProductQuantity(ctx context.Context, productId, quantity uint64) error {
	err := ps.products.ToDownProductQuantity(ctx, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ReserveProducts(ctx context.Context, reservations []*models.Reservation) error {
	err := ps.products.ReserveProducts(ctx, reservations)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) CancelReservationsForOrder(ctx context.Context, orderId uint64) error {
	err := ps.products.CancelReservationForOrder(ctx, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) DeleteReservationsForOrder(ctx context.Context, orderId uint64) error {
	err := ps.products.DeleteReservationsForOrder(ctx, orderId)
	if err != nil {
		return err
	}

	return nil
}
