package core

import (
	"context"
)

// Plugin 表示支付流程中的一个独立步骤
type Plugin interface {
	Assembly(ctx context.Context, r *Rocket, next Next) (*Rocket, error)
}

type Next func(ctx context.Context, r *Rocket) (*Rocket, error)
