package grpcapi

import (
	"context"
	"product_service/api/pb"
	"product_service/internal/models"
	"product_service/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewHandler(s *service.ProductService) *Handler {
	server := &Handler{
		service: s,
	}
	return server
}

func (s *Handler) ShowProducts(ctx context.Context, req *pb.ShowProductsRequest) (*pb.ShowProductsResponse, error) {
	products, err := s.service.ShowProducts(ctx, req.PageNumber, req.ItemsPerPage)
	if err != nil {
		return nil, err
	}

	serialisedProducts := make([]*pb.Product, 0)

	for _, elem := range products {
		p := &pb.Product{
			Id:            elem.Id,
			Name:          elem.Name,
			PricePerUnit:  elem.PricePerUnit,
			TotalQuantity: elem.TotalQuantity,
		}
		serialisedProducts = append(serialisedProducts, p)
	}

	resp := &pb.ShowProductsResponse{
		Products: serialisedProducts,
	}

	return resp, nil
}

func (s *Handler) AddProducts(ctx context.Context, req *pb.AddProductsRequest) (*emptypb.Empty, error) {
	productList := make([]*models.Product, 0, len(req.Products))
	for _, elem := range req.Products {
		p := &models.Product{
			Id:            elem.Id,
			Name:          elem.Name,
			PricePerUnit:  elem.PricePerUnit,
			TotalQuantity: elem.TotalQuantity,
		}
		productList = append(productList, p)
	}

	err := s.service.AddProducts(ctx, productList)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Handler) DeleteProducts(ctx context.Context, req *pb.DeleteProductsRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Handler) ToUpProducts(ctx context.Context, req *pb.ToUpProductsRequest) (*emptypb.Empty, error) {
	err := s.service.ToUpProduct(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Handler) ToDownProducts(ctx context.Context, req *pb.ToDownProductsRequest) (*emptypb.Empty, error) {
	err := s.service.ToDownProduct(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Handler) ReserveProducts(ctx context.Context, req *pb.ReserveProductsRequest) (*emptypb.Empty, error) {
	reservations := make([]*models.Reservation, 0)
	for _, elem := range req.Reservations {
		r := &models.Reservation{
			OrderId:   elem.OrderId,
			ProductId: elem.ProductId,
			Quantity:  elem.Quantity,
			Status:    int(elem.Status),
			ExpiresAt: elem.ExpiresAt.AsTime(),
		}
		reservations = append(reservations, r)
	}

	err := s.service.ReserveProducts(ctx, reservations)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
