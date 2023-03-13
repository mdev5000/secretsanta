package apperror

import "net/http"

var AlreadySetup = ErrorCode{
	status:  http.StatusInternalServerError,
	code:    CodeInternalError,
	message: "already setup",
}
