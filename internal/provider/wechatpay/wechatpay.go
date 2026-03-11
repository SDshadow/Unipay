package wechatpay

import (
	"bytes"
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
	"io"
	"net/http"
	"strings"
	"time"
)

type WechatPayClient struct {
	appid      string
	mchid      string
	serialNo   string
	privateKey *rsa.PrivateKey
	apiV3Key   string
	httpClient *http.Client
	baseURL    string
}

// 微信支付响应结构
type WechatPayResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// APP支付预支付响应
type PrepayResponse struct {
	PrepayID string `json:"prepay_id"`
}

func NewWechatPayClient(appid, mchid, serialNo, privateKeyPEM, apiV3Key string) (*WechatPayClient, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}

	return &WechatPayClient{
		appid:      appid,
		mchid:      mchid,
		serialNo:   serialNo,
		privateKey: privateKey,
		apiV3Key:   apiV3Key,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://api.mch.weixin.qq.com",
	}, nil
}

// 解析私钥
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

// 生成签名
func (c *WechatPayClient) generateSignature(method, url, body string) (string, int64, string, error) {
	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()

	// 构建签名字符串
	// 格式：HTTP请求方法\nURL\n请求时间戳\n请求随机串\n请求报文主体\n
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n",
		method,
		strings.TrimPrefix(url, "https://api.mch.weixin.qq.com"),
		timestamp,
		nonceStr,
		body,
	)

	// 使用私钥签名
	hash := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", 0, "", fmt.Errorf("sign failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), timestamp, nonceStr, nil
}

// 生成Authorization头
func (c *WechatPayClient) generateAuthorization(method, url, body string) (string, error) {
	signature, timestamp, nonceStr, err := c.generateSignature(method, url, body)
	if err != nil {
		return "", err
	}

	// 格式：WECHATPAY2-SHA256-RSA2048 mchid="商户号",nonce_str="随机串",signature="签名",timestamp="时间戳",serial_no="证书序列号"
	auth := fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%d",serial_no="%s"`,
		c.mchid,
		nonceStr,
		signature,
		timestamp,
		c.serialNo,
	)

	return auth, nil
}

// 调用微信支付API
func (c *WechatPayClient) CallAPI(ctx context.Context, method, endpoint string, body map[string]any) (map[string]any, error) {
	url := c.baseURL + endpoint

	var jsonBody []byte
	var err error
	if body != nil && len(body) > 0 {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body failed: %w", err)
		}
	} else {
		jsonBody = []byte{}
	}

	auth, err := c.generateAuthorization(method, url, string(jsonBody))
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if len(jsonBody) > 0 {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", auth)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return result, nil
}

// APP支付统一下单
func (c *WechatPayClient) AppPay(ctx context.Context, params map[string]any) (*PrepayResponse, error) {
	endpoint := "/v3/pay/transactions/app"
	result, err := c.CallAPI(ctx, "POST", endpoint, params)
	if err != nil {
		return nil, err
	}

	// 检查是否有错误
	if code, ok := result["code"]; ok {
		return nil, fmt.Errorf("wechatpay error: code=%v, message=%v", code, result["message"])
	}

	prepayID, ok := result["prepay_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid response: prepay_id not found")
	}

	return &PrepayResponse{
		PrepayID: prepayID,
	}, nil
}

// 查询订单
func (c *WechatPayClient) QueryOrder(ctx context.Context, outTradeNo string) (map[string]any, error) {
	endpoint := fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s?mchid=%s", outTradeNo, c.mchid)
	return c.CallAPI(ctx, "GET", endpoint, nil)
}

// 关闭订单
func (c *WechatPayClient) CloseOrder(ctx context.Context, outTradeNo string) error {
	body := map[string]any{
		"mchid":        c.mchid,
		"out_trade_no": outTradeNo,
	}
	endpoint := fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s/close", outTradeNo)
	_, err := c.CallAPI(ctx, "POST", endpoint, body)
	return err
}

// 生成APP支付调起参数（返回给客户端）
func (c *WechatPayClient) GenerateAppPayParams(prepayID string) (map[string]any, error) {
	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()
	packageValue := "Sign=WXPay"

	// 生成签名
	// 格式：应用ID\n时间戳\n随机字符串\n预支付交易会话标识\n
	signMessage := fmt.Sprintf("%s\n%d\n%s\n%s\n",
		c.appid,
		timestamp,
		nonceStr,
		prepayID,
	)

	hash := sha256.Sum256([]byte(signMessage))
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	return map[string]any{
		"appid":     c.appid,
		"partnerid": c.mchid,
		"prepayid":  prepayID,
		"package":   packageValue,
		"noncestr":  nonceStr,
		"timestamp": timestamp,
		"sign":      base64.StdEncoding.EncodeToString(signature),
	}, nil
}

// 生成随机字符串
func generateNonceStr() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
