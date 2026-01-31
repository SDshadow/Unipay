package exceptions

import "fmt"

// InvalidArgumentException 参数异常
type InvalidArgumentException struct {
	*Exception
}

// NewInvalidArgumentException 创建参数异常
func NewInvalidArgumentException(message string) *InvalidArgumentException {
	return &InvalidArgumentException{
		Exception: NewException(message, 400),
	}
}

// NewInvalidArgumentExceptionWithError 创建包含错误的参数异常
func NewInvalidArgumentExceptionWithError(message string, err error) *InvalidArgumentException {
	return &InvalidArgumentException{
		Exception: NewExceptionWithError(message, err, 400),
	}
}

// Error 实现 error 接口
func (e *InvalidArgumentException) Error() string {
	if e.err != nil {
		return fmt.Sprintf("InvalidArgumentException: [code: %d] %s (cause: %v)", e.code, e.message, e.err)
	}
	return fmt.Sprintf("InvalidArgumentException: [code: %d] %s", e.code, e.message)
}
