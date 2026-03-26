package models

type Product struct {
	Id       uint64 `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    uint64 `json:"price,omitempty"`
	Quantity uint64 `json:"quantity,omitempty"`
}

type BasketInfo struct {
	TotalPrice uint64     `json:"total_price"`
	Products   []*Product `json:"products"`
}
