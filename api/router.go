package api

import (
	"net/http"

	"Unipay/api/handler"
	"Unipay/internal/service"
)

// NewRouter 创建并返回一个简单的路由器
func NewRouter(svc *service.PaymentService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/pay", handler.PaymentHandler(svc))
	return mux
}
