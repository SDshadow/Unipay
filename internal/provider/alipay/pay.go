package alipay

import (
	"context"

	"Unipay/internal/core"
)

// AlipayProvider 实现了 Provider 接口（简化版）
type AlipayProvider struct{}

func NewAlipayProvider() *AlipayProvider { return &AlipayProvider{} }

// func (a *AlipayProvider) Webpay(ctx context.Context, rocket *core.Rocket) (core.Result, error) {

// 	rocket.MergePayload(getPayload(rocket.Params))
// 	// r := core.NewRocket(map[string]any{"order": order})
// 	pipeline := core.NewPipeline(
// 		BuildPlugin{},
// 		SignPlugin{},
// 		HttpPlugin{},
// 		VerifyPlugin{},
// 		ParsePlugin{},
// 	)
// 	return pipeline.Execute(rocket)
// }

func (a *AlipayProvider) Webpay(ctx context.Context, rocket *core.Rocket) (core.Result, error) {
	rocket.MergePayload(getPayload(rocket.Params))

}

func getPayload(params map[string]any) map[string]any {
	return map[string]any{
		// "app_id":      "9021000159670453",
		"app_id":      params["app_id"],
		"method":      "",
		"format":      "JSON",
		"return_url":  "https://example.com/return",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"sign":        "",
		"timestamp":   "2024-06-01 12:00:00",
		"version":     "1.0",
		"biz_content": map[string]any{},
	}
}
