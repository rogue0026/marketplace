package service

import (
	"context"
	"product_service/internal/models"
)

type ProductsRepository interface {
	AddProducts(ctx context.Context, products ...*models.Product) error
	GetProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error)
	ToUpProductQuantity(ctx context.Context, productId uint64, quantity uint64) error
	ToDownProductQuantity(ctx context.Context, productId uint64, quantity uint64) error
	DeleteProduct(ctx context.Context, productId uint64) error
}

type ProductService struct {
	Repository ProductsRepository
}

func NewProductService(repo ProductsRepository) *ProductService {
	ps := &ProductService{
		Repository: repo,
	}

	return ps
}

func (ps *ProductService) ShowProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error) {
	products, err := ps.Repository.GetProducts(ctx, pageNumber, itemsPerPage)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) AddProducts(ctx context.Context, productsList []*models.Product) error {
	err := ps.Repository.AddProducts(ctx, productsList...)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, productId uint64) error {
	err := ps.Repository.DeleteProduct(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ToUpProduct(ctx context.Context, productId uint64, quantity uint64) error {
	err := ps.Repository.ToUpProductQuantity(ctx, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) ToDownProduct(ctx context.Context, productId, quantity uint64) error {
	err := ps.Repository.ToDownProductQuantity(ctx, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}
