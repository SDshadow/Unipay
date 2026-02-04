package alipay

import (
	"Unipay/internal/core"
)

// HttpPlugin 发起 HTTP 请求（示例 stub，会写入模拟响应）
type HttpPlugin struct{}

// Handle implements [core.Plugin].
func (h HttpPlugin) Handle(r *core.Rocket) error {
	panic("unimplemented")
}

func (h HttpPlugin) Web(r *core.Rocket) error {
	return nil
}
