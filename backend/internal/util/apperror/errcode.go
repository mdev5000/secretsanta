package apperror

import "net/http"

const CodeInternalError = "internal-001"
const CodeBadRequest = "internal-002"

type ErrorCode struct {
	status      int
	code        string
	message     string
	description string
}

func (e *ErrorCode) WithDescription(description string) ErrorCode {
	return ErrorCode{
		status:      e.status,
		code:        e.code,
		message:     e.message,
		description: description,
	}
}

var BadRequest = ErrorCode{
	status:  http.StatusBadRequest,
	code:    CodeBadRequest,
	message: "invalid request",
}
