package xerr

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	XError struct {
		code int
		msg  string
	}
)

const (
	InternalError = -1
)

func NewXError(code int, msg string) *XError {
	return &XError{code: code, msg: msg}
}

func Errorf(msg string, args ...interface{}) error {
	return &XError{
		code: InternalError,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func Wrap(err error, msg string) *XError {
	if err == nil {
		return nil
	}
	if xe, ok := err.(*XError); ok {
		return NewXError(xe.code, WrapMessage(err, msg))
	}
	return NewXError(InternalError, WrapMessage(err, msg))
}

func WrapMessage(err error, msg string) string {
	if err != nil {
		return fmt.Sprintf("%s: %s", err, err.Error())
	}
	return msg
}

func (e XError) Error() string {
	return e.msg
}

func (e XError) Code() int {
	return e.code
}

// ToGrpcError ...
func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(XError); ok {
		return status.Error(codes.Code(e.Code()), e.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
