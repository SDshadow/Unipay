package core

// Rocket 是支付执行过程的上下文载体，插件之间通过它传递数据。
type Rocket struct {
	Params   map[string]any // 原始业务参数
	Payload  map[string]any // 支付协议参数（逐步构建）
	Response any            // 第三方返回结果
	Result   Result         // 统一支付结果
}

func NewRocket(params map[string]any) *Rocket {
	return &Rocket{
		Params:  params,
		Payload: map[string]any{},
		Result:  Result{},
	}
}
