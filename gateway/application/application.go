package application

import (
	"context"
	"fmt"
	"gateway/internal/clients"
	"gateway/internal/config"
	"gateway/internal/httphandlers"
	"gateway/internal/middleware"
	"gateway/internal/service"
	"gateway/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	httpServer *http.Server
}

func New() (*Application, error) {
	appLogger := logger.New()
	mw := middleware.LoggingMiddleware(appLogger)

	clientsConfig, err := config.LoadGRPCClientsConfig()
	if err != nil {
		return nil, err
	}

	httpServerConfig, err := config.LoadHttpServerConfig()
	if err != nil {
		return nil, err
	}

	c, err := clients.New(clientsConfig)
	if err != nil {
		return nil, err
	}

	svc := service.NewGatewayService(c)

	mux := chi.NewRouter()
	mux.Get("/products", httphandlers.ProductCatalogHandler(svc))
	mux.Post("/users", httphandlers.CreateUserHandler(svc))
	mux.Delete("/users", httphandlers.DeleteUserHandler(svc))
	mux.Get("/users/basket", httphandlers.BasketInfoHandler(svc))
	mux.Post("/users/basket", httphandlers.AddProductToBasketHandler(svc))

	s := &http.Server{
		Addr:    httpServerConfig.Addr,
		Handler: mw(mux),
	}

	app := &Application{
		httpServer: s,
	}

	return app, nil
}

func (a *Application) Run() {
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

}

func (a *Application) Stop(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
