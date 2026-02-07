package web

import (
	"Unipay/internal/core"
	"context"
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
	return buf.String()
}

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
