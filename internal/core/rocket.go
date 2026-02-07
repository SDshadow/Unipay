package core

import (
	"context"
	"net/http"
)

type Direction interface {
	Execute(ctx context.Context, r *Rocket) (*Rocket, error)
}

type Rocket struct {
	Params    map[string]any
	Payload   map[string]any
	Radar     *http.Request
	Direction Direction
}

func NewRocket(params map[string]any) *Rocket {
	return &Rocket{
		Params:  params,
		Payload: make(map[string]any),
	}
}

func (r *Rocket) MergePayload(payload map[string]any) {
	if r.Payload == nil {
		r.Payload = make(map[string]any)
	}
	for k, v := range payload {
		r.Payload[k] = v
	}
}
