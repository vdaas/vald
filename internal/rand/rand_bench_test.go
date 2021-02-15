package rand

import (
	"strconv"
	"testing"
)

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
	b.ResetTimer()

	type paralleTest struct {
		name       string
		args       args
		paralle    []int
		beforeFunc func()
		afterFunc  func()
	}
	ptests := []paralleTest{
		{
			name:    "test rand",
			paralle: []int{1, 2, 4, 6, 8, 16},
		},
	}
	for _, ptest := range ptests {
		test := ptest
		for _, p := range test.paralle {
			name := test.name + "-" + strconv.Itoa(p)
			b.Run(name, func(b *testing.B) {
				b.SetParallelism(p)
				b.ResetTimer()

				if test.beforeFunc != nil {
					test.beforeFunc()
				}
				if test.afterFunc != nil {
					defer test.afterFunc()
				}
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						Uint32()
					}
				})
			})
		}
	}
}
