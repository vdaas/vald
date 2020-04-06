package test

import (
	"context"
	"reflect"
	"testing"
)

type Test interface {
	Run(context.Context, *testing.T)
}

type test struct {
	cs     []Caser
	target func(ctx context.Context, c Caser) ([]interface{}, error)
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

			gots, err := test.target(ctx, c)
			if err != nil {
				tt.Error(err)
			}

			if fn := c.CheckFunc(); fn != nil {
				if err := fn(gots, c.Wants()); err != nil {
					t.Errorf("checkFunc returns error: %v", err)
				}
			} else {
				if len(c.Wants()) != len(gots) {
					tt.Fatalf("wants and gots length are not equals. wants: %d, gots: %d", len(c.Wants()), len(gots))
				}

				for i, want := range c.Wants() {
					if !reflect.DeepEqual(want, gots[i]) {
						tt.Errorf("not equals. want: %v, but got: %v", want, gots[i])
					}
				}
			}
		})
	}
}
