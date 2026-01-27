package payment

type OrderRepository interface {
	Save(order *Order) error
	Get(orderID string) (*Order, error)
}
