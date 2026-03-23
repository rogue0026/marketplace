package models

type Product struct {
	Id             uint64 `json:"id"`
	Name           string `json:"name"`
	PricePerUnit   uint64 `json:"price_per_unit"`
	StockRemaining uint64 `json:"total_quantity"`
}

type Reservation struct {
	Id        uint64 `json:"id"`
	OrderId   uint64 `json:"order_id"`
	ProductId uint64 `json:"product_id"`
	Quantity  uint64 `json:"quantity"`
}
