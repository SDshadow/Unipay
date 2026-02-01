package tests

import (
	"context"
	"testing"

	"Unipay/internal/model"
	"Unipay/internal/provider"
	"Unipay/internal/service"
)

func TestPaymentFlow(t *testing.T) {
	factory := provider.NewProviderFactory()
	factory.Register("mock", &MockProvider{})

	svc := service.NewPaymentService(factory)

	req := model.PayRequest{Channel: "mock", Order: model.Order{ID: "o-123", Amount: 1000}}
	res, err := svc.Pay(context.Background(), req)
	if err != nil {
		t.Fatalf("pay failed: %v", err)
	}
	if !res.Success {
		t.Fatalf("unexpected success=false")
	}
}
