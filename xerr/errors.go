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
	SystemError = -1
)

func NewXError(code int, msg string) *XError {
	return &XError{code: code, msg: msg}
}

func Errorf(msg string, args ...interface{}) error {
	return &XError{
		code: SystemError,
		msg:  fmt.Sprintf(msg, args...),
	}
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
