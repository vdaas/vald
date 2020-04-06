package test

import (
	"context"
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

func WithTarget(fn func(context.Context, Caser) ([]interface{}, error)) Option {
	return func(t *test) {
		if fn != nil {
			t.target = fn
		}
	}
}
