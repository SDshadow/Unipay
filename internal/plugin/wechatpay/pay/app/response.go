package app

import (
	"Unipay/internal/core"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type ResponsePlugin struct{}

func (s *ResponsePlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	response, ok := r.Payload["response"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("response not found in payload")
	}

	if code, ok := response["code"]; ok {
		return nil, fmt.Errorf("wechatpay error: code=%v, message=%v", code, response["message"])
	}

	prepayID, ok := response["prepay_id"].(string)
	if !ok {
		return nil, fmt.Errorf("prepay_id not found in response")
	}

	config := core.Get("wechatpay", "default")
	if config == nil {
		return nil, fmt.Errorf("wechatpay config not found")
	}

	appid := config["appid"].(string)
	mchid := config["mchid"].(string)
	privateKeyPEM := config["private_key"].(string)

	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}

	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()
	packageValue := "Sign=WXPay"

	signMessage := fmt.Sprintf("%s\n%d\n%s\n%s\n",
		appid,
		timestamp,
		nonceStr,
		prepayID,
	)

	hash := sha256.Sum256([]byte(signMessage))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	appPayParams := map[string]any{
		"appid":     appid,
		"partnerid": mchid,
		"prepayid":  prepayID,
		"package":   packageValue,
		"noncestr":  nonceStr,
		"timestamp": timestamp,
		"sign":      base64.StdEncoding.EncodeToString(signature),
	}

	r.MergePayload(map[string]any{
		"prepay_id":      prepayID,
		"app_pay_params": appPayParams,
		"status":         "success",
	})

	return next(ctx, r)
}
