package main

import (
	"fmt"
	"os"
	"os/signal"
	"user_service/application"
)

func main() {
	_ = os.Setenv("DATABASE_URL", "postgresql://user:password@localhost:5431/products")
	app, err := application.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.Run(":50052")
	if err != nil {
		fmt.Println(err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("user service started")
	<-stop
	app.Stop()
	fmt.Println("user service stopped")
}
