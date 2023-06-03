package apperror

import "net/http"

var AlreadySetup = ErrorCode{
	Status:  http.StatusBadRequest,
	Code:    CodeBadRequest,
	Message: "already setup",
}
