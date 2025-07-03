package apperrors

type (
	ErrorType int

	ErrorCode interface {
		String() string
	}

	AppError struct {
		Type    ErrorType
		Code    ErrorCode
		Message string
	}
)

const (
	InternalError ErrorType = iota
	NotFoundError
	ConflictError
	BadRequestError
)

func (err *AppError) Error() string {
	return err.Message
}

func NewInternal(code ErrorCode, msg string) *AppError {
	return &AppError{
		Type:    InternalError,
		Code:    code,
		Message: msg,
	}
}

func NewNotFound(code ErrorCode, msg string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Code:    code,
		Message: msg,
	}
}

func NewConflict(code ErrorCode, msg string) *AppError {
	return &AppError{
		Type:    ConflictError,
		Code:    code,
		Message: msg,
	}
}

func NewBadRequest(code ErrorCode, msg string) *AppError {
	return &AppError{
		Type:    BadRequestError,
		Code:    code,
		Message: msg,
	}
}
