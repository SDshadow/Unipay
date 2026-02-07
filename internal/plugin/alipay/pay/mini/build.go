package mini

import (
	"Unipay/internal/core"
	"context"
)

type BuildPlugin struct{}

func (p BuildPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	r.MergePayload(map[string]any{
		"method": "alipay.trade.create",
		"biz_content": map[string]any{
			"product_code": "JSAPI_PAY",
			"out_trade_no": r.Params["out_trade_no"],
			"total_amount": r.Params["total_amount"],
			"subject":      r.Params["subject"],
		},
	})
	return next(ctx, r)
}
