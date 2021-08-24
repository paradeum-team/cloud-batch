package e

import (
	"github.com/pkg/errors"
	"net/http"
)

type ErrorCode struct {
	code         string
	httpCode     int
	TemplateData interface{}
	PluralCount  interface{}
}

// New ErroCode
func New(code string, httpCode int) ErrorCode {
	return ErrorCode{code: code, httpCode: httpCode}
}

func (ec *ErrorCode) GetCode() string {
	return ec.code
}

func (ec *ErrorCode) GetHttpCode() int {
	return ec.httpCode
}

// ErrorCode list
var (
	StatusOK                  = New("StatusOK", http.StatusOK)                          // 正常响应，无异常
	InternalError             = New("InternalError", http.StatusInternalServerError)    // 内部错误
	DeprecatedURI             = New("DeprecatedURI", http.StatusBadRequest)             // URI已经弃用
	MalformedJSON             = New("MalformedJSON", http.StatusBadRequest)             // JSON 格式错误
	InappropriateJSON         = New("InappropriateJSON", http.StatusBadRequest)         // JSON 值检验不通过
	InvalidParameter          = New("InvalidParameter", http.StatusBadRequest)          // 参数值校验不通过, 一般用于多个或未知参数不通过
	AuthError                 = New("AuthError", http.StatusUnauthorized)               // 账号或密码错误
	InvalidToken              = New("InvalidToken", http.StatusUnauthorized)            // Token验证失败，请重新登录
	Asynchronization          = New("Asynchronization", http.StatusAccepted)            // 异步处理中，请稍后重试
	RequestTimeout            = New("RequestTimeout", http.StatusRequestTimeout)        // 请求超时，请稍后重试。
	InvalidParameterLength    = New("InvalidParameterLength", http.StatusBadRequest)    // 参数必须在0 - {{.MaxLength))之间
	ParametherNotAllowedEmpty = New("ParametherNotAllowedEmpty", http.StatusBadRequest) // {{.Name))参数不允许为空
	InvalidParameterValue     = New("InvalidParameterValue", http.StatusBadRequest)     // 参数{{.Name))值校验不通过, 用于确定单一参数
	NotFound                  = New("NotFound", http.StatusNotFound)
)

var (
	ErrServerCreating          = errors.New("主机创建中，请稍后。")
	ErrServerCreateTimeout     = errors.New("主机创建超时，请联系管理员。")
	ErrServerCreateError       = errors.New("主机创建失败，请联系管理员。")
	ErrUnknownError            = errors.New("未知错误")
	ErrServerCreateStatusEmpty = errors.New("该批次主机创建列表为空，请重新创建。")

	ErrStatusEmpty    = errors.New("status empty.")
	ErrStatusCreating = errors.New("创建中，请稍后。")
	ErrStatusTimeout  = errors.New("操作超时，请联系管理员。")
	ErrStatusStart    = errors.New("正在进行中, 请稍后。")
	ErrStatusError    = errors.New("操作失败。")
	ErrStatusDone     = errors.New("已经完成。")
	ErrStatusConflict = errors.New("操作冲突, 请联系管理员。")
)

func GetCodeByErr(err error) ErrorCode {
	if errors.Is(err, ErrServerCreating) {
		return Asynchronization
	}
	if errors.Is(err, ErrStatusStart) {
		return Asynchronization
	}
	if errors.Is(err, ErrStatusDone) {
		return StatusOK
	}
	return InternalError
}

const (
	StatusStart    = "start"
	StatusError    = "error"
	StatusDone     = "done"
	StatusConflict = "conflict"
	StatusTimeout = "timeout"
)
