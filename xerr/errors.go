package xerr

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

type (
	XError struct {
		code int
		err  error
	}
)

const (
	InternalError = -1
)

func NewXErrorByError(code int, err error) *XError {
	if err == nil {
		return nil
	}
	return &XError{code: code, err: err}
}

func NewXError(code int, msg string) *XError {
	return &XError{code: code, err: errors.New(msg)}
}

func Errorf(msg string, args ...interface{}) error {
	return &XError{
		code: InternalError,
		err:  errors.Errorf(msg, args...),
	}
}

func Wrap(err error, msg string) *XError {
	if err == nil {
		return nil
	}
	if xe, ok := err.(*XError); ok {
		xe.Wrap(msg)
		return xe
	}
	return NewXErrorByError(InternalError, errors.Wrap(err, msg))
}

func (e *XError) Wrap(msg string) {
	e.err = errors.Wrap(e.err, msg)
}

func (e *XError) Error() string {
	return e.err.Error()
}

func (e *XError) Format(s fmt.State, verb rune) {
	f, ok := e.err.(fmt.Formatter)
	if ok {
		f.Format(s, verb)
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.Error())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *XError) Code() int {
	return e.code
}

// ToGrpcError ...
func ToGrpcError(err error) error {
	return status.Convert(err).Err()
}
