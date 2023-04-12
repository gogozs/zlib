package xerr

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	XError struct {
		code int
		msg  string
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
	return &XError{code: code, err: err, msg: err.Error()}
}

func NewXError(code int, msg string) *XError {
	return &XError{code: code, err: errors.New(msg), msg: msg}
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
	return e.msg
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
			_, _ = io.WriteString(s, e.msg)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.msg)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *XError) Code() int {
	return e.code
}

// ToGrpcError ...
func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*XError); ok {
		return status.Error(codes.Code(e.Code()), e.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
