package apperrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	ErrorType int

	ErrorCode interface {
		String() string
	}

	GenericErrorCode string

	AppError struct {
		Type    ErrorType
		Code    ErrorCode
		Message string
	}
)

func (g GenericErrorCode) String() string {
	return string(g)
}

const (
	InternalError ErrorType = iota
	NotFoundError
	ConflictError
	BadRequestError
)

func (err *AppError) Error() string {
	return err.Message
}

func FromGRPCError(err error, errorCode ErrorCode) *AppError {
	if st, ok := status.FromError(err); ok {
		var errorType ErrorType

		switch st.Code() {
		case codes.InvalidArgument:
			errorType = BadRequestError
		case codes.NotFound:
			errorType = NotFoundError
		case codes.AlreadyExists:
			errorType = ConflictError
		default:
			errorType = InternalError
		}

		return &AppError{
			Type:    errorType,
			Code:    errorCode,
			Message: st.Message(),
		}
	}

	return NewInternal(GenericErrorCode("UNKNOWN_ERROR"), err.Error())
}

func (err *AppError) ToGRPCStatus() error {
	var code codes.Code

	switch err.Type {
	case BadRequestError:
		code = codes.InvalidArgument
	case NotFoundError:
		code = codes.NotFound
	case ConflictError:
		code = codes.AlreadyExists

	default:
		code = codes.Internal
	}

	return status.Error(code, err.Message)
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
