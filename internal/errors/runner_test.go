package errors

import (
	"testing"

	"go.uber.org/goleak"
)

func TestErrDaemonStartFailed(t *testing.T) {
	t.Parallel()
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			err := New("runner error")
			return test{
				name: "return an ErrDaemonStartFailed error when err is not nil",
				args: args{
					err: err,
				},
				want: want{
					want: Wrap(err, "failed to start daemon"),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrDaemonStartFailed error when err is nil",
				want: want{
					want: Wrap(nil, "failed to start daemon"),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrDaemonStartFailed(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrDaemonStopFailed(t *testing.T) {
	t.Parallel()
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			err := New("runner error")
			return test{
				name: "return an ErrDaemonStopFailed error when err is not nil",
				args: args{
					err: err,
				},
				want: want{
					want: Wrap(err, "failed to stop daemon"),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrDaemonStopFailed error when err is nil",
				want: want{
					want: Wrap(nil, "failed to stop daemon"),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrDaemonStopFailed(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrStartFunc(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
		err  error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			err := New("runner Start error")
			name := "gateway"
			return test{
				name: "return an ErrStartFunc error when err is not nil and name is not empty",
				args: args{
					name: name,
					err:  err,
				},
				want: want{
					want: Wrapf(err, "error occurred in runner.Start at %s", name),
				},
			}
		}(),
		func() test {
			err := New("runner Start error")
			var name string
			return test{
				name: "return an ErrStartFunc error when err is not nil and name is empty string",
				args: args{
					err: err,
				},
				want: want{
					want: Wrapf(err, "error occurred in runner.Start at %s", name),
				},
			}
		}(),
		func() test {
			name := "gateway"
			return test{
				name: "return an ErrStartFunc error when err is nil and name is not empty",
				args: args{
					name: name,
				},
				want: want{
					want: Wrapf(nil, "error occurred in runner.Start at %s", name),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrStartFunc(test.args.name, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrPreStopFunc(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
		err  error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runner PreStop error")
	defaultName := "gateway"
	tests := []test{
		func() test {
			return test{
				name: "return an ErrPreStopFunc error when err is not nil and name is not empty",
				args: args{
					name: defaultName,
					err:  defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.PreStop at %s", defaultName),
				},
			}
		}(),
		func() test {
			var name string
			return test{
				name: "return an ErrPreStopFunc error when err is not nil and name is empty string",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.PreStop at %s", name),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrPreStopFunc error when err is nil and name is not empty",
				args: args{
					name: defaultName,
				},
				want: want{
					want: Wrapf(nil, "error occurred in runner.PreStop at %s", defaultName),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrPreStopFunc(test.args.name, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrStopFunc(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
		err  error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runner Stop error")
	defaultName := "gateway"
	tests := []test{
		func() test {
			return test{
				name: "return an ErrStopFunc error when err is not nil and name is not empty",
				args: args{
					name: defaultName,
					err:  defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.Stop at %s", defaultName),
				},
			}
		}(),
		func() test {
			var name string
			return test{
				name: "return an ErrStopFunc error when err is not nil and name is empty string",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.Stop at %s", name),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrStopFunc error when err is nil and name is not empty",
				args: args{
					name: defaultName,
				},
				want: want{
					want: Wrapf(nil, "error occurred in runner.Stop at %s", defaultName),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrStopFunc(test.args.name, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrPostStopFunc(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
		err  error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runner PostStop error")
	defaultName := "gateway"
	tests := []test{
		func() test {
			return test{
				name: "return an ErrPostStopFunc error when err is not nil and name is not empty",
				args: args{
					name: defaultName,
					err:  defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.PostStop at %s", defaultName),
				},
			}
		}(),
		func() test {
			var name string
			return test{
				name: "return an ErrPostStopFunc error when err is not nil and name is empty string",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.PostStop at %s", name),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrPostStopFunc error when err is nil and name is not empty",
				args: args{
					name: defaultName,
				},
				want: want{
					want: Wrapf(nil, "error occurred in runner.PostStop at %s", defaultName),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrPostStopFunc(test.args.name, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRunnerWait(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
		err  error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runner RunnerWait error")
	defaultName := "gateway"
	tests := []test{
		func() test {
			return test{
				name: "return an ErrRunnerWait error when err is not nil and name is not empty",
				args: args{
					name: defaultName,
					err:  defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.Wait at %s", defaultName),
				},
			}
		}(),
		func() test {
			var name string
			return test{
				name: "return an ErrRunnerWait error when err is not nil and name is empty string",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrapf(defaultErr, "error occurred in runner.Wait at %s", name),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrRunnerWait error when err is nil and name is not empty",
				args: args{
					name: defaultName,
				},
				want: want{
					want: Wrapf(nil, "error occurred in runner.Wait at %s", defaultName),
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrRunnerWait(test.args.name, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
