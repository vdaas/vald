package version

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestCheck(t *testing.T) {
	type args struct {
		cur string
		max string
		min string
	}

	type test struct {
		name string
		args args
		want error
	}

	tests := []test{
		{
			name: "return nil",
			args: args{
				cur: "1.0.5",
				max: "1.0.10",
				min: "1.0.0",
			},
			want: nil,
		},

		{
			name: "return error when cur format is invalid",
			args: args{
				cur: "vald",
				max: "1.0.10",
				min: "1.0.0",
			},
			want: errors.New("Malformed version: vald"),
		},

		{
			name: "return error when min format is invalid",
			args: args{
				cur: "1.5.10",
				max: "vald",
				min: "1.0.0",
			},
			want: errors.New("Malformed constraint:  <= vald"),
		},

		{
			name: "return error when min format is invalid",
			args: args{
				cur: "1.0.10",
				max: "1.0.5",
				min: "1.0.0",
			},
			want: errors.ErrInvalidConfigVersion("1.0.10", ">= 1.0.0, <= 1.0.5"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Check(tt.args.cur, tt.args.max, tt.args.min)
			if !errors.Is(tt.want, err) {
				t.Error(err)
			}
		})
	}

}
