package info

import (
	"runtime"

	"github.com/vdaas/vald/internal/errors"
)

// Option represent the functional option for info.
type Option func(i *info) error

var (
	defaultOpts = []Option{
		WithRuntimeCaller(runtime.Caller),
		WithRuntimeFuncForPC(runtime.FuncForPC),
	}
)

// WithServerName returns the option to set the server name.
func WithServerName(s string) Option {
	return func(i *info) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("ServerName", s)
		}
		i.detail.ServerName = s
		return nil
	}
}

// WithRuntimeCaller returns the option to set the runtime Caller function.
func WithRuntimeCaller(f func(skip int) (pc uintptr, file string, line int, ok bool)) Option {
	return func(i *info) error {
		if f == nil {
			return errors.NewErrInvalidOption("RuntimeCaller", f)
		}
		i.rtCaller = f
		return nil
	}
}

// WithRuntimeFuncForPC returns the option to set the runtime FuncForPC function.
func WithRuntimeFuncForPC(f func(pc uintptr) *runtime.Func) Option {
	return func(i *info) error {
		if f == nil {
			return errors.NewErrInvalidOption("RuntimeFuncForPC", f)
		}
		i.rtFuncForPC = f
		return nil
	}
}
