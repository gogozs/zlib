package tools

import (
	"runtime/debug"

	"github.com/gogozs/zlib/xlog"
	"github.com/pkg/errors"
)

func SafeRun(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ToPanicError(r)
		}
	}()
	return f()
}

func SafeGo(f func()) {
	go func() {
		_ = SafeRun(func() error {
			f()
			return nil
		})
	}()
}

func ToPanicError(r any) error {
	xlog.Error("%+v\n\n%s", r, debug.Stack())
	return errors.Errorf("panic: %v", r)
}

func HandleCrash(handler func(err error)) {
	if r := recover(); r != nil {
		handler(ToPanicError(r))
	}
}
