package provider

import (
	"context"
	"sync"

)

type Params map[string]any

type Provider interface {
	Name() string
}
type PayProvider interface {
	Provider
	Pay(ctx context.Context, params Params) error
}
type RefundProvider interface {
	Provider
	Refund(ctx context.Context, params Params) error
}

type ProviderFactory struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]Provider),
	}
}
func (f *ProviderFactory) Register(p Provider) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.providers[p.Name()] = p
}
func (f *ProviderFactory) Get(name string) (Provider, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	p, ok := f.providers[name]
	return p, ok
}
func (f *ProviderFactory) Pay(name string) (PayProvider, bool) {
	p, ok := f.Get(name)
	if !ok {
		return nil, false
	}
	pp, ok := p.(PayProvider)
	return pp, ok
}
func (f *ProviderFactory) Refund(name string) (RefundProvider, bool) {
	p, ok := f.Get(name)
	if !ok {
		return nil, false
	}
	rp, ok := p.(RefundProvider)
	return rp, ok
}
