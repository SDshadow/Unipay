package core

import (
	"context"
)

type Pipeline struct {
	plugins []Plugin
}

func NewPipeline(plugins ...Plugin) *Pipeline {
	return &Pipeline{plugins: plugins}
}

func (p *Pipeline) Execute(ctx context.Context, r *Rocket) (*Rocket, error) {
	var chain Next

	index := 0
	chain = func(ctx context.Context, r *Rocket) (*Rocket, error) {
		if index >= len(p.plugins) {
			return r, nil
		}

		plugin := p.plugins[index]
		index++

		return plugin.Assembly(ctx, r, chain)
	}

	return chain(ctx, r)
}
