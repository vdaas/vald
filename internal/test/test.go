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
	target func(ctx context.Context, c Caser) error
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

			err := test.target(ctx, c)
			if err != nil {
				tt.Error(err)
			}
		})
	}
}
