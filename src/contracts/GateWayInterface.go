package contracts

type PaymentResult struct {
	Success bool
	Code    string
	Message string
	Data    map[string]interface{}
}

type OrderInfo struct {
	TradeNo    string
	OutTradeNo string
	Amount     float64
	Status     string
	CreatedAt  string
	PaidAt     string
	Extra      map[string]interface{}
}

type VerifyResult struct {
	Valid   bool
	OrderID string
	Amount  float64
	Data    map[string]interface{}
}

type GateWayInterface interface {
	Pay(configBiz map[string]interface{}) (*PaymentResult, error)

	Refund(configBiz map[string]interface{}) (*PaymentResult, error)

	Close(outTradeNo string) (bool, error)

	Find(outTradeNo string) (*OrderInfo, error)

	Verify(data map[string]interface{}) (*VerifyResult, error)
}

type SyncVerifier interface {
	SyncVerify(data map[string]interface{}, sign string) (*VerifyResult, error)
}
