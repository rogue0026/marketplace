package grpcapi

import (
	"context"
	"product_service/internal/models"
	"product_service/internal/service"

	pb "github.com/rogue0026/marketplace-proto/products"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	s *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	h := &ProductHandler{
		s: s,
	}

	return h
}

func (h *ProductHandler) ShowProducts(ctx context.Context, req *pb.ShowProductsRequest) (*pb.ShowProductsResponse, error) {
	products, err := h.s.ShowProducts(ctx, req.PageNumber, req.ItemsPerPage)
	if err != nil {
		return nil, err
	}
	pbProducts := make([]*pb.Product, 0)
	for _, p := range products {
		pbProduct := &pb.Product{
			ProductId:      p.Id,
			Name:           p.Name,
			RemainingStock: p.StockRemaining,
			CurrentPrice:   p.PricePerUnit,
		}
		pbProducts = append(pbProducts, pbProduct)
	}
	resp := &pb.ShowProductsResponse{
		Products: pbProducts,
	}

	return resp, nil
}

func (h *ProductHandler) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	p := &models.Product{
		Name:           req.Name,
		PricePerUnit:   req.CurrentPrice,
		StockRemaining: req.RemainingStock,
	}
	err := h.s.AddProduct(ctx, p)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	err := h.s.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ProductHandler) IncreaseProductQuantity(ctx context.Context, req *pb.IncreaseProductRequest) (*emptypb.Empty, error) {
	err := h.s.ToUpProductQuantity(ctx, req.ProductId, req.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ProductHandler) DecreaseProductQuantity(ctx context.Context, req *pb.DecreaseProductRequest) (*emptypb.Empty, error) {
	err := h.s.ToDownProductQuantity(ctx, req.ProductId, req.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ProductHandler) ReserveProducts(ctx context.Context, req *pb.ReserveProductsRequest) (*emptypb.Empty, error) {
	reservations := make([]*models.Reservation, 0)
	for _, elem := range req.Reservations {
		r := &models.Reservation{
			OrderId:   elem.OrderId,
			ProductId: elem.ProductId,
			Quantity:  elem.Amount,
		}
		reservations = append(reservations, r)
	}

	err := h.s.ReserveProducts(ctx, reservations)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ProductHandler) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*emptypb.Empty, error) {
	err := h.s.CancelReservationsForOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ProductHandler) ConfirmReservation(ctx context.Context, req *pb.ConfirmReservationRequest) (*emptypb.Empty, error) {
	return nil, nil
}

/* service layer
ShowProducts(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]*models.Product, error)
AddProduct(ctx context.Context, p *models.Product) error
DeleteProduct(ctx context.Context, productId uint64) error
ToUpProductQuantity(ctx context.Context, productId uint64, quantity uint64) error
ToDownProductQuantity(ctx context.Context, productId, quantity uint64) error
ReserveProducts(ctx context.Context, reservations []*models.Reservation) error
CancelReservationsForOrder(ctx context.Context, orderId uint64) error
*/

/* grpc layer
ShowProducts(context.Context, *ShowProductsRequest) (*ShowProductsResponse, error)
AddProduct(context.Context, *AddProductRequest) (*AddProductResponse, error)
DeleteProduct(context.Context, *DeleteProductRequest) (*emptypb.Empty, error)
func (h *ProductHandler) IncreaseProductQuantity(ctx context.Context, req *pb.IncreaseProductRequest) (*pb.emptypb.Empty, error) {}
func (h *ProductHandler) DecreaseProductQuantity(ctx context.Context, req *pb.DecreaseProductRequest) (*pb.emptypb.Empty, error) {}
func (h *ProductHandler) ReserveProducts(ctx context.Context, req *pb.ReserveProductsRequest) (*pb.ReserveProductResponse, error) {}
func (h *ProductHandler) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*pb.emptypb.Empty, error) {}
func (h *ProductHandler) ConfirmReservation(ctx context.Context, req *pb.ConfirmReservationRequest) (*pb.emptypb.Empty, error) {}
*/
