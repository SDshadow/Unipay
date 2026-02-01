package alipay

import (
	"Unipay/internal/core"
)

// HttpPlugin 发起 HTTP 请求（示例 stub，会写入模拟响应）
type HttpPlugin struct{}

func (h HttpPlugin) Handle(r *core.Rocket) error {
	// 这里不调用真实网络，写入一个模拟的响应对象
	r.Response = map[string]any{"status": "ok", "raw": "alipay-response"}
	return nil
}
