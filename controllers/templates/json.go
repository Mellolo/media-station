package templates

const (
	// 成功
	StatusCodeOK        = "200" // 查询时返回
	StatusCodeCreated   = "201" // 新增时返回
	StatusCodeNoContent = "204" // 更新/删除时返回

	// 客户端错误
	StatusCodeBadRequest   = "400"
	StatusCodeUnauthorized = "401"
	StatusCodeForbidden    = "403"
	StatusCodeNotFound     = "404"

	// 服务端错误
	StatusCodeInternalServerError = "500"
)

type JsonTemplate struct {
	Code    string      `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewJsonTemplate200(data interface{}) JsonTemplate {
	return JsonTemplate{
		Code:    StatusCodeOK,
		Success: true,
		Data:    data,
	}
}

func NewJsonTemplate201(msg string) JsonTemplate {
	return JsonTemplate{
		Code:    StatusCodeCreated,
		Success: true,
		Message: msg,
	}
}

func NewJsonTemplate204(msg string) JsonTemplate {
	return JsonTemplate{
		Code:    StatusCodeNoContent,
		Success: true,
		Message: msg,
	}
}

func NewJsonTemplate400(msg string) JsonTemplate {
	return errorResponse(StatusCodeBadRequest, msg)
}

func NewJsonTemplate401(msg string) JsonTemplate {
	return errorResponse(StatusCodeUnauthorized, msg)
}

func NewJsonTemplate403(msg string) JsonTemplate {
	return errorResponse(StatusCodeForbidden, msg)
}

func NewJsonTemplate404(msg string) JsonTemplate {
	return errorResponse(StatusCodeNotFound, msg)
}

func NewJsonTemplate500(msg string) JsonTemplate {
	return errorResponse(StatusCodeInternalServerError, msg)
}

func errorResponse(statusCode string, msg string) JsonTemplate {
	return JsonTemplate{
		Code:    statusCode,
		Message: msg,
		Success: false,
	}
}
