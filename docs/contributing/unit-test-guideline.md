# Unit Test Guideline

## About

When we create the unit test for each function or package, we should consider test coverage not only code coverage.
This guideline will help you to create good unit tests.

## Unit Test

Before considering the Unit test, we should define what is the unit.
Unit is a set of procedures or functions, in a procedural or functional language.
Also, it is a class and its needed classes, in an object or object-oriented language.
Unit testing validates the basic units of the program in isolation.

## Guideline

### Strategy

Before creating the unit tests, we must know the essentials of strategy.

1. Know what you are testing
	- purpose / consumers of target function or classess or etc / designed by contract  defensive program
1. Test behaviors and results, not implementation
	- ensures that tests only fail when there is an actual effect and not due to internal changes
1. Test one thing at a time
	- each test should have a clear, concise, singular objective
1. Make tests readable and understandable
	- hold test code to a similar standard as production code
1. Make tests deterministic
	- A test should pass all the time or fail all the time until fixed.
1. Make tests independent and self-sufficient
	- Setup, execution, and verification steps in a given test should not depend on running othre tests before it. To keep unit tests simple, fast running, and easy to debug, it may be necessary to isolate the class under test
1. Repeat yourself when necessary
	- it is okay to violate the 'do not repeat yourself' principle if it makes tests simpler and easier to read
1. Measure code covearage but focus on test coverage
	- Do not simply to archive code coverage

### Test case

#### Basic

It is not perfect, but we should try cover all codes at first.
To cover all code, the basic way is the creating test cases to complete all branch in target unit.

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

When the above function is given, we should create 2 test for archiving 100% code coverage at least.

- When val is not given
  - In other word, the default value is given
- When val(int32) is given
  - `var val int32 = 1`
  
It seems enough for the given function, but we have to take care some test cases are remaining for imporoving test coverage.
In this case, there is one test case is remaining.

- When val([]int32) is given
  - `var val []int32 = {1, 2}`
  
That is the focusing on test coverage.

Therefore, we should concern all cases for impoving test coverage of target unit.
To improve test coverage, the basic but critical thinking way is think about input patterns.
It is not only single input, but also multi inputs.

If a function or method requires multi input, we should try to create test case many patterns.

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

|len(val1)|len(val2| option |
|:-----:|:-----:|:-----:|
| 0  | 0  | - |
| 0  | >0 | - |
| >0 | 0  | - |
| >0 | >0 | ave1 > ave2 |
| >0 | >0 | ave1 < ave2 |

At a glance, one of the last 2 patterns is enough, but these will help us to notice the bug in dependency.
It will avoid the unexpected error due to update dependencies.

You should create unit tests for error pattern as same as succesable pattern.

#### Advanced

##### Robust boudary test

The previous section is about a basic test case.
For more cover test coverage, the (robust) boundary test should be applied.
A boundary test is testing values on or near the boundaries of allowed inputs.

```go
func AgeTest(val int32) (ok bool) {
	if val >= 20 && val < 65 {
		ok = true
	}
	return
}
```

When the above function is given as target, the minimum number of tests based on the "Robust Boundary" test is 7.
  - the input `val` cases: 19, 20, 21, 30, 64, 65, 66
  
The robust boundary test requires 6N+1 test cases (N is the number of inputs for target function).
More increasing the required input, more increasing the minimum number of cases.
Sometimes, it requires a lot of resources to create/maintain test.
So, you may test the critical cases using boundary value and have a sense of purpose.

#### Equivalence Class Testing

Equivalence Class Testing is grouping the values together such that all members of the group are considered the same and testing one value from each.
```go
func AgeTest(val int32) (ok bool) {
	if val >= 20 && val < 65 {
		ok = true
	}
	return
}
```

When the above function is given as target, you have to create 3 groups and pick one value from each group.

|Desc|Input range|Chose Input|
|:-----:|:-----:|:-----:|
| Under 20    | {0, ...,  19 }| 5  |
| worker range| {20, ..., 64} | 29 |
| Senior      | {65, ...    } | 72 |

### Vald Style

In the Vald, we create unit tests based on the basic test case.
And, you also create unit tests based on robust boundary tests or equivalence class tests as needed.

## Coding Style

Please refer [here](../coding-style.md#Test)
