package main

import (
	"log"
	"net/http"

	"Unipay/api"
	"Unipay/internal/provider"
	"Unipay/internal/provider/alipay"
	"Unipay/internal/service"
)

func main() {
	factory := provider.NewProviderFactory()
	factory.Register("alipay", alipay.NewAlipayProvider())

	svc := service.NewPaymentService(factory)
	r := api.NewRouter(svc)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
