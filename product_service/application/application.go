package application

import (
	"context"
	"log"
	"net"
	"product_service/api/pb"
	"product_service/internal/service"
	"product_service/internal/storage/products/pg"
	"product_service/internal/transport/grpcapi"
	"product_service/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Application struct {
	grpcServer *grpc.Server
	connPool   *pgxpool.Pool
}

func New() (*Application, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.NewPool(ctx)
	if err != nil {
		return nil, err
	}
	productsRepo := pg.NewProductsRepository(pool)
	productService := service.NewProductService(productsRepo)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	router := grpcapi.NewHandler(productService)

	pb.RegisterProductServiceServer(grpcServer, router)
	reflection.Register(grpcServer)
	a := &Application{
		grpcServer: grpcServer,
		connPool:   pool,
	}

	return a, nil
}

func (a *Application) Run(grpcAddress string) error {
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("starting grpc server on %s\n", grpcAddress)
		err := a.grpcServer.Serve(listener)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	return nil
}

func (a *Application) Stop() {
	a.grpcServer.GracefulStop()
	a.connPool.Close()
}
