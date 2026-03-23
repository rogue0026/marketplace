package models

type OrderItem struct {
}

type OrderData struct {
	Id         uint64       `json:"order_id"`
	UserId     uint64       `json:"user_id"`
	Items      []*OrderItem `json:"items"`
	TotalPrice uint64       `json:"total_price"`
}
