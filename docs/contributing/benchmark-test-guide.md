# Benchmark Test Guideline

## About

Benchmarking is used to measure performance of target elements to produce a performance metics that used to compare to others or future performance investigation.

Unlike unit testing, benchmark test only focus of measure on the benchmark of the function. We can measure the CPU and the memory usage/allocation information of the target function.

This is the guideline to help you to implement benchmark code. Implement benchmark test code allows continous measurement on the performance.


## Which package / file to test?

We should only test on the function that is performance critical or the function that is used frequently.


## How to write bench code?

Create / generate a file called `[filename]_bench_test.go`.

For example if we want to test the `internal/rand/rand.go` file, you need to create / generate a file called `internal/rand/rand_bench_test.go` file.

To execute the benchmarking, use the `go test -bench . -benchmem` command.

We implement benchmark code for each function we want to test.
We should follow this template to write the benchmark code.

```golang
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
```
