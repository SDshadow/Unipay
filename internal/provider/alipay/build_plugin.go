package alipay

import (
	"Unipay/internal/core"
)

// BuildPlugin 构建支付请求参数（示例）
type BuildPlugin struct {
	provider *AlipayProvider
}

// Handle implements [core.Plugin].
func (b BuildPlugin) Handle(r *core.Rocket) error {
	panic("unimplemented")
}

func (b BuildPlugin) Web(rocket *core.Rocket) error {
	if rocket.Payload == nil {
		rocket.Payload = map[string]any{}
	}

	outTradeNo := rocket.Params["out_trade_no"].(string)
	amount := rocket.Params["total_amount"].(int64)

	payload := map[string]any{
		// "app_id":  b.provider.AppID,
		"app_id":       "9021000159670453",
		"method":       "alipay.trade.page.pay",
		"charset":      "utf-8",
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}

	payload["biz_content"] = map[string]any{
		"out_trade_no": outTradeNo,
		"total_amount": amount,
		"subject":      rocket.Params["subject"],
	}

	rocket.Payload["alipay"] = payload
	return nil
}
