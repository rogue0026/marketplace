package grpcapi

import (
	"context"
	_ "google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "product_service/api/product/v1"
	"product_service/internal/models"
	"product_service/internal/service"
)

type GRPCRouter struct {
	pb.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewRouter(s *service.ProductService) *GRPCRouter {
	server := &GRPCRouter{
		service: s,
	}
	return server
}

func (s *GRPCRouter) ShowProducts(ctx context.Context, req *pb.ShowProductsRequest) (*pb.ShowProductsResponse, error) {
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

func (s *GRPCRouter) AddProducts(ctx context.Context, req *pb.AddProductsRequest) (*emptypb.Empty, error) {
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

func (s *GRPCRouter) DeleteProducts(ctx context.Context, req *pb.DeleteProductsRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *GRPCRouter) ToUpProducts(ctx context.Context, req *pb.ToUpProductsRequest) (*emptypb.Empty, error) {
	err := s.service.ToUpProduct(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *GRPCRouter) ToDownProducts(ctx context.Context, req *pb.ToDownProductsRequest) (*emptypb.Empty, error) {
	err := s.service.ToDownProduct(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
