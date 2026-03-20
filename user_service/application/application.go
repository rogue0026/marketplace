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
	Cfg        *Config
}

func New(cfg *Config) (*Application, error) {
	ctx := context.Background()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	p, err := postgres.NewPool(ctx, cfg.DatabaseURL)
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
		Cfg:        cfg,
	}

	return app, nil
}

func (a *Application) Run() error {
	listener, err := net.Listen("tcp", a.Cfg.GRPCServerAddr)
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
