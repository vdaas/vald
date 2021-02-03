package errors

import (
	"testing"
)

func TestErrParseUnitFailed(t *testing.T) {
	type args struct {
		str string
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
			return Errorf("got: \"%#v\", \n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrParseUnitFailed error when str is not empty.",
			args: args{
				str: "parse target string",
			},
			want: want{
				want: New("failed to parse: 'parse target string'"),
			},
		},
		{
			name: "return an ErrParseUnitFailed error when str is empty.",
			args: args{},
			want: want{
				want: New("failed to parse: ''"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrParseUnitFailed(test.args.str)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
