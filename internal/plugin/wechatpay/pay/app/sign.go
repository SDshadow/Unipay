package app

import (
	"Unipay/internal/core"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

type SignPlugin struct{}

func (s *SignPlugin) Assembly(ctx context.Context, r *core.Rocket, next core.Next) (*core.Rocket, error) {
	config := core.Get("wechatpay", "default")
	if config == nil {
		return nil, fmt.Errorf("wechatpay config not found")
	}

	privateKeyPEM, ok := config["private_key"].(string)
	if !ok {
		return nil, fmt.Errorf("private_key not found in config")
	}

	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}

	method := "POST"
	url := "/v3/pay/transactions/app"
	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()

	bodyBytes, _ := json.Marshal(r.Payload)
	body := string(bodyBytes)

	signMessage := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n",
		method,
		url,
		timestamp,
		nonceStr,
		body,
	)

	hash := sha256.Sum256([]byte(signMessage))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	serialNo := config["serial_no"].(string)
	mchid := config["mchid"].(string)
	authorization := fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%d",serial_no="%s"`,
		mchid,
		nonceStr,
		base64.StdEncoding.EncodeToString(signature),
		timestamp,
		serialNo,
	)

	r.MergePayload(map[string]any{
		"authorization":  authorization,
		"sign_timestamp": timestamp,
		"sign_nonce_str": nonceStr,
	})

	return next(ctx, r)
}

func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaKey, nil
}
