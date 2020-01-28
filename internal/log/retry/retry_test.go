package retry

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type test struct {
		name string
		opts []Option
		want *retry
	}

	tests := []test{
		func() test {
			return test{
				name: "returns retry object when options is empty",
				want: &retry{
					warnFn:  nopFunc,
					errorFn: nopFunc,
				},
			}
		}(),

		func() test {
			errorFn := func(vals ...interface{}) {}

			return test{
				name: "returns retry object when WithError options is on",
				opts: []Option{
					WithError(errorFn),
				},
				want: &retry{
					warnFn:  nopFunc,
					errorFn: errorFn,
				},
			}
		}(),

		func() test {
			warnFn := func(vals ...interface{}) {}

			return test{
				name: "returns retry object when WithWarn options is on",
				opts: []Option{
					WithWarn(warnFn),
				},
				want: &retry{
					warnFn:  warnFn,
					errorFn: nopFunc,
				},
			}
		}(),

		func() test {
			warnFn := func(vals ...interface{}) {}
			errorFn := func(vals ...interface{}) {}

			return test{
				name: "returns retry object when WithError and WithWarn options is on",
				opts: []Option{
					WithWarn(warnFn),
					WithError(errorFn),
				},
				want: &retry{
					warnFn:  warnFn,
					errorFn: errorFn,
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := New(tt.opts...).(*retry)
			if !ok {
				t.Errorf("type is invalid")
			}

			if reflect.ValueOf(got.errorFn).Pointer() != reflect.ValueOf(tt.want.errorFn).Pointer() {
				t.Error("errorfn is not equals")
			}

			if reflect.ValueOf(got.warnFn).Pointer() != reflect.ValueOf(tt.want.warnFn).Pointer() {
				t.Error("warnfn is not equals")
			}
		})
	}
}
