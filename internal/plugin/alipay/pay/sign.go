package alipay

import (
	"Unipay/internal/core"
	"context"
	"encoding/base64"
	"sort"
	"strings"
)

type SignPlugin struct{}

func (p SignPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {

	sign := sign(r)
	r.MergePayload(map[string]any{
		"sign": sign,
	})
	return next(ctx, r)
}

func sign(r *core.Rocket) string {
	keys := prepareSignString(r)
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(r.Payload[k].(string))
	}
	return base64.StdEncoding.EncodeToString([]byte(buf.String()))
}

// 过滤掉空值和sign参数，并对剩余参数进行排序
func prepareSignString(r *core.Rocket) []string {
	filtered := make([]string, 0)
	for k, v := range r.Payload {
		if v != "" && k != "sign" {
			filtered = append(filtered, k)
		}
	}
	sort.Strings(filtered)
	return filtered
}

// func getPrivateKey() *rsa.PrivateKey {
// 	_ = core.LoadFromFile("config.json")
// 	config := core.Get("alipay", "default")
// 	if config["private_key"] == nil {
// 		// zap   log.Error("private_key is missing in config")
// 		return nil
// 	}
// 	privateKeyStr := config["private_key"].(string)
// 	privateKey, err := core.ParsePrivateKey(privateKeyStr)
// }
