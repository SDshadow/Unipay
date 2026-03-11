package app

import (
	"Unipay/internal/core"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RequestPlugin struct{}

func (s *RequestPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	authorization, ok := r.Payload["authorization"].(string)
	if !ok {
		return nil, fmt.Errorf("authorization not found in payload")
	}

	bodyBytes, err := json.Marshal(r.Payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}

	url := "https://api.mch.weixin.qq.com/v3/pay/transactions/app"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("User-Agent", "Unipay/1.0")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	r.MergePayload(map[string]any{
		"response":      result,
		"response_code": resp.StatusCode,
	})

	return next(ctx, r)
}
