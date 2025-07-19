package apperrors

type CommonErrorCode string

func (c CommonErrorCode) String() string {
	return string(c)
}

const (
	InternalErrorCode CommonErrorCode = "INTERNAL_ERROR"
)

var (
	InternalServerError = NewInternal(InternalErrorCode, "internal server error")
)
