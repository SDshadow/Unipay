package tests

import (
	"context"

	"Unipay/internal/core"
	"Unipay/internal/model"
)

// MockProvider 简化的测试 Provider
type MockProvider struct{}

func (m *MockProvider) Pay(ctx context.Context, order model.Order) (core.Result, error) {
	return core.Result{Success: true, Code: "200", Message: "mock ok", Data: map[string]any{"order_id": order.ID}}, nil
}
