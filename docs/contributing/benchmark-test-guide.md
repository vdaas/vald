# Benchmark Test Guideline

## Introduction

### What is the problem now in Vald?

In Vald, performance is a critical requirement and factor for users. If Vald does not perform good and stable performance, user may leave Vald.

So it is important to capture the performance statistic for Vald to perform performance tunning and help to find the bottleneck of Vald when there is some performance issue in Vald.

### How can we resolve this problem?

Benchmarking is used to measure the performance of target elements to produce a performance metrics for future performance investigation.

We can measure the performance statistic in 3 different levels to measure different level of performance:

1. Code benchmarking
  - to capture performance metrics by function level (/internal)
  - to perform lowest level of performance investigation/tuning

1. Code level E2E benchmarking
  - to capture performance metrics by use case level (/pkg)

1. Component level E2E benchmarking
  - to capture performance metrics by component level

Unlike unit testing, benchmark test only focus of measure on the performance of the function but not the functionality of the function.

### What can we do with benchmark?

In benchamrking, we can measure performance using specific indicator and resulting in a metrics of performance.

The metrics can be used in different purpose:

- demonstrate whether or not a system meets the criteria set forward
- compare two applications to determine which one works better
- measure a system to find what performs badly

Ref: https://www.castsoftware.com/glossary/software-performance-benchmarking-modeling#:~:text=Software%20performance%20benchmarking%20serves%20different,to%20find%20what%20performs%20badly.

Here is the example of what to measure in benchmarking:

- How much time consumed for each operation (e.g. time/operation)
- How much memory allocated for each operation (e.g. ??? byte/operation)
- How many times the memory is allocated (e.g. ??? times)

After getting the metrics, we can ask ourself the following questions about the metrics:

- How should this software perform under the maximum load?
- How many concurrent users are expected on a daily basis?
- What does “good” performance look like?
- What does “acceptable” performance look like?

### How should we work on benchmark?

- Work on every functions?
- With template
	- Try few important packages first to test if the template is working
- Without template

### Schedule

(If we use the template)
- Test the template with difficult packages (1-2weeks)
- Prioritize the order to implementing bench code

### Questions/Concerns

1. Missing code benchmark coverage in golang
1. How to detect the changes on benchmark result? 





## Code bench

Code bench is the lowest level of benchmarking in above. It performs the benchmark testing on functions and it provide useful statistic for investigation.

In golang, it supports the get the benchmark metrics by using `go` command.

We can measure the CPU and the memory usage/allocation information of the target function.

e.g. The result is as follow.

```
Benchmark_Uint32/test_rand-4         	17290003	        82.3 ns/op	       0 B/op	       0 allocs/op
```

Reference: https://golang.org/pkg/testing/#BenchmarkResult

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
