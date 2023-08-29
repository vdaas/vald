package malloc

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_convert(t *testing.T) {
	type args struct {
		body string
	}
	type want struct {
		wantM *MallocInfo
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *MallocInfo, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotM *MallocInfo, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotM, w.wantM) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotM, w.wantM)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		// {
		// 	name: "test_case_1",
		// 	args: args{
		// 		body: "",
		// 	},
		// 	want:      want{},
		// 	checkFunc: defaultCheckFunc,
		// 	beforeFunc: func(t *testing.T, args args) {
		// 		t.Helper()
		// 	},
		// 	afterFunc: func(t *testing.T, args args) {
		// 		t.Helper()
		// 	},
		// },
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

			gotM, err := convert(test.args.body)
			if err := checkFunc(test.want, gotM, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGetMallocInfo(t *testing.T) {
	type want struct {
		want *MallocInfo
		err  error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *MallocInfo, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *MallocInfo, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "convert type check",
			want: want{
				want: &MallocInfo{},
				err:  nil,
			},
			checkFunc: func(w want, got *MallocInfo, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if reflect.TypeOf(got) != reflect.TypeOf(w.want) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := GetMallocInfo()
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
