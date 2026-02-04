package alipay

import (
	"Unipay/internal/core"
)

// SignPlugin 对请求进行签名（示例 stub）
type SignPlugin struct{}

// Handle implements [core.Plugin].
func (s SignPlugin) Handle(r *core.Rocket) error {
	panic("unimplemented")
}

func (s SignPlugin) Web(r *core.Rocket) error {
	// 简化的签名逻辑：设置一个 dummy 字段
	r.Payload["signature"] = "signed-by-alipay"
	return nil
}
