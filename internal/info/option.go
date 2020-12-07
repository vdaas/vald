package info

import "runtime"

type Option func(i *info) error

var (
	defaultOpts = []Option{
		WithRuntimeCaller(runtime.Caller),
		WithRuntimeFuncForPC(runtime.FuncForPC),
	}
)

func WithServerName(s string) Option {
	return func(i *info) error {
		i.detail.ServerName = s
		return nil
	}
}

func WithRuntimeCaller(f func(skip int) (pc uintptr, file string, line int, ok bool)) Option {
	return func(i *info) error {
		i.rtCaller = f
		return nil
	}
}

func WithRuntimeFuncForPC(f func(pc uintptr) *runtime.Func) Option {
	return func(i *info) error {
		i.rtFuncForPC = f
		return nil
	}
}
