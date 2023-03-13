package apperror

import (
	"github.com/mdev5000/flog/attr"
	"net/http"
)

type AppError struct {
	Status      int
	Code        string
	Message     string
	Description string
	Err         error
	Attr        []attr.Attr
}

func (e AppError) Error() string {
	if e.Err == nil {
		return "nil error"
	}
	return e.Err.Error()
}

func (e AppError) Unwrap() error { return e.Err }

func (e AppError) Is(target error) bool {
	_, ok := target.(AppError)
	return ok
}

func InternalError(err error, attrs ...attr.Attr) error {
	return AppError{
		Status:  http.StatusInternalServerError,
		Code:    CodeInternalError,
		Message: "internal server error",
		Err:     err,
		Attr:    attrs,
	}
}

func Error(ec ErrorCode, err error, attrs ...attr.Attr) error {
	return AppError{
		Status:      ec.status,
		Code:        ec.code,
		Message:     ec.message,
		Description: ec.description,
		Err:         err,
		Attr:        attrs,
	}
}
