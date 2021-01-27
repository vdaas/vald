package errors

import (
	"testing"
)

func TestErrMetaDataAlreadyExists(t *testing.T) {
	t.Parallel()
	type args struct {
		meta string
	}
	type want struct {
		want error
	}
	type test struct{
		name string
		args args
		want want
		checkFunc func(want, error) error
		beforeFunc func(args)
		afterFunc func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "return an ErrMetaDataAlreadyExists error when meta is not empty.",
			args: args{
				meta: "vald-meta-01",
			},
			want: want{
				want: New("vald metadata:\tvald-meta-01\talready exists "),
			},
		},
		{
			name: "return an ErrMetaDataAlreadyExists error when meta is empty.",
			args: args{},
			want: want{
				want: New("vald metadata:\t\talready exists "),
			},
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrMetaDataAlreadyExists(test.args.meta)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})

	}
}

func TestErrMetadataCannotFetch(t *testing.T) {
	t.Parallel()
	type want struct {
		want error
	}
	type test struct{
		name string
		want want
		checkFunc func(want, error) error
		beforeFunc func()
		afterFunc func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "return an ErrMetaDataCannotFetch error",
			want: want{
				want: New("vald metadata cannot fetch"),
			},
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrMetaDataCannotFetch()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})

	}
}
