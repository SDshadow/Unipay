package channel

import "Unipay/src/domain/payment"

type PayResp struct {
	PayURL string
}

type CallbackData struct {
	OrderID string
	Success bool
}

type Channel interface {
	Pay(order *payment.Order) (*PayResp, error)
	ParseCallback(raw []byte) (*CallbackData, error)
}
