package main

import (
	"context"
	"fmt"
	"user_service/internal/models"
	"user_service/internal/storage/pg"
	"user_service/pkg/postgres"
)

func main() {
	conn := "postgresql://user:password@localhost:5431/user_service_db"
	p, err := postgres.NewPool(context.Background(), conn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	baskets := pg.NewBasketsRepo(p)

	product := &models.Product{
		UserId:          uint64(1),
		ProductId:       2,
		ProductQuantity: 5,
		PricePerUnit:    1800,
	}
	err = baskets.AddProductToBasket(context.Background(), product)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
