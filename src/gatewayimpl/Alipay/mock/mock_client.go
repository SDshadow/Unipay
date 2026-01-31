// 文件路径: gateways/alipay/mock/mock_client.go
package mock

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MockHTTPClient 模拟 HTTP 客户端
type MockHTTPClient struct {
	// 预设的响应
	responses map[string]*MockResponse
	// 记录的请求
	requests []*http.Request
}

// MockResponse 模拟响应
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// NewMockHTTPClient 创建模拟客户端
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]*MockResponse),
		requests:  make([]*http.Request, 0),
	}
}

// SetResponse 设置响应
func (m *MockHTTPClient) SetResponse(url string, response *MockResponse) {
	m.responses[url] = response
}

// GetRequests 获取所有请求
func (m *MockHTTPClient) GetRequests() []*http.Request {
	return m.requests
}

// ClearRequests 清空请求记录
func (m *MockHTTPClient) ClearRequests() {
	m.requests = make([]*http.Request, 0)
}

// Do 实现 http.Client 接口
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// 记录请求
	m.requests = append(m.requests, req)

	// 查找预设的响应
	resp, exists := m.responses[req.URL.String()]
	if !exists {
		// 默认响应
		resp = &MockResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("No mock response for URL: %s", req.URL.String()),
		}
	}

	// 创建响应
	return &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(strings.NewReader(resp.Body)),
		Header:     make(http.Header),
	}, nil
}

// MockConfig 模拟配置
func MockConfig() map[string]interface{} {
	return map[string]interface{}{
		"app_id":         "2021000123456789",
		"private_key":    "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC4xQ5CQ3sGxU1L\n... 模拟的私钥 ...",
		"ali_public_key": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMUOQkN7BsVNS0\n... 模拟的公钥 ...",
		"notify_url":     "https://example.com/notify",
		"return_url":     "https://example.com/return",
	}
}
