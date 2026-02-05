package web

import (
	"Unipay/internal/core"
)

type BuildPlugin struct{}

func (p BuildPlugin) Handle(r *core.Rocket) error
