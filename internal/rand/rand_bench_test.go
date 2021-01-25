package rand

import "testing"

func Benchmark_Uint32(b *testing.B) {
	type args struct {
	}
	type test struct {
		name       string
		args       args
		beforeFunc func()
		afterFunc  func()
	}

	tests := []test{
		{
			name: "test rand",
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}

			for i := 0; i < b.N; i++ {
				Uint32()
			}
		})
	}
}
