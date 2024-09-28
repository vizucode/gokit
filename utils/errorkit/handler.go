package errorkit

import "fmt"

type ErrorStd struct {
	HttpStatusCode int
	RpcStatusCode  string
	Message        string
}

func NewErrorStd(HttpStatusCode int, RpcStatusCode string, Message string) *ErrorStd {
	return &ErrorStd{
		HttpStatusCode: HttpStatusCode,
		RpcStatusCode:  RpcStatusCode,
		Message:        Message,
	}
}

func (e *ErrorStd) ErrorCode() string {
	return fmt.Sprintf("%d%s", e.HttpStatusCode, e.RpcStatusCode)
}

func (e *ErrorStd) Error() string {
	return e.Message
}
