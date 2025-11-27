package status

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
)

func TestToCode(t *testing.T) {
	type args struct {
		code codes.Code
		err  error
	}
	type want struct {
		want codes.Code
	}
	type test struct {
		name      string
		args      args
		want      want
		checkFunc func(want, codes.Code) error
	}
	defaultCheckFunc := func(w want, got codes.Code) error {
		if got != w.want {
			return errors.Errorf("got: %s, want: %s", got.String(), w.want.String())
		}
		return nil
	}

	tests := []test{
		{
			name: "returns original code if not OK",
			args: args{
				code: codes.NotFound,
				err:  errors.New("some error"),
			},
			want: want{
				want: codes.NotFound,
			},
		},
		{
			name: "returns OK if code is OK and err is nil",
			args: args{
				code: codes.OK,
				err:  nil,
			},
			want: want{
				want: codes.OK,
			},
		},
		{
			name: "resolves code from error if code is OK",
			args: args{
				code: codes.OK,
				err:  New(codes.InvalidArgument, "invalid argument").Err(),
			},
			want: want{
				want: codes.InvalidArgument,
			},
		},
		{
			name: "resolves canceled from context error",
			args: args{
				code: codes.OK,
				err:  context.Canceled,
			},
			want: want{
				want: codes.Canceled,
			},
		},
		{
			name: "resolves deadline exceeded from context error",
			args: args{
				code: codes.OK,
				err:  context.DeadlineExceeded,
			},
			want: want{
				want: codes.DeadlineExceeded,
			},
		},
		{
			name: "returns Unknown for generic error",
			args: args{
				code: codes.OK,
				err:  errors.New("generic error"),
			},
			want: want{
				want: codes.Unknown,
			},
		},
		{
			name: "returns Unknown for invalid code",
			args: args{
				code: codes.Code(100),
				err:  nil,
			},
			want: want{
				want: codes.Unknown,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			checkFunc := test.checkFunc
			if checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := ToCode(test.args.code, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
