package clients

import (
	"gateway/internal/config"

	"github.com/rogue0026/marketplace-proto/orders"
	"github.com/rogue0026/marketplace-proto/products"
	"github.com/rogue0026/marketplace-proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	UsersClient    users.UserServiceClient
	ProductsClient products.ProductServiceClient
	OrdersClient   orders.OrderServiceClient
}

func New(cfg *config.GRPCClientsConfig) (*Clients, error) {
	var dialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())

	ccUsers, err := grpc.NewClient(cfg.UsersClientAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	ccProducts, err := grpc.NewClient(cfg.ProductsClientAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	ccOrders, err := grpc.NewClient(cfg.OrdersClientAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	usersClient := users.NewUserServiceClient(ccUsers)
	productsClient := products.NewProductServiceClient(ccProducts)
	ordersClient := orders.NewOrderServiceClient(ccOrders)

	c := &Clients{
		UsersClient:    usersClient,
		ProductsClient: productsClient,
		OrdersClient:   ordersClient,
	}

	return c, nil
}
