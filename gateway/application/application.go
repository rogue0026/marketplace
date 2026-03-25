package application

import (
	"context"
	"fmt"
	"gateway/internal/clients"
	"gateway/internal/config"
	"net/http"
)

type Application struct {
	C          *clients.Clients
	httpServer *http.Server
}

func New() (*Application, error) {

	clientsConfig, err := config.LoadGRPCClientsConfig()
	if err != nil {
		return nil, err
	}

	httpServerConfig, err := config.LoadHttpServerConfig()
	if err != nil {
		return nil, err
	}

	clients, err := clients.New(clientsConfig)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s := &http.Server{
		Addr:    httpServerConfig.Addr,
		Handler: mux,
	}

	app := &Application{
		C:          clients,
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
