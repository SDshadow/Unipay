package model

// Order 表示一个最简化的订单结构
type Order struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}
