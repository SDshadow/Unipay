package alipay

import (
	"Unipay/internal/core"
)

// ParsePlugin 将第三方响应解析为统一 Result
type ParsePlugin struct{}

func (p ParsePlugin) Handle(r *core.Rocket) error {
	// 将之前写入的模拟响应解析为 Result
	r.Result = core.Result{
		Success: true,
		Code:    "200",
		Message: "ok",
		Data:    r.Response,
	}
	return nil
}
