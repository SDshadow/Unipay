package model

// Callback 表示支付回调的简化模型
type Callback struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
