package alipay

import (
	"Unipay/internal/core"
)

// ParsePlugin 将第三方响应解析为统一 Result
type ParsePlugin struct{}

// Handle implements [core.Plugin].
func (p ParsePlugin) Handle(r *core.Rocket) error {
	panic("unimplemented")
}

func (p ParsePlugin) Web(r *core.Rocket) error {
	// 将之前写入的模拟响应解析为 Result
	return nil
}
