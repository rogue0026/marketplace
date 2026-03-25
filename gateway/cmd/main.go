package main

import (
	"context"
	"fmt"
	"gateway/application"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	app, err := application.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app.Run()
	fmt.Printf("starting api gateway")
	<-stop
	app.Stop(context.Background())
}
