package safety

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func BenchmarkRecoverFunc(b *testing.B) {
	type args struct {
		fn func() error
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "fn return error func",
			args: args{
				fn: func() error {
					return errors.New("fn err")
				},
			},
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := RecoverFunc(test.args.fn)
					if got == nil {
						b.Error("got is empty")
					}
				}
			})
		})
	}
}

func BenchmarkRecoverWithoutPanicFunc(b *testing.B) {
	type args struct {
		fn func() error
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "fn return error func",
			args: args{
				fn: func() error {
					return errors.New("fn err")
				},
			},
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := RecoverWithoutPanicFunc(test.args.fn)
					if got == nil {
						b.Error("got is empty")
					}
				}
			})
		})
	}
}

func BenchmarkRecoverWithoutPanicFunc_with_execution(b *testing.B) {
	type args struct {
		fn func() error
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "fn return error func",
			args: args{
				fn: func() error {
					return errors.New("fn err")
				},
			},
		},
		{
			name: "fn panic runtime error func",
			args: args{
				fn: func() error {
					_ = []string{}[10]
					return nil
				},
			},
		},
		{
			name: "fn panic string func",
			args: args{
				fn: func() error {
					panic("panic")
				},
			},
		},
		{
			name: "fn panic error func",
			args: args{
				fn: func() error {
					panic(errors.Errorf("error"))
				},
			},
		},
		{
			name: "fn panic int func",
			args: args{
				fn: func() error {
					panic(10)
				},
			},
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := RecoverWithoutPanicFunc(test.args.fn)
					if got == nil {
						b.Error("got is empty")
					}
					err := got()
					if err == nil {
						b.Error("err is nil")
					}
				}
			})
		})
	}
}
