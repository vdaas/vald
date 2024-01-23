package pogreb

import (
	"reflect"
	"testing"
	"time"

	"github.com/akrylysov/pogreb"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithPath(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		want Option
		err  error
		path string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *db, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *db, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if err == nil {
			if !reflect.DeepEqual(got.path, w.path) {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
		}
		return nil
	}
	tests := []test{
		func() test {
			path := "db-dir"
			return test{
				name: "Succeeds to apply option",
				args: args{
					path: path,
				},
				want: want{
					path: path,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := new(db)
			err := WithPath(test.args.path)(got)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBackgroundSyncInterval(t *testing.T) {
	t.Parallel()
	type args struct {
		s string
	}
	type want struct {
		want Option
		err  error
		opts *pogreb.Options
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *db, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *db, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if err == nil {
			if !reflect.DeepEqual(got.opts.BackgroundSyncInterval, w.opts.BackgroundSyncInterval) {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
		}
		return nil
	}
	tests := []test{
		func() test {
			dur := "100ms"
			return test{
				name: "Succeeds to apply option",
				args: args{
					s: dur,
				},
				want: want{
					opts: &pogreb.Options{
						BackgroundSyncInterval: 100 * time.Millisecond,
					},
				},
			}
		}(),
		func() test {
			dur := "invalid"
			return test{
				name: "Fails to apply option with invalid value",
				args: args{
					s: dur,
				},
				want: want{
					err: errors.New("time: invalid duration \"invalid\""),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := new(db)
			err := WithBackgroundSyncInterval(test.args.s)(got)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBackgroundCompactionInterval(t *testing.T) {
	t.Parallel()
	type args struct {
		s string
	}
	type want struct {
		want Option
		err  error
		opts *pogreb.Options
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *db, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *db, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if err == nil {
			if !reflect.DeepEqual(got.opts.BackgroundCompactionInterval, w.opts.BackgroundCompactionInterval) {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
		}
		return nil
	}
	tests := []test{
		func() test {
			dur := "100ms"
			return test{
				name: "Succeeds to apply option",
				args: args{
					s: dur,
				},
				want: want{
					opts: &pogreb.Options{
						BackgroundCompactionInterval: 100 * time.Millisecond,
					},
				},
			}
		}(),
		func() test {
			dur := "invalid"
			return test{
				name: "Fails to apply option with invalid value",
				args: args{
					s: dur,
				},
				want: want{
					err: errors.New("time: invalid duration \"invalid\""),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := new(db)
			err := WithBackgroundCompactionInterval(test.args.s)(got)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
