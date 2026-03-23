package models

type User struct {
	Id           uint64 `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type Product struct {
	Id              uint64 `json:"id"`
	UserId          uint64 `json:"user_id"`
	ProductId       uint64 `json:"product_id"`
	ProductQuantity uint64 `json:"product_quantity"`
	PricePerUnit    uint64 `json:"price_per_unit"`
}

type Wallet struct {
	Id      uint64 `json:"id"`
	UserId  uint64 `json:"user_id"`
	Balance uint64 `json:"balance"`
}

type UserBasketInfo struct {
	Products   []*Product `json:"products"`
	TotalPrice uint64     `json:"total_price"`
}
