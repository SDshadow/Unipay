package main

import (
	"Unipay/internal/core"
	alipay "Unipay/internal/plugin/alipay/pay"
	alipayweb "Unipay/internal/plugin/alipay/pay/web"
	"context"
	"fmt"
)

func main() {
	plugin := core.NewPipeline(
		&alipay.StartPlugin{},
		&alipayweb.BuildPlugin{})
	rocket := core.NewRocket(map[string]any{
		"out_trade_no": "20230618001",
		"total_amount": "88.88",
	})
	result, _ := plugin.Execute(context.Background(), rocket)
	fmt.Println(result.Payload)
}
