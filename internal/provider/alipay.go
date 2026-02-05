package provider

import (
	"Unipay/internal/core"
	"context"
)

type AlipayProvider struct{}

func (a *AlipayProvider) Name() string {
	return "alipay"
}

func (a *AlipayProvider) Pay(ctx context.Context, params Params) (core.Result, error) {
	rocket := core.NewRocket(params)
}
