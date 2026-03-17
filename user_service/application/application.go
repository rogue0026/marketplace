package application

import (
	"context"
	"log"
	"net"
	"user_service/api/pb"
	"user_service/internal/service"
	"user_service/internal/storage/pg"
	"user_service/internal/transport/grpcapi"
	"user_service/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Application struct {
	grpcServer *grpc.Server
	connPool   *pgxpool.Pool
}

func New() (*Application, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	p, err := postgres.NewPool(ctx)
	if err != nil {
		return nil, err
	}

	wallets := pg.NewWalletsRepo(p)
	baskets := pg.NewBasketsRepo(p)
	users := pg.NewUsersRepo(p)

	userService := service.NewUserService(
		users,
		wallets,
		baskets,
	)

	h := grpcapi.NewHandler(userService)

	pb.RegisterUserServiceServer(grpcServer, h)

	app := &Application{
		grpcServer: grpcServer,
		connPool:   p,
	}

	return app, nil
}

func (a *Application) Run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		err := a.grpcServer.Serve(listener)
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (a *Application) Stop() {
	a.grpcServer.GracefulStop()
	a.connPool.Close()
}
