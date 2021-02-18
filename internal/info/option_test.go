package info

import (
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

func TestWithServerName(t *testing.T) {
	type T = info
	type args struct {
		s string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set server name success",
			args: args{
				s: "srvName",
			},
			want: want{
				obj: &T{
					detail: Detail{
						ServerName: "srvName",
					},
				},
			},
		},
		{
			name: "return ErrInvalidOption error when server name is empty",
			args: args{
				s: "",
			},
			want: want{
				obj: &T{
					detail: Detail{},
				},
				err: errors.NewErrInvalidOption("ServerName", ""),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithServerName(test.args.s)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithRuntimeCaller(t *testing.T) {
	type T = info
	type args struct {
		f func(skip int) (pc uintptr, file string, line int, ok bool)
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			comparator.Comparer(func(x, y sync.Once) bool {
				return reflect.DeepEqual(x, y)
			}),
			comparator.Comparer(func(x, y func(skip int) (pc uintptr, file string, line int, ok bool)) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
			comparator.Comparer(func(x, y func(pc uintptr) *runtime.Func) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
		}
		if diff := comparator.Diff(w.obj, obj, opts...); len(diff) != 0 {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	tests := []test{
		func() test {
			f := func(skip int) (pc uintptr, file string, line int, ok bool) {
				return uintptr(0), "", 0, false
			}
			return test{
				name: "set func success",
				args: args{
					f: f,
				},
				want: want{
					obj: &T{
						rtCaller: f,
						detail:   Detail{},
					},
				},
			}
		}(),
		{
			name: "return ErrInvalidOption error when func is nil",
			args: args{
				f: nil,
			},
			want: want{
				obj: &T{
					detail: Detail{},
				},
				err: errors.NewErrInvalidOption("RuntimeCaller", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithRuntimeCaller(test.args.f)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithRuntimeFuncForPC(t *testing.T) {
	type T = info
	type args struct {
		f func(pc uintptr) *runtime.Func
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			comparator.Comparer(func(x, y sync.Once) bool {
				return reflect.DeepEqual(x, y)
			}),
			comparator.Comparer(func(x, y func(skip int) (pc uintptr, file string, line int, ok bool)) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
			comparator.Comparer(func(x, y func(pc uintptr) *runtime.Func) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
		}
		if diff := comparator.Diff(w.obj, obj, opts...); len(diff) != 0 {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	tests := []test{
		func() test {
			f := func(pc uintptr) *runtime.Func {
				return nil
			}
			return test{
				name: "set func success",
				args: args{
					f: f,
				},
				want: want{
					obj: &T{
						rtFuncForPC: f,
						detail:      Detail{},
					},
				},
			}
		}(),
		{
			name: "return ErrInvalidOption error when func is nil",
			args: args{
				f: nil,
			},
			want: want{
				obj: &T{
					detail: Detail{},
				},
				err: errors.NewErrInvalidOption("RuntimeFuncForPC", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithRuntimeFuncForPC(test.args.f)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
