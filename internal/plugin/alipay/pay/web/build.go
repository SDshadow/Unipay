package web

import (
	"Unipay/internal/core"
	"context"
)

type BuildPlugin struct{}

func (p BuildPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	r.MergePayload(map[string]any{
		"method": "alipay.trade.page.pay",
		"biz_content": map[string]any{
			"product_code": "FAST_INSTANT_TRADE_PAY",
		},
	})
	return next(ctx, r)
}
