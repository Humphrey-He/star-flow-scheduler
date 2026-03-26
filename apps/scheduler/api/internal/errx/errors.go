package errx

type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

func New(code int, message string) *BizError {
	return &BizError{Code: code, Message: message}
}

func InvalidParam(message string) *BizError {
	return New(CodeInvalidParam, message)
}

func NotFound(message string) *BizError {
	return New(CodeNotFound, message)
}

func Conflict(message string) *BizError {
	return New(CodeStatusConflict, message)
}

func IdempotentConflict(message string) *BizError {
	return New(CodeIdempotentConflict, message)
}

func Internal(message string) *BizError {
	return New(CodeInternalError, message)
}

func FromError(err error) *BizError {
	if err == nil {
		return nil
	}
	if be, ok := err.(*BizError); ok {
		return be
	}
	return Internal(err.Error())
}
