package models

type OrderItem struct {
	ProductId    uint64 `json:"product_id"`
	Quantity     uint64 `json:"quantity"`
	PricePerUnit uint64 `json:"price_per_unit"`
}

type Order struct {
	OrderId    uint64       `json:"order_id"`
	UserId     uint64       `json:"user_id"`
	Items      []*OrderItem `json:"items"`
	TotalPrice uint64       `json:"total_price"`
}
