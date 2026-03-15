package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"product_service/application"
)

func main() {
	_ = os.Setenv("DATABASE_URL", "postgresql://user:password@localhost:5432/products")
	app, err := application.New()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = app.Run(":50051")
	if err != nil {
		fmt.Println(err.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	app.Stop()
	fmt.Println("programm terminated")
}
