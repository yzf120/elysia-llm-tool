package errs

import (
	"encoding/json"
	"strconv"
	"strings"
)

// 错误码定义
const (
	// 通用错误
	ErrInternalServer = 500 // 服务内部错误
	ErrBadRequest     = 400 // 请求参数错误

	// LLM相关错误 (30000+)
	ErrModelNotFound      = 30001 // 模型不存在
	ErrModelRequestFailed = 30002 // 模型请求失败
	ErrModelTimeout       = 30003 // 模型请求超时
	ErrModelStreamFailed  = 30004 // 流式响应失败
	ErrInvalidModelID     = 30005 // 无效的模型ID
	ErrProviderNotSupport = 30006 // 不支持的提供商

	// 参数错误 (40000+)
	ErrParamMissing  = 40001 // 参数缺失
	ErrParamInvalid  = 40002 // 参数无效
	ErrMessageEmpty  = 40003 // 消息为空
	ErrMessageFormat = 40004 // 消息格式错误

	// 配置错误 (50000+)
	ErrConfigMissing = 50001 // 配置缺失
	ErrAPIKeyMissing = 50002 // API Key缺失
)

// 错误消息映射
var ErrorMessages = map[int]string{
	ErrInternalServer:     "服务内部错误",
	ErrBadRequest:         "请求参数错误",
	ErrModelNotFound:      "模型不存在",
	ErrModelRequestFailed: "模型请求失败",
	ErrModelTimeout:       "模型请求超时",
	ErrModelStreamFailed:  "流式响应失败",
	ErrInvalidModelID:     "无效的模型ID",
	ErrProviderNotSupport: "不支持的提供商",
	ErrParamMissing:       "参数缺失",
	ErrParamInvalid:       "参数无效",
	ErrMessageEmpty:       "消息不能为空",
	ErrMessageFormat:      "消息格式错误",
	ErrConfigMissing:      "配置缺失",
	ErrAPIKeyMissing:      "API Key未配置",
}

// CommonError 通用错误结构
type CommonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CommonError) Error() string {
	return e.Message
}

// NewCommonError 创建通用错误
func NewCommonError(code int, msg string) *CommonError {
	return &CommonError{
		Code:    code,
		Message: msg,
	}
}

// NewCommonErrorWithDefault 创建通用错误（使用默认消息）
func NewCommonErrorWithDefault(code int) *CommonError {
	msg := ErrorMessages[code]
	if msg == "" {
		msg = "未知错误"
	}
	return &CommonError{
		Code:    code,
		Message: msg,
	}
}

// ParseCommonError 解析 CommonError.Error() 返回的错误字符串
// 格式: "[code]message" 例如: "[30001]模型不存在"
// 返回: code (int), message (string)
// 如果解析失败，返回 500 和原始错误字符串
func ParseCommonError(errStr string) (int, string) {
	// 检查格式是否为 "[code]message"
	if !strings.HasPrefix(errStr, "[") {
		return ErrInternalServer, errStr
	}

	// 查找右括号的位置
	closeBracketIdx := strings.Index(errStr, "]")
	if closeBracketIdx == -1 {
		return ErrInternalServer, errStr
	}

	// 提取 code 部分
	codeStr := errStr[1:closeBracketIdx]
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return ErrInternalServer, errStr
	}

	// 提取 message 部分
	message := ""
	if closeBracketIdx+1 < len(errStr) {
		message = errStr[closeBracketIdx+1:]
	}

	return code, message
}

// CommonResponse 通用响应结构
type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *CommonResponse {
	return &CommonResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, msg string) *CommonResponse {
	errMsg := msg
	// 优先使用预设的错误消息
	if codeMsg, exists := ErrorMessages[code]; exists && msg == "" {
		errMsg = codeMsg
	}
	return &CommonResponse{
		Code:    code,
		Message: errMsg,
	}
}

// ToJSON 转换为JSON字符串
func (r *CommonResponse) ToJSON() string {
	jsonStr, _ := json.Marshal(r)
	return string(jsonStr)
}

// IsSuccess 判断是否成功
func (r *CommonResponse) IsSuccess() bool {
	return r.Code == 0
}
