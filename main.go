package main

import (
	"Unipay/internal/core"
	"fmt"
)

func main() {
	// load config.json from repo root
	_ = core.LoadFromFile("config.json")
	config := core.Get("alipay", "default")
	fmt.Println(config["app_id"])
}
