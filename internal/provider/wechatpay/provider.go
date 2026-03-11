package wechatpay

import (
	"Unipay/internal/core"
	"context"
	"fmt"
)

type Provider struct {
	client *WechatPayClient
}

func NewProvider(config map[string]any) (*Provider, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	appid, ok := config["appid"].(string)
	if !ok {
		return nil, fmt.Errorf("appid is required")
	}

	mchid, ok := config["mchid"].(string)
	if !ok {
		return nil, fmt.Errorf("mchid is required")
	}

	serialNo, ok := config["serial_no"].(string)
	if !ok {
		return nil, fmt.Errorf("serial_no is required")
	}

	privateKey, ok := config["private_key"].(string)
	if !ok {
		return nil, fmt.Errorf("private_key is required")
	}

	apiV3Key, ok := config["api_v3_key"].(string)
	if !ok {
		return nil, fmt.Errorf("api_v3_key is required")
	}

	client, err := NewWechatPayClient(appid, mchid, serialNo, privateKey, apiV3Key)
	if err != nil {
		return nil, fmt.Errorf("create wechatpay client failed: %w", err)
	}

	return &Provider{
		client: client,
	}, nil
}

// APP支付
func (p *Provider) AppPay(ctx context.Context, rocket *core.Rocket) (*core.Rocket, error) {
	payload := rocket.Payload

	// 调用统一下单API
	prepayResp, err := p.client.AppPay(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("app pay failed: %w", err)
	}

	// 生成APP调起支付的参数
	appPayParams, err := p.client.GenerateAppPayParams(prepayResp.PrepayID)
	if err != nil {
		return nil, fmt.Errorf("generate app pay params failed: %w", err)
	}

	// 合并结果到Payload
	rocket.MergePayload(map[string]any{
		"prepay_id":      prepayResp.PrepayID,
		"app_pay_params": appPayParams,
	})

	return rocket, nil
}

// 查询订单
func (p *Provider) QueryOrder(ctx context.Context, rocket *core.Rocket) (*core.Rocket, error) {
	outTradeNo, ok := rocket.Params["out_trade_no"].(string)
	if !ok {
		return nil, fmt.Errorf("out_trade_no is required")
	}

	result, err := p.client.QueryOrder(ctx, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("query order failed: %w", err)
	}

	rocket.MergePayload(result)
	return rocket, nil
}

// 关闭订单
func (p *Provider) CloseOrder(ctx context.Context, rocket *core.Rocket) (*core.Rocket, error) {
	outTradeNo, ok := rocket.Params["out_trade_no"].(string)
	if !ok {
		return nil, fmt.Errorf("out_trade_no is required")
	}

	err := p.client.CloseOrder(ctx, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("close order failed: %w", err)
	}

	rocket.MergePayload(map[string]any{
		"status": "closed",
	})

	return rocket, nil
}

// 实现Plugin接口（如果需要作为插件使用）
func (p *Provider) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	// 根据Params中的action决定执行哪个操作
	action, ok := r.Params["action"].(string)
	if !ok {
		action = "pay" // 默认支付
	}

	var err error
	switch action {
	case "pay":
		r, err = p.AppPay(ctx, r)
	case "query":
		r, err = p.QueryOrder(ctx, r)
	case "close":
		r, err = p.CloseOrder(ctx, r)
	default:
		err = fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		return nil, err
	}

	return next(ctx, r)
}
