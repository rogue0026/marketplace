package models

type Product struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	PricePerUnit  uint64 `json:"price_per_unit"`
	TotalQuantity uint64 `json:"total_quantity"`
}
