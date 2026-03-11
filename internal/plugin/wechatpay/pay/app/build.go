package app

import (
	"Unipay/internal/core"
	"context"
	"time"
)

type BuildPlugin struct{}

func (p BuildPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	r.MergePayload(getPayload(r.Params))
	return next(ctx, r)
}

func getPayload(params map[string]any) map[string]any {
	_ = core.LoadFromFile("config.json")
	config := core.Get("wechatpay", "default")

	return map[string]any{
		"appid":        config["appid"],
		"mchid":        config["mchid"],
		"description":  params["description"],
		"out_trade_no": params["out_trade_no"],
		"notify_url":   config["notify_url"],
		"amount": map[string]any{
			"total":    params["amount"], // 单位：分
			"currency": "CNY",
		},
		"scene_info": map[string]any{
			"payer_client_ip": params["payer_client_ip"],
		},
		"timestamp": time.Now().Unix(),
		"nonce_str": generateNonceStr(),
	}
}

func generateNonceStr() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
