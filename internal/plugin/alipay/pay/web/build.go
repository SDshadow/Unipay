package web

import (
	"Unipay/internal/core"
	"context"
)

type BuildPlugin struct{}

func (p BuildPlugin) Handle(ctx context.Context, r *core.Rocket, next core.Next) error {
	r.MergePayload(map[string]any{
		"method": "alipay.trade.page.pay",
	})
	return next(ctx, r)
}
