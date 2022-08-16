# Unit Test Guideline

## About

The goal of the unit test is to check whether the unit is implemented correctly in various cases.
A good unit test is checking various cases we do not expect as supposed usage.
When we create the unit test for each function or package, we have to consider test coverage, not only code coverage.
This guideline will help you to create good unit tests.

## Unit Test

Before considering the unit test, we should define what the unit is.
A unit is a set of procedures or functions in a procedural or functional language.
Also, it is a class, and its needed classes, in an object or object-oriented language.
Unit testing validates the basic units of the program in isolation.

## Guideline

### Strategy

Before creating the unit tests, we must know the essentials of strategy.

1. Know what you are testing
   - Purpose/Consumers of target function or classes or etc / designed by defensive contract program.
1. Test behaviors and results, not implementation
   - Ensures that tests only fail when there is an actual effect and not due to internal changes.
1. Test one thing at a time
   - Each test should have a clear, concise, singular objective.
1. Make tests readable and understandable
   - Hold test code to a similar standard as production code.
1. Make tests deterministic
   - A test should pass all the time or fail all the time until fixed.
1. Make tests independent and self-sufficient
   - Setup, execution, and verification steps in a given test should not depend on running other tests. To keep unit tests simple, fast running, and easy to debug, it may be necessary to isolate the class under test.
1. Repeat yourself when necessary
   - It is okay to violate the 'do not repeat yourself' principle if it makes tests simpler and easier to read.
1. Measure code coverage but focus on test coverage
   - Do not simply archive code coverage.

### Test case

#### Basic

It is not perfect for improving test coverage, but we should first try to cover all codes.
To cover all code, the basic way is to create test cases to complete all branches in the target unit.

Let's see the below function.

```go
func calcSum(val ...int32) (sum int32) {
	if len(val) == 0 {
		return
	} else {
		for _, v := range val {
			sum += v
		}
		return
	}
}
```

When the above function is given, we should create 2 tests to achieve 100% code coverage at least.

- When val is not given
  - In other words, the default value is given
- When val(int32) is given
  - `var val int32 = 1`

It seems enough for the given function, but we have to take care of some remaining test cases to improve test coverage.
There is at least one more test case to complete test coverage in this case.

- When val([]int32) is given
  - `var val []int32 = {1, 2}`
  - `var val []int32 = {math.MaxInt32, -1}`
  - `var val []int32 = {math.MaxInt32, 0}`
  - `var val []int32 = {math.MaxInt32, 1}`
  - `var val []int32 = {math.MinInt32, -1}`
  - `var val []int32 = {math.MinInt32, 0}`
  - `var val []int32 = {math.MinInt32, -1}`
  - `var val []int32 = {math.MaxInt32, math.MaxInt32}`
  - `var val []int32 = {math.MinInt32, math.MinInt32}`

That is the focus on test coverage.

Therefore, we should consider all cases to improve test coverage for the target unit.
To enhance test coverage, the basic but critical thinking way is thinking about input patterns.
It is not only a single input, but also multi inputs.

If a function or method accepts multiple inputs, we should try to create test cases to cover all the inputs.

(Note: The below `calcSum()` is diffrent function from `calsSum(val ...int32)` which we mention before.)

```go
func calcAverageDiff(val1 []int32, val2 []int32) (diff float64) {
	var ave1, ave2 float64
	if len(val1) == 0 && len(val2) == 0 {
		return
	}
	if len(val1) != 0 {
		ave1 = float64(calcSum(val1)) / float64(len(val1))
	}
	if len(val2) != 0 {
		ave2 = float64(calcSum(val2)) / float64(len(val2))
	}

	diff = math.Abs(ave1 - ave2)
	return
}
```

When `calcAverageDiff()` is given, the test patterns are below:

| len(val1) | len(val2) | option                                            |
| :-------: | :-------: | :------------------------------------------------ |
|     0     |     0     | -                                                 |
|     0     |    >0     | -                                                 |
|    >0     |     0     | -                                                 |
|    >0     |    >0     | ave1 > ave2                                       |
|    >0     |    >0     | ave1 < ave2                                       |
|    >0     |    >0     | ave1 = ave2                                       |
|    >0     |    >0     | ave1 = {math.MinInt32}, ave2 = {math.MaxInt32}    |
|    >0     |    >0     | ave1 = {math.MaxInt32, math.MaxInt32}, ave2 = {0} |

At a glance, one of the last 2 patterns is enough, but these will help us notice the dependency bug.
It will avoid unexpected errors due to updated dependencies.

You have to create unit tests for error patterns as the same as success patterns.

#### Advanced

##### Robust boudary test

The previous section is about the basic test cases.
The (robust) boundary test should be applied to cover more test coverage.
A boundary test tests values on or near the boundaries of allowed inputs.

```go
func AgeTest(val int32) (ok bool) {
	if val >= 20 && val < 65 {
		ok = true
	}
	return
}
```

When the above function is given as a target, the minimum number of tests based on the "Robust Boundary" test is 7.

- the input `val` cases: 19, 20, 21, 30, 64, 65, 66

The robust boundary test requires 6N+1 test cases (N is the number of inputs for the target function).
Increasing the input arguments will also increase the number of required test cases.
Sometimes, it requires a lot of resources to create/maintain tests.
So, you may test the critical cases using boundary values and have a sense of purpose.

#### Equivalence Class Testing

Equivalence Class Testing groups the values together such that all members of the group are considered the same and tests one value from each.

```go
func AgeTest(val int32) (ok bool) {
	if val >= 20 && val < 65 {
		ok = true
	}
	return
}
```

When the above function is given as a target, you have to create three groups and pick one value from each group.

|     Desc     |  Input range  | Chose Input |
| :----------: | :-----------: | :---------: |
|   Under 20   | {0, ..., 19 } |      5      |
| worker range | {20, ..., 64} |     29      |
|    Senior    |  {65, ... }   |     72      |

### Vald Style

In the Vald, we create unit tests based on the basic test case.
And, you may also create unit tests based on robust boundary tests or equivalence class tests if needed.

But, we have to take care the Vald is developed using Go.
As you know as Go has many coding features as other languages.
One of the features is that Go will convert a single value to a slice value when the Function or Method receives a variadic argument (e.g., `...[]int`, `...[]string`, `...interface{}`, or etc.) as the input.

And we apply the table-driven test for running unit tests.
For example, when we create the unit test of `func getMeta(...[]int)`, the test code will be more complex than other functions' tests which don't use variadic argument as the input, if we create the test for all input patterns.
So, we define the basic unit case slightly different from [the basic test case](#Basic).

This change is very clear and you can apply it easily.
Our basic test case depends on the type of two variadic arguments.

1. When input is `...interface{}`
   - We have to write all test cases which satisfy `...interface{}` as same as [a basic test case](#Basic). For example, `val = 1`, `val = "input"`, `val = []float64{2020.12}` and so on.
1. When input is not `...interface{}` but `...[]int`, `...[]string` or etc.
   - We have to create only slice pattern test cases, which is the same as not creating test cases with a single value.
   - We should test with boundary cases, for example, we should test with `val = []int{math.MaxInt64()}` when the input value is `...[]int`.

Summarize Vald unit test guideline:

- Apply basic test case, but take care of input variable pattern, in particular, the variadic argument (`...interface{}` or not)
- Apply robust boundary tests, including edge cases (e.g., `math.MaxInt64()`)
- Apply equivalence class testing when needed.

## Coding Style

Please refer [here](../contributing/coding-style.md#Test)
