package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"user_service/application"
)

func main() {
	app, err := application.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	app.Run()
	log.Printf("starting user service at %s\n", app.Cfg.GRPCAddr)
	<-stop

	app.Stop()

	log.Println("user service stopped")

}
