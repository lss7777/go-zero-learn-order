package response

// ApiResponse 统一响应格式
type ApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 成功响应
func Success(data any) ApiResponse {
	return ApiResponse{Code: 0, Msg: "success", Data: data}
}

// Error 错误响应
func Error(code int, msg string) ApiResponse {
	return ApiResponse{Code: code, Msg: msg}
}
