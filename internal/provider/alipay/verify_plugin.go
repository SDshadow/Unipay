package alipay

import (
	"Unipay/internal/core"
	"errors"
)

// VerifyPlugin 验证返回签名（简化）
type VerifyPlugin struct{}

// Handle implements [core.Plugin].
func (v VerifyPlugin) Handle(r *core.Rocket) error {
	panic("unimplemented")
}

func (v VerifyPlugin) Web(r *core.Rocket) error {
	// 简化：如果响应中没有 status 或不是 ok 则认为验签/校验失败
	if resp, ok := r.Response.(map[string]any); ok {
		if resp["status"] == "ok" {
			return nil
		}
	}
	return errors.New("verify failed")
}
