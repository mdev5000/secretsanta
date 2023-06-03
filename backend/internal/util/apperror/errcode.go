package apperror

import "net/http"

const CodeInternalError = "internal-001"
const CodeBadRequest = "internal-002"
const CodeEchoError = "internal-003"
const CodeInvalidLogin = "auth-001"
const CodeNotAuthenticated = "auth-002"

type ErrorCode struct {
	Status      int
	Code        string
	Message     string
	Description string
}

func (e *ErrorCode) WithDescription(description string) ErrorCode {
	return ErrorCode{
		Status:      e.Status,
		Code:        e.Code,
		Message:     e.Message,
		Description: description,
	}
}

var InvalidLogin = ErrorCode{
	Status:  http.StatusBadRequest,
	Code:    CodeInvalidLogin,
	Message: "invalid login",
}

var BadRequest = ErrorCode{
	Status:  http.StatusBadRequest,
	Code:    CodeBadRequest,
	Message: "invalid request",
}

var NotAuthenticated = ErrorCode{
	Status:  http.StatusUnauthorized,
	Code:    CodeNotAuthenticated,
	Message: "access denied",
}
