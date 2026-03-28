package application

import (
	"context"
	"fmt"
	"net"
	"product_service/internal/config"
	"product_service/internal/service"
	"product_service/internal/storage/pg"
	"product_service/internal/transport/interceptors"
	"product_service/pkg/logger"
	"product_service/pkg/postgres"

	"product_service/internal/transport/grpcapi"

	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/rogue0026/marketplace-proto/products"
	"google.golang.org/grpc"
)

type Application struct {
	Cfg         *config.AppConfig
	connPool    *pgxpool.Pool
	s           *service.ProductService
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

	products := pg.NewProductsRepo(pool)
	userService := service.NewProductService(products)

	appLogger := logger.New()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.Logging(appLogger)),
	)

	h := grpcapi.NewProductHandler(userService)

	pb.RegisterProductServiceServer(grpcServer, h)

	l, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return nil, err
	}
	app := &Application{
		Cfg:         cfg,
		connPool:    pool,
		tcpListener: l,
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
