package errors

import "fmt"

// 错误码定义
const (
	// 1xxxx 参数/请求层
	CodeInvalidParam = 10001

	// 2xxxx 业务数据层
	CodeOrderNotFound     = 20001
	CodeOrderStatusInvalid = 20002

	// 3xxxx 上游依赖层
	CodeRpcCallFailed  = 30001
	CodeRpcRespInvalid = 30002

	// 4xxxx 系统/基础设施层
	CodeDBError       = 40001
	CodeInternalError = 40002
)

// AppError 自定义业务错误
type AppError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

// New 创建 AppError
func New(code int, msg string) *AppError {
	return &AppError{Code: code, Msg: msg}
}

// Newf 创建带格式化消息的 AppError
func Newf(code int, format string, args ...any) *AppError {
	return &AppError{Code: code, Msg: fmt.Sprintf(format, args...)}
}
