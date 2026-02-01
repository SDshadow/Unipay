package model

// PayRequest 是业务传入的支付请求
type PayRequest struct {
	Channel string `json:"channel"`
	Order   Order  `json:"order"`
}

// PayResult 是业务层对外的返回结构
type PayResult struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
