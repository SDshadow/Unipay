package service

import (
	"context"

	"Unipay/internal/model"
	"Unipay/internal/provider"
)

// PaymentService 是业务层的边界，负责选择 Provider 并发起支付
type PaymentService struct {
	factory *provider.ProviderFactory
}

func NewPaymentService(f *provider.ProviderFactory) *PaymentService {
	return &PaymentService{factory: f}
}

func (s *PaymentService) Pay(ctx context.Context, req model.PayRequest) (model.PayResult, error) {
	p := s.factory.Get(req.Channel)
	if p == nil {
		return model.PayResult{}, ErrProviderNotFound
	}
	res, err := p.Pay(ctx, req.Order)
	if err != nil {
		return model.PayResult{}, err
	}
	return model.PayResult{Success: res.Success, Code: res.Code, Message: res.Message, Data: res.Data}, nil
}

// ErrProviderNotFound 是简化的错误变量
var ErrProviderNotFound = &ProviderError{"provider not found"}

type ProviderError struct{ s string }

func (e *ProviderError) Error() string { return e.s }
