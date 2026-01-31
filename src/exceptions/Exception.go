package exceptions

import (
	"fmt"
)

// Exception 基础异常类
type Exception struct {
	message string
	code    int
	err     error
}

// NewException 创建基础异常
func NewException(message string, code ...int) *Exception {
	exceptionCode := 0
	if len(code) > 0 {
		exceptionCode = code[0]
	}

	return &Exception{
		message: message,
		code:    exceptionCode,
	}
}

// NewExceptionWithError 创建包含底层错误的异常
func NewExceptionWithError(message string, err error, code ...int) *Exception {
	exceptionCode := 0
	if len(code) > 0 {
		exceptionCode = code[0]
	}

	return &Exception{
		message: message,
		code:    exceptionCode,
		err:     err,
	}
}

// Error 实现 error 接口
func (e *Exception) Error() string {
	if e.err != nil {
		return fmt.Sprintf("Exception: [code: %d] %s (cause: %v)", e.code, e.message, e.err)
	}
	return fmt.Sprintf("Exception: [code: %d] %s", e.code, e.message)
}

// Unwrap 支持 errors.Unwrap
func (e *Exception) Unwrap() error {
	return e.err
}

// GetCode 获取异常代码
func (e *Exception) GetCode() int {
	return e.code
}

// GetMessage 获取异常消息
func (e *Exception) GetMessage() string {
	return e.message
}

// // 文件路径: exceptions/business_exception.go
// package exceptions

// // BusinessException 业务异常
// type BusinessException struct {
// 	*Exception
// }

// // NewBusinessException 创建业务异常
// func NewBusinessException(message string, code ...int) *BusinessException {
// 	businessCode := 500
// 	if len(code) > 0 {
// 		businessCode = code[0]
// 	}

// 	return &BusinessException{
// 		Exception: NewException(message, businessCode),
// 	}
// }

// // 文件路径: exceptions/notify_exception.go
// package exceptions

// // NotifyException 通知异常
// type NotifyException struct {
// 	*Exception
// 	notifyData map[string]interface{}
// }

// // NewNotifyException 创建通知异常
// func NewNotifyException(message string, notifyData map[string]interface{}, code ...int) *NotifyException {
// 	notifyCode := 422
// 	if len(code) > 0 {
// 		notifyCode = code[0]
// 	}

// 	return &NotifyException{
// 		Exception:  NewException(message, notifyCode),
// 		notifyData: notifyData,
// 	}
// }

// // GetNotifyData 获取通知数据
// func (e *NotifyException) GetNotifyData() map[string]interface{} {
// 	return e.notifyData
// }
