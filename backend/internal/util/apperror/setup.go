package apperror

import "net/http"

var AlreadySetup = ErrorCode{
	status:  http.StatusBadRequest,
	code:    CodeBadRequest,
	message: "already setup",
}
