package models

import "time"

const (
	RESERVATION_STATUS_REQUIRES_PAYMENT = iota
	RESERVATION_STATUS_PAID_FOR
)

type Product struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	PricePerUnit  uint64 `json:"price_per_unit"`
	TotalQuantity uint64 `json:"total_quantity"`
}

type Reservation struct {
	Id        uint64    `json:"id"`
	OrderId   uint64    `json:"order_id"`
	ProductId uint64    `json:"product_id"`
	Quantity  uint64    `json:"quantity"`
	Status    int       `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}
