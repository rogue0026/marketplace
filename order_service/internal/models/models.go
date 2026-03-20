package models

import "time"

type Order struct {
	Id        uint64    `json:"id"`
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

type OrdersContent struct {
	OrderId       uint64 `json:"order_id"`
	ProductId     uint64 `json:"product_id"`
	PriceSnapshot uint64 `json:"price_snapshot"`
}
