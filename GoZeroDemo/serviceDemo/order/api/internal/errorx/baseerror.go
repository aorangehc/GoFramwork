package errorx

const (
	defaultErrCode = 1001
	RPCErrCODE     = 1002
)

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CodeErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewCodeError(code int, message string) error {
	return &CodeError{
		Code:    code,
		Message: message,
	}
}

func NewDefaultCodeError(message string) error {
	return &CodeError{
		Code:    defaultErrCode,
		Message: message,
	}
}

func (e *CodeError) Error() string {
	return e.Message
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}
