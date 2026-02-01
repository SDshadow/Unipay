package alipay

import (
	"Unipay/internal/core"
)

// BuildPlugin 构建支付请求参数（示例）
type BuildPlugin struct{}

func (b BuildPlugin) Handle(r *core.Rocket) error {
	// 从业务参数组装协议参数，这里只做示例
	if r.Payload == nil {
		r.Payload = map[string]any{}
	}
	if order, ok := r.Params["order"]; ok {
		r.Payload["order_id"] = order
	}
	return nil
}
