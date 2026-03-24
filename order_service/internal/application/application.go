package application

import (
	"context"
	"fmt"
	"net"
	"order_service/internal/config"
	"order_service/internal/service"
	"order_service/internal/storage/pg"
	"order_service/internal/transport/grpcapi"
	"order_service/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/rogue0026/marketplace-proto/orders"
	pbproducts "github.com/rogue0026/marketplace-proto/products"
	pbusers "github.com/rogue0026/marketplace-proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Application struct {
	Cfg         *config.AppConfig
	connPool    *pgxpool.Pool
	s           *service.OrderService
	tcpListener net.Listener
	grpcServer  *grpc.Server
}

func New() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	pool, err := postgres.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	orders := pg.NewOrdersRepo(pool)

	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	ccUsers, err := grpc.NewClient(cfg.UsersAddr, dialOpts...)
	if err != nil {
		return nil, err
	}

	usersClient := pbusers.NewUserServiceClient(ccUsers)

	ccProducts, err := grpc.NewClient(cfg.ProductsAddr, dialOpts...)
	if err != nil {
		return nil, err
	}

	productsClient := pbproducts.NewProductServiceClient(ccProducts)

	s, err := service.NewOrderService(orders, usersClient, productsClient)
	if err != nil {
		return nil, err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	h := grpcapi.NewOrdersHandler(s)
	pb.RegisterOrderServiceServer(grpcServer, h)

	listener, err := net.Listen("tcp", cfg.GRPCAddr)

	app := &Application{
		Cfg:         cfg,
		connPool:    pool,
		s:           s,
		tcpListener: listener,
		grpcServer:  grpcServer,
	}

	return app, nil
}

func (a *Application) Run() {
	go func() {
		err := a.grpcServer.Serve(a.tcpListener)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()
}

func (a *Application) Stop() error {
	a.grpcServer.GracefulStop()
	a.connPool.Close()
	err := a.tcpListener.Close()
	if err != nil {
		return err
	}

	return nil
}
