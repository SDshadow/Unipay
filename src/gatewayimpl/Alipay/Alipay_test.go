// 文件路径: gateways/alipay/alipay_test.go
package alipay

import (
	"fmt"
	"strings"
	"testing"

	"Unipay/src/exceptions"
)

// TestNewAlipay 测试创建支付宝实例
func TestNewAlipay(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "正常配置",
			config: map[string]interface{}{
				"app_id":         "2021000123456789",
				"private_key":    "test_private_key",
				"ali_public_key": "test_public_key",
				"notify_url":     "https://example.com/notify",
			},
			expectError: false,
		},
		{
			name: "缺少app_id",
			config: map[string]interface{}{
				"private_key":    "test_private_key",
				"ali_public_key": "test_public_key",
			},
			expectError: true,
			errorMsg:    "Missing Config -- [app_id]",
		},
		{
			name:        "空配置",
			config:      map[string]interface{}{},
			expectError: true,
			errorMsg:    "Missing Config -- [app_id]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alipay, err := NewAlipay(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("期望返回错误，但实际没有错误")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("期望错误包含 '%s'，实际得到 '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望返回错误，但得到: %v", err)
				}
				if alipay == nil {
					t.Error("期望返回 Alipay 实例，但得到 nil")
				}
				if alipay.userConfig == nil {
					t.Error("userConfig 未正确初始化")
				}
			}
		})
	}
}

// TestGetSignContent 测试签名字符串生成
func TestGetSignContent(t *testing.T) {
	config := map[string]interface{}{
		"app_id": "test_app_id",
	}

	alipay, _ := NewAlipay(config)

	// 测试数据
	testData := map[string]interface{}{
		"app_id":    "test_app",
		"method":    "alipay.trade.create",
		"timestamp": "2023-12-01 10:00:00",
		"version":   "1.0",
		"sign":      "should_be_ignored",
	}

	// 测试签名模式
	signContent := alipay.getSignContent(testData, false)
	if strings.Contains(signContent, "sign=") {
		t.Error("签名模式不应该包含 sign 字段")
	}

	// 测试验证模式
	verifyContent := alipay.getSignContent(testData, true)
	if !strings.Contains(verifyContent, "app_id=test_app") {
		t.Error("验证模式应该包含 app_id")
	}

	// 测试排序
	testData2 := map[string]interface{}{
		"z_key": "z_value",
		"a_key": "a_value",
		"m_key": "m_value",
	}
	signContent2 := alipay.getSignContent(testData2, false)
	expectedStart := "a_key=a_value&"
	if !strings.HasPrefix(signContent2, expectedStart) {
		t.Errorf("期望按字母排序，期望开头 '%s'，实际得到 '%s'", expectedStart, signContent2)
	}
}

// TestFormatKeys 测试密钥格式化
func TestFormatKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "短密钥",
			input:    "1234567890",
			expected: "-----BEGIN RSA PRIVATE KEY-----\n1234567890\n-----END RSA PRIVATE KEY-----\n",
		},
		{
			name:  "长密钥",
			input: strings.Repeat("a", 130),
			expected: fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n%s\n%s\n-----END RSA PRIVATE KEY-----\n",
				strings.Repeat("a", 64),
				strings.Repeat("a", 64), strings.Repeat("a", 2)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPrivateKey(tt.input)
			if result != tt.expected {
				t.Errorf("格式化私钥失败\n期望:\n%s\n实际:\n%s", tt.expected, result)
			}
		})
	}
}

// TestGatewayException 测试网关异常处理
func TestGatewayException(t *testing.T) {
	config := map[string]interface{}{
		"app_id": "test_app_id",
	}

	alipay, _ := NewAlipay(config)

	// 测试没有私钥的情况
	_, err := alipay.getSign()
	if err == nil {
		t.Error("期望返回缺少私钥的错误")
	} else if !strings.Contains(err.Error(), "Missing Config -- [private_key]") {
		t.Errorf("期望错误包含 'Missing Config -- [private_key]'，实际得到: %v", err)
	}

	// 验证是否是 InvalidArgumentException
	if _, ok := err.(*exceptions.InvalidArgumentException); !ok {
		t.Error("期望返回 InvalidArgumentException 类型")
	}
}

// TestFindMethod 测试查询订单
func TestFindMethod(t *testing.T) {
	config := map[string]interface{}{
		"app_id":         "test_app_id",
		"private_key":    "test_private_key",
		"ali_public_key": "test_public_key",
	}

	alipay, _ := NewAlipay(config)

	// 测试空订单号
	_, err := alipay.Find("")
	if err == nil {
		t.Error("空订单号应该返回错误")
	}
}
