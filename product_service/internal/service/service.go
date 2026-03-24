package service

import (
	"context"
	"product_service/internal/models"
)

type ProductsRepository interface {
	AddProduct(ctx context.Context, p *models.Product) error
	GetProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error)
	ToUpProductQuantity(ctx context.Context, amount uint64, productId uint64) error
	ToDownProductQuantity(ctx context.Context, amount uint64, productId uint64) error
	DeleteProduct(ctx context.Context, productId uint64) error
	ReserveProducts(ctx context.Context, reservations []*models.Reservation) error
	CancelReservationForOrder(ctx context.Context, orderId uint64) error
	DeleteReservationsForOrder(ctx context.Context, orderId uint64) error
}

type ProductService struct {
	Products ProductsRepository
}

func NewProductService(products ProductsRepository) *ProductService {
	ps := &ProductService{
		Products: products,
	}

	return ps
}

func (ps *ProductService) ShowProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
	products, err := ps.Products.GetProducts(ctx, pageNumber, itemsPerPage)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) AddProduct(ctx context.Context, p *models.Product) error {
	err := ps.Products.AddProduct(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, productId uint64) error {
	err := ps.Products.DeleteProduct(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ToUpProductQuantity(ctx context.Context, productId uint64, quantity uint64) error {
	err := ps.Products.ToUpProductQuantity(ctx, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ToDownProductQuantity(ctx context.Context, productId, quantity uint64) error {
	err := ps.Products.ToDownProductQuantity(ctx, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ReserveProducts(ctx context.Context, reservations []*models.Reservation) error {
	err := ps.Products.ReserveProducts(ctx, reservations)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) CancelReservationsForOrder(ctx context.Context, orderId uint64) error {
	err := ps.Products.CancelReservationForOrder(ctx, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) DeleteReservationsForOrder(ctx context.Context, orderId uint64) error {
	err := ps.Products.DeleteReservationsForOrder(ctx, orderId)
	if err != nil {
		return err
	}

	return nil
}
