package apperror

const CodeInternalError = "internal-001"

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
