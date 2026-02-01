package provider

import (
	"context"
	"sync"

	"Unipay/internal/core"
	"Unipay/internal/model"
)

// Provider 负责为某个支付渠道组装 Plugin Pipeline 并执行
type Provider interface {
	Pay(ctx context.Context, order model.Order) (core.Result, error)
}

// ProviderFactory 简单工厂，用于根据渠道获取 Provider
type ProviderFactory struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{providers: map[string]Provider{}}
}

func (f *ProviderFactory) Register(name string, p Provider) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.providers[name] = p
}

func (f *ProviderFactory) Get(name string) Provider {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.providers[name]
}
