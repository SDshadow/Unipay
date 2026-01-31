package alipay

import (
	"bytes"
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
	"sort"
	"strings"
	"time"

	"Unipay/src/contracts"
	"Unipay/src/exceptions"
	"Unipay/src/support"
)

type Alipay struct {
	// 网关地址
	gateway string

	// 支付宝配置
	config map[string]interface{}

	// 用户配置
	userConfig *support.Config

	// HTTP 客户端
	httpClient *http.Client
}

func NewAlipay(config map[string]interface{}) (*Alipay, error) {
	alipay := &Alipay{
		gateway:    "https://openapi.alipay.com/gateway.do?charset=utf-8",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	// 初始化用户配置
	userConfig := support.NewConfig(config)
	alipay.userConfig = userConfig

	// 验证必要参数
	if userConfig.GetString("app_id") == "" {
		return nil, exceptions.NewInvalidArgumentException("Missing Config -- [app_id]")
	}

	// 初始化配置
	alipay.config = map[string]interface{}{
		"app_id":      userConfig.GetString("app_id"),
		"method":      "",
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"version":     "1.0",
		"return_url":  userConfig.GetString("return_url", ""),
		"notify_url":  userConfig.GetString("notify_url", ""),
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"sign":        "",
		"biz_content": "",
	}

	return alipay, nil
}

// Pay 支付订单
func (a *Alipay) Pay(configBiz map[string]interface{}) (*contracts.PaymentResult, error) {
	// 设置产品代码
	configBiz["product_code"] = a.getProductCode()

	// 设置方法名
	a.config["method"] = a.getMethod()

	// 序列化业务参数
	bizContent, err := json.Marshal(configBiz)
	if err != nil {
		return nil, exceptions.NewGatewayException("json marshal error", err)
	}
	a.config["biz_content"] = string(bizContent)

	// 生成签名
	sign, err := a.getSign()
	if err != nil {
		return nil, err
	}
	a.config["sign"] = sign

	// 构建支付表单（网页支付场景）
	// 如果是移动端APP支付，这里需要返回不同的结果
	return a.buildPaymentResult()
}

// Refund 退款订单
func (a *Alipay) Refund(configBiz interface{}) (interface{}, error) {
	var params map[string]interface{}

	switch v := configBiz.(type) {
	case map[string]interface{}:
		params = v
	case string:
		// 如果是字符串，认为是订单号
		params = map[string]interface{}{
			"out_trade_no": v,
		}
	default:
		return nil, exceptions.NewInvalidArgumentException("invalid refund config")
	}

	return a.getResult(params, "alipay.trade.refund")
}

// Close 关闭订单
func (a *Alipay) Close(configBiz interface{}) (interface{}, error) {
	var params map[string]interface{}

	switch v := configBiz.(type) {
	case map[string]interface{}:
		params = v
	case string:
		params = map[string]interface{}{
			"out_trade_no": v,
		}
	default:
		return nil, exceptions.NewInvalidArgumentException("invalid close config")
	}

	return a.getResult(params, "alipay.trade.close")
}

// Find 查询订单
func (a *Alipay) Find(outTradeNo string) (interface{}, error) {
	params := map[string]interface{}{
		"out_trade_no": outTradeNo,
	}

	return a.getResult(params, "alipay.trade.query")
}

// Verify 验证通知
func (a *Alipay) Verify(data interface{}, sign ...string) (interface{}, error) {
	// 检查支付宝公钥
	if a.userConfig.GetString("ali_public_key") == "" {
		return nil, exceptions.NewInvalidArgumentException("Missing Config -- [ali_public_key]")
	}

	// 获取签名
	var signature string
	if len(sign) > 0 && sign[0] != "" {
		signature = sign[0]
	} else {
		// 从data中提取签名
		if dataMap, ok := data.(map[string]interface{}); ok {
			if s, exists := dataMap["sign"].(string); exists {
				signature = s
			}
		}
	}

	if signature == "" {
		return false, exceptions.NewInvalidArgumentException("signature is required")
	}

	// 解析数据
	var dataMap map[string]interface{}
	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case string:
		err := json.Unmarshal([]byte(v), &dataMap)
		if err != nil {
			return false, exceptions.NewInvalidArgumentException("invalid data format")
		}
	default:
		return false, exceptions.NewInvalidArgumentException("invalid data type")
	}

	// 验证签名
	isSync := len(sign) > 1 && sign[1] == "sync"
	valid, err := a.verifySignature(dataMap, signature, isSync)
	if err != nil {
		return false, err
	}

	if valid {
		return dataMap, nil
	}

	return false, nil
}

// 抽象方法 - 需要在具体实现中实现

// getMethod 获取方法名
func (a *Alipay) getMethod() string {
	// 需要在子类中实现
	return ""
}

// getProductCode 获取产品代码
func (a *Alipay) getProductCode() string {
	// 需要在子类中实现
	return ""
}

// buildPaymentResult 构建支付结果
func (a *Alipay) buildPaymentResult() (*contracts.PaymentResult, error) {
	// 需要在具体支付方式中实现
	return nil, fmt.Errorf("method not implemented")
}

// 辅助方法

// getResult 调用支付宝API获取结果
func (a *Alipay) getResult(configBiz map[string]interface{}, method string) (map[string]interface{}, error) {
	// 序列化业务参数
	bizContent, err := json.Marshal(configBiz)
	if err != nil {
		return nil, exceptions.NewGatewayException("json marshal error", err)
	}

	// 更新配置
	a.config["biz_content"] = string(bizContent)
	a.config["method"] = method

	// 生成签名
	sign, err := a.getSign()
	if err != nil {
		return nil, err
	}
	a.config["sign"] = sign

	// 发送请求
	resp, err := a.post(a.gateway, a.config)
	if err != nil {
		return nil, exceptions.NewGatewayException("request alipay error", err)
	}

	// 解析响应
	responseKey := strings.ReplaceAll(method, ".", "_") + "_response"

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, exceptions.NewGatewayException("json unmarshal error", err)
	}

	responseData, ok := result[responseKey].(map[string]interface{})
	if !ok {
		return nil, exceptions.NewGatewayException("invalid response format", nil)
	}

	// 检查返回码
	code, _ := responseData["code"].(string)
	if code != "10000" {
		msg, _ := responseData["msg"].(string)
		subCode, _ := responseData["sub_code"].(string)
		return nil, exceptions.NewGatewayException(
			fmt.Sprintf("alipay error: %s - %s", msg, subCode),
			fmt.Errorf("code: %s", code),
		)
	}

	// 验证签名
	signature, _ := result["sign"].(string)
	verified, err := a.verifySignature(responseData, signature, true)
	if err != nil || !verified {
		return nil, exceptions.NewGatewayException("signature verification failed", err)
	}

	return responseData, nil
}

// getSign 生成签名
func (a *Alipay) getSign() (string, error) {
	privateKey := a.userConfig.GetString("private_key")
	if privateKey == "" {
		return "", exceptions.NewInvalidArgumentException("Missing Config -- [private_key]")
	}

	// 格式化私钥
	formattedKey := formatPrivateKey(privateKey)

	// 解析PKCS1私钥
	block, _ := pem.Decode([]byte(formattedKey))
	if block == nil {
		return "", exceptions.NewInvalidArgumentException("invalid private key format")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", exceptions.NewInvalidArgumentException(fmt.Sprintf("parse private key error: %v", err))
	}

	// 获取待签名字符串
	signContent := a.getSignContent(a.config, false)

	// 计算SHA256哈希
	hash := sha256.New()
	hash.Write([]byte(signContent))
	hashed := hash.Sum(nil)

	// 使用RSA签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed)
	if err != nil {
		return "", exceptions.NewGatewayException("sign error", err)
	}

	// Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// verifySignature 验证签名
func (a *Alipay) verifySignature(data map[string]interface{}, sign string, sync bool) (bool, error) {
	publicKey := a.userConfig.GetString("ali_public_key")
	if publicKey == "" {
		return false, exceptions.NewInvalidArgumentException("Missing Config -- [ali_public_key]")
	}

	// 格式化公钥
	formattedKey := formatPublicKey(publicKey)

	// 解析公钥
	block, _ := pem.Decode([]byte(formattedKey))
	if block == nil {
		return false, exceptions.NewInvalidArgumentException("invalid public key format")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, exceptions.NewInvalidArgumentException(fmt.Sprintf("parse public key error: %v", err))
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return false, exceptions.NewInvalidArgumentException("not an RSA public key")
	}

	// 获取待验证字符串
	var toVerify string
	if sync {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return false, err
		}
		toVerify = string(jsonData)
	} else {
		toVerify = a.getSignContent(data, true)
	}

	// 解码签名
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, exceptions.NewInvalidArgumentException("invalid base64 signature")
	}

	// 计算SHA256哈希
	hash := sha256.New()
	hash.Write([]byte(toVerify))
	hashed := hash.Sum(nil)

	// 验证签名
	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hashed, signBytes)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// getSignContent 获取待签名字符串
func (a *Alipay) getSignContent(toBeSigned map[string]interface{}, verify bool) string {
	// 按键排序
	keys := make([]string, 0, len(toBeSigned))
	for k := range toBeSigned {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var stringToBeSigned strings.Builder
	for _, k := range keys {
		v := toBeSigned[k]

		// 跳过空值和签名相关字段
		if v == nil || v == "" {
			continue
		}

		vStr := fmt.Sprintf("%v", v)

		if verify {
			// 验证时跳过sign和sign_type
			if k == "sign" || k == "sign_type" {
				continue
			}
			stringToBeSigned.WriteString(fmt.Sprintf("%s=%s&", k, vStr))
		} else {
			// 签名时跳过sign字段，且不处理以@开头的文件上传参数
			if k == "sign" || strings.HasPrefix(vStr, "@") {
				continue
			}
			stringToBeSigned.WriteString(fmt.Sprintf("%s=%s&", k, vStr))
		}
	}

	result := stringToBeSigned.String()
	if len(result) > 0 {
		result = result[:len(result)-1] // 移除最后一个&
	}

	return result
}

// post 发送POST请求
func (a *Alipay) post(urlStr string, params map[string]interface{}) ([]byte, error) {
	// 创建 form 数据
	data := make(map[string][]string)
	for k, v := range params {
		data[k] = []string{fmt.Sprintf("%v", v)}
	}

	// 使用 http.PostForm
	resp, err := http.PostForm(urlStr, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// buildPayHtml 构建支付表单（网页支付用）
func (a *Alipay) buildPayHtml() string {
	var html strings.Builder

	html.WriteString(fmt.Sprintf("<form id='alipaysubmit' name='alipaysubmit' action='%s' method='POST'>", a.gateway))

	for k, v := range a.config {
		vStr := fmt.Sprintf("%v", v)
		vStr = strings.ReplaceAll(vStr, "'", "&apos;")
		html.WriteString(fmt.Sprintf("<input type='hidden' name='%s' value='%s'/>", k, vStr))
	}

	html.WriteString("<input type='submit' value='ok' style='display:none;'></form>")
	html.WriteString("<script>document.forms['alipaysubmit'].submit();</script>")

	return html.String()
}

// 密钥格式化辅助函数

func formatPrivateKey(privateKey string) string {
	// 清理密钥中的空白字符
	privateKey = strings.TrimSpace(privateKey)
	privateKey = strings.ReplaceAll(privateKey, "\n", "")
	privateKey = strings.ReplaceAll(privateKey, "\r", "")
	privateKey = strings.ReplaceAll(privateKey, " ", "")

	var buf bytes.Buffer
	buf.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")

	// 每64个字符换行
	for i := 0; i < len(privateKey); i += 64 {
		end := i + 64
		if end > len(privateKey) {
			end = len(privateKey)
		}
		buf.WriteString(privateKey[i:end])
		buf.WriteString("\n")
	}

	buf.WriteString("-----END RSA PRIVATE KEY-----\n")
	return buf.String()
}

func formatPublicKey(publicKey string) string {
	// 清理密钥中的空白字符
	publicKey = strings.TrimSpace(publicKey)
	publicKey = strings.ReplaceAll(publicKey, "\n", "")
	publicKey = strings.ReplaceAll(publicKey, "\r", "")
	publicKey = strings.ReplaceAll(publicKey, " ", "")

	var buf bytes.Buffer
	buf.WriteString("-----BEGIN PUBLIC KEY-----\n")

	// 每64个字符换行
	for i := 0; i < len(publicKey); i += 64 {
		end := i + 64
		if end > len(publicKey) {
			end = len(publicKey)
		}
		buf.WriteString(publicKey[i:end])
		buf.WriteString("\n")
	}

	buf.WriteString("-----END PUBLIC KEY-----\n")
	return buf.String()
}
