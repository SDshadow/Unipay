package wechatpay

import (
	"Unipay/internal/core"
	"context"
	"time"
)

type StartPlugin struct{}

func (s *StartPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	r.MergePayload(getPayload(r.Params))
	r, err := next(ctx, r)
	return r, err
}

func getPayload(params map[string]any) map[string]any {
	_ = core.LoadFromFile("config.json")
	config := core.Get("wechatpay", "default")

	// 生成随机字符串
	nonceStr := generateNonceStr()
	timestamp := time.Now().Unix()

	// APP支付请求参数
	payload := map[string]any{
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
		"nonce_str": nonceStr,
		"timestamp": timestamp,
		"sign_type": "RSA",
	}

	// 如果有用户openid（APP支付通常不需要，但保留字段）
	if openid, ok := params["openid"]; ok {
		payload["payer"] = map[string]any{
			"openid": openid,
		}
	}

	return payload
}

func generateNonceStr() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
