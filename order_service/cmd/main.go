package main

import (
	"log"
	"order_service/internal/application"
	"os"
	"os/signal"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Println(err.Error())
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app.Run()
	log.Printf("starting order service at %s\n", app.Cfg.GRPCAddr)
	<-stop

	log.Println("order service stopped")
}
