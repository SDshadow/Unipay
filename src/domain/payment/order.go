package payment

type Status string

const (
	StatusCreated Status = "CREATED"
	StatusPaying  Status = "PAYING"
	StatusPaid    Status = "PAID"
	StatusFailed  Status = "FAILED"
)

type Money struct {
	Amount   int64
	Currency string
}

type Order struct {
	ID      string
	Amount  Money
	Channel string
	Status  Status
}

func (o *Order) CanPay() bool {
	return o.Status == StatusCreated
}

func (o *Order) MarkPaying() {
	o.Status = StatusPaying
}

func (o *Order) MarkPaid() {
	o.Status = StatusPaid
}
