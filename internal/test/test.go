package test

import (
	"context"
	"testing"
)

type Test interface {
	Run(context.Context, *testing.T)
}

type test struct {
	cs     []Caser
	target func(ctx context.Context,
		t *testing.T,
		args, fields []interface{},
		checkFunc func(t *testing.T, gots ...interface{}),
	)
}

func New(opts ...Option) Test {
	t := new(test)
	for _, opt := range append(defaultOptions, opts...) {
		opt(t)
	}
	return t
}

func (test *test) Run(ctx context.Context, t *testing.T) {
	t.Helper()
	for _, c := range test.cs {
		t.Run(c.Name(), func(tt *testing.T) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			test.target(ctx, tt, c.Args(), c.Fields(), c.CheckFunc())
		})
	}
}
