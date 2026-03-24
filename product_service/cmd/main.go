package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"product_service/application"
)

func main() {
	app, err := application.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app.Run()
	log.Printf("starting product service at %s", app.Cfg.GRPCAddr)
	<-stop

	err = app.Stop()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("product service stopped")
}
