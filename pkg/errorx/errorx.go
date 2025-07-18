package errorx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/phuslu/log"
)

type Error interface {
	error
	Code() int
	HTTPStatus() int
	Message() string
	StackTrace() string
	Unwrap() error
	Log() Error
}

type baseError struct {
	code       Code
	httpStatus int
	message    string
	stack      string
	cause      error
}

func (e *baseError) Error() string      { return e.message }
func (e *baseError) Code() int          { return e.code.Int() }
func (e *baseError) HTTPStatus() int    { return e.httpStatus }
func (e *baseError) Message() string    { return e.message }
func (e *baseError) StackTrace() string { return e.stack }
func (e *baseError) Unwrap() error      { return e.cause }

func (e *baseError) Log() Error {
	log.Error().Int("code", e.code.Int()).
		Int("status", e.httpStatus).
		Str("msg", e.message).
		Str("stack", e.stack).
		Msg("error")
	return e
}

func NewErrMsg(msg string) error {
	return &baseError{
		code:       0,
		httpStatus: 0,
		message:    msg,
		stack:      captureStack(3),
	}
}

func NewErrMsgF(format string, args ...interface{}) error {
	return &baseError{
		code:       0,
		httpStatus: 0,
		message:    fmt.Sprintf(format, args...),
		stack:      captureStack(3),
	}
}

func WrapErr(err error, message string) error {
	if err == nil {
		return nil
	}

	return &baseError{
		code:       GetCode(err),
		httpStatus: 0,
		message:    message,
		stack:      captureStack(3),
		cause:      err,
	}
}

func WrapErrMsgF(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &baseError{
		code:       GetCode(err),
		httpStatus: 0,
		message:    fmt.Sprintf(format, args...),
		stack:      captureStack(3),
		cause:      err,
	}
}

// New creates a new base error with stack trace
func New(code int, httpStatus int, message string) Error {
	return &baseError{
		code:       Code(code),
		httpStatus: httpStatus,
		message:    message,
		stack:      captureStack(3),
	}
}

// Wrap wraps an existing error with a new message and captures stack
func Wrap(code int, httpStatus int, msg string, err error) Error {
	return &baseError{
		code:       Code(code),
		httpStatus: httpStatus,
		message:    fmt.Sprintf("%s: %v", msg, err),
		cause:      err,
		stack:      captureStack(3),
	}
}

func GetCode(err error) Code {
	var x baseError
	if ok := As(err, &x); ok {
		return x.code
	}
	return CodeOK
}

// captureStack captures stack trace, skipping `skip` frames
func captureStack(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, &target)
}
