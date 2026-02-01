package exceptions

import (
	"fmt"
)

// GatewayException 网关异常
type GatewayException struct {
	// 异常消息
	message string

	// 异常代码
	code int

	// 原始响应数据
	raw map[string]interface{}

	// 底层错误
	err error
}

// NewGatewayException 创建网关异常
// message: 异常消息
// code: 异常代码（可以是int或string，会被转换为int）
// raw: 原始响应数据
func NewGatewayException(message string, code interface{}, raw ...map[string]interface{}) *GatewayException {
	exception := &GatewayException{
		message: message,
	}

	// 处理代码
	switch v := code.(type) {
	case int:
		exception.code = v
	case int8:
		exception.code = int(v)
	case int16:
		exception.code = int(v)
	case int32:
		exception.code = int(v)
	case int64:
		exception.code = int(v)
	case uint:
		exception.code = int(v)
	case uint8:
		exception.code = int(v)
	case uint16:
		exception.code = int(v)
	case uint32:
		exception.code = int(v)
	case uint64:
		exception.code = int(v)
	case string:
		// 尝试将字符串转换为数字
		var intCode int
		_, err := fmt.Sscanf(v, "%d", &intCode)
		if err == nil {
			exception.code = intCode
		} else {
			// 如果转换失败，使用哈希值作为代码
			exception.code = hashString(v)
		}
	default:
		// 其他类型，使用默认值
		exception.code = 0
	}

	// 处理原始数据
	if len(raw) > 0 {
		exception.raw = raw[0]
	} else {
		exception.raw = make(map[string]interface{})
	}

	return exception
}

// NewGatewayExceptionWithError 创建包含底层错误的网关异常
func NewGatewayExceptionWithError(message string, err error, code ...interface{}) *GatewayException {
	var exceptionCode interface{} = 0
	if len(code) > 0 {
		exceptionCode = code[0]
	}

	exception := NewGatewayException(message, exceptionCode)
	exception.err = err
	return exception
}

// Error 实现 error 接口
func (e *GatewayException) Error() string {
	if e.err != nil {
		return fmt.Sprintf("GatewayException: [code: %d] %s (cause: %v)", e.code, e.message, e.err)
	}
	return fmt.Sprintf("GatewayException: [code: %d] %s", e.code, e.message)
}

// Unwrap 支持 errors.Unwrap
func (e *GatewayException) Unwrap() error {
	return e.err
}

// GetCode 获取异常代码
func (e *GatewayException) GetCode() int {
	return e.code
}

// GetMessage 获取异常消息
func (e *GatewayException) GetMessage() string {
	return e.message
}

// GetRaw 获取原始响应数据
func (e *GatewayException) GetRaw() map[string]interface{} {
	return e.raw
}

// SetRaw 设置原始响应数据
func (e *GatewayException) SetRaw(raw map[string]interface{}) {
	e.raw = raw
}

// AddRawData 添加原始数据
func (e *GatewayException) AddRawData(key string, value interface{}) {
	if e.raw == nil {
		e.raw = make(map[string]interface{})
	}
	e.raw[key] = value
}

// Is 检查是否为网关异常
func (e *GatewayException) Is(target error) bool {
	_, ok := target.(*GatewayException)
	return ok
}

// 辅助函数：字符串哈希
func hashString(s string) int {
	hash := 0
	for _, char := range s {
		hash = 31*hash + int(char)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash % 10000
}
