package application

type Application struct{}

func New() (*Application, error) {

	app := &Application{}

	return app, nil
}
