package core

// Plugin 表示支付流程中的一个独立步骤
type Plugin interface {
	Handle(r *Rocket) error
}
