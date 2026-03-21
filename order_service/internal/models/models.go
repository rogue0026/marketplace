package models

import "time"

type ProductInfo struct {
	ProductId     uint64 `json:"product_id"`
	PriceSnapshot uint64 `json:"price_snapshot"`
}

type OrderInfo struct {
	Id           uint64         `json:"id"`
	UserId       uint64         `json:"user_id"`
	ProductsInfo []*ProductInfo `json:"products_info"`
	CreatedAt    time.Time      `json:"created_at"`
	TotalPrice   uint64         `json:"total_price"`
	Status       string         `json:"status"`
}
