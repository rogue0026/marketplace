package application

import (
	"context"
	"fmt"
	"net"
	"product_service/api/pb"
	"product_service/internal/service"
	"product_service/internal/storage/pg"
	"product_service/internal/transport/grpcapi"
	"product_service/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Application struct {
	grpcServer *grpc.Server
	connPool   *pgxpool.Pool
	Cfg        *Config
}

func New(cfg *Config) (*Application, error) {
	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	products := pg.NewProductsRepo(pool)
	productService := service.NewProductService(products)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	h := grpcapi.NewHandler(productService)

	pb.RegisterProductServiceServer(grpcServer, h)

	a := &Application{
		grpcServer: grpcServer,
		connPool:   pool,
		Cfg:        cfg,
	}

	return a, nil
}

func (a *Application) Run() error {
	listener, err := net.Listen("tcp", a.Cfg.GRPCServerAddr)
	if err != nil {
		return err
	}

	go func() {
		err := a.grpcServer.Serve(listener)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	return nil
}

func (a *Application) Stop() {
	a.grpcServer.GracefulStop()
	a.connPool.Close()
}
