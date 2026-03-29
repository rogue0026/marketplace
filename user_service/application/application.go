package application

import (
	"context"
	"fmt"
	"net"
	"user_service/internal/config"
	"user_service/internal/service"
	"user_service/internal/storage/pg"
	"user_service/internal/transport/grpcapi"
	"user_service/internal/transport/interceptors"
	"user_service/pkg/logger"
	"user_service/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/rogue0026/marketplace-proto/users"

	"google.golang.org/grpc"
)

type Application struct {
	Cfg         *config.AppConfig
	connPool    *pgxpool.Pool
	tcpListener net.Listener
	grpcServer  *grpc.Server
}

func New() (*Application, error) {
	appConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	pool, err := postgres.NewPool(context.Background(), appConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}

	users := pg.NewUsersRepo(pool)
	wallets := pg.NewWalletsRepo(pool)
	baskets := pg.NewBasketsRepo(pool)

	userService := service.NewUserService(users, wallets, baskets)

	h := grpcapi.NewHandler(userService)

	appLogger := logger.New()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.Logging(appLogger)),
	)

	pb.RegisterUserServiceServer(grpcServer, h)

	l, err := net.Listen("tcp", appConfig.GRPCAddr)
	if err != nil {
		return nil, err
	}

	app := &Application{
		Cfg:         appConfig,
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
