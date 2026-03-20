package main

import (
	"fmt"
	"os"
	"os/signal"
	"user_service/application"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	appCfg, err := application.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	app, err := application.New(appCfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.Run(app.Cfg.GRPCServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user service started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	app.Stop()
	fmt.Println("user service stopped")
}
