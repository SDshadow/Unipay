package alipay

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
	config := core.Get("alipay", "default")

	return map[string]any{
		"app_id":      config["app_id"],
		"method":      "",
		"format":      "JSON",
		"return_url":  config["return_url"],
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"sign":        "",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": map[string]any{},
	}
}
