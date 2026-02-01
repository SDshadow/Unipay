package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"Unipay/internal/model"
	"Unipay/internal/service"
)

// PaymentHandler 简单的 HTTP 处理器，接收 JSON 并调用 PaymentService
func PaymentHandler(svc *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.PayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request"))
			return
		}
		res, err := svc.Pay(context.Background(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
