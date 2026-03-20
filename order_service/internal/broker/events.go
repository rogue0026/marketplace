package broker

type EventOrderCreated struct {
	OrderId uint64 `json:"order_id"`
	UserId  uint64 `json:"user_id"`
}

func (e *EventOrderCreated) GetEventName() string {
	return "order.created"
}

func (e *EventOrderCreated) GetTopicName() string {
	return "order.created"
}
