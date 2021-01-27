# Benchmark Test Guideline

## What is the problem now?

In Vald, performance is a critical requirement and factor for users. If Vald does not perform good and stable performance, user may leave Vald.

So it is important to capture the performance statistic for Vald to measure the performance of Vald and help to find the bottleneck of Vald.

## How can we resolve this problem?

Benchmarking is used to measure the performance of target elements to produce a performance metrics for future performance investigation.

We can measure the performance statistic in 3 different levels:

1. Code bench
1. Code level E2E
1. Component level E2E

Unlike unit testing, benchmark test only focus of measure on the performance of the function but not the functionality of the function.

## Code bench

Code bench is the lowest level of benchmarking in above. It performs the benchmark testing on functions and it provide useful statistic for investigation.

In golang, it supports the get the benchmark metrics by using `go` command.

We can measure the CPU and the memory usage/allocation information of the target function.

e.g. The result is as follow.

```
Benchmark_Uint32/test_rand-4         	17290003	        82.3 ns/op	       0 B/op	       0 allocs/op
```

### How to write bench code?

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
