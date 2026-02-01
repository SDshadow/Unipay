package alipay

import (
	"context"

	"Unipay/internal/core"
	"Unipay/internal/model"
)

// AlipayProvider 实现了 Provider 接口（简化版）
type AlipayProvider struct{}

func NewAlipayProvider() *AlipayProvider { return &AlipayProvider{} }

func (a *AlipayProvider) Pay(ctx context.Context, order model.Order) (core.Result, error) {
	r := core.NewRocket(map[string]any{"order": order})
	pipeline := core.NewPipeline(
		BuildPlugin{},
		SignPlugin{},
		HttpPlugin{},
		VerifyPlugin{},
		ParsePlugin{},
	)
	return pipeline.Execute(r)
}
