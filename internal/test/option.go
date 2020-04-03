package test

import (
	"context"
	"testing"
)

type Option func(*test)

var (
	defaultOptions = []Option{}
)

func WithCase(cs ...Caser) Option {
	return func(t *test) {
		if len(cs) != 0 {
			t.cs = cs
		}
	}
}

func WithTarget(
	fn func(ctx context.Context,
		t *testing.T,
		args, fields []interface{},
		checkFunc func(t *testing.T, gots ...interface{}),
	),
) Option {
	return func(t *test) {
		if fn != nil {
			t.target = fn
		}
	}
}
