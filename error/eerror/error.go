package eerror

import (
	"errors"
	"sync"
)

type (
	ErrorMethod func(*Error) string
	ErrorCode   int
	Stack       []interface{}

	WrapError struct {
		messages []string
	}

	Error struct {
		Cause      error                  `json:"cause,omitempty"`
		StatusCode ErrorCode              `json:"code"`
		Message    string                 `json:"message"`
		StackTrace *reverseSlice          `json:"stack_trace,omitempty"`
		Info       map[string]interface{} `json:"info,omitempty"`
	}

	ErrorBuilder interface {
		Code(ErrorCode) ErrorBuilder
		Stack(item interface{}) ErrorBuilder
		Msg(message string) ErrorBuilder
		Wrap(string) ErrorBuilder
		Field(key string, value interface{}) ErrorBuilder
		Fields(map[string]interface{}) ErrorBuilder
		Err() error
	}
)

var (
	mu               = sync.Mutex{}
	stackDepth       = DefaultStackSize
	errorFunc        = DefaultErrorMethodFunc
	defaultErrorCode = DefaultErrorCode
	wrapDepth        = WrapMessageCap
	wrapSeparator    = DefaultWrapSeparator
	infoCap          = DefaultInfoCap
)

func (e *Error) Error() string {
	return errorFunc(e)
}

func New(msg string) ErrorBuilder {
	return &Error{
		Cause:      newWrapError(msg),
		Message:    msg,
		StatusCode: defaultErrorCode,
	}
}

func Err(err error) ErrorBuilder {
	eErr, ok := err.(*Error)
	if ok {
		return eErr
	}
	return &Error{
		Cause:      err,
		StatusCode: defaultErrorCode,
	}
}

func GetStack(err error) Stack {
	eErr, ok := err.(*Error)
	if ok {
		if eErr.StackTrace == nil {
			return nil
		}
		return eErr.StackTrace.Get()
	}
	return nil
}

func Status(err error) ErrorCode {
	eErr, ok := err.(*Error)
	if ok {
		return eErr.StatusCode
	}
	return NotCode
}

func Cause(err error) error {
	eErr, ok := err.(*Error)
	if ok && eErr.Cause != nil {
		return eErr.Cause
	}
	return nil
}

func OrigCause(err error) error {
	eErr, ok := err.(*Error)
	if ok && eErr.Cause != nil {
		return errors.New(eErr.Cause.(*WrapError).messages[0])
	}
	return nil
}

func Info(err error) map[string]interface{} {
	eErr, ok := err.(*Error)
	if ok && eErr.Info != nil {
		return eErr.Info
	}
	return nil
}

func DefaultAll() {
	mu.Lock()
	defer mu.Unlock()
	defaultErrorCode = DefaultErrorCode
	errorFunc = DefaultErrorMethodFunc
	stackDepth = DefaultStackSize
	wrapDepth = WrapMessageCap
	wrapSeparator = DefaultWrapSeparator
	infoCap = DefaultInfoCap
}
