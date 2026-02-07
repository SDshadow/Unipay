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
			"out_trade_no": r.Params["out_trade_no"],
			"total_amount": r.Params["total_amount"],
			"subject":      r.Params["subject"],
		},
	})
	return next(ctx, r)
}
