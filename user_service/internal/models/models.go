package models

type User struct {
	Id           uint64 `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type Product struct {
	Id       uint64 `json:"id"`
	Quantity uint64 `json:"quantity"`
}

type Wallet struct {
	Id      uint64 `json:"id"`
	UserId  uint64 `json:"user_id"`
	Balance uint64 `json:"balance"`
}

type UserBasketInfo struct {
	Products []*Product `json:"products"`
}
