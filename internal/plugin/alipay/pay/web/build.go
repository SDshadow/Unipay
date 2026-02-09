package web

import (
	"Unipay/internal/core"
	"context"
	"encoding/json"
	"maps"
)

type BuildPlugin struct{}

func (p BuildPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	static_content := map[string]any{
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}
	maps.Copy(r.Payload["biz_content"].(map[string]any), static_content)
	maps.Copy(r.Payload["biz_content"].(map[string]any), r.Params)
	marshal, err := json.Marshal(r.Payload["biz_content"])
	if err != nil {
		return r, err
	}
	r.Payload["biz_content"] = string(marshal)
	return next(ctx, r)
}
