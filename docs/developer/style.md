# Go Style Guide in Vald

## Table of Contents

- [Introduction](#Introduction)
- [Style](#Style)
- [Test](#Test)
	- [Table-Driven-Test](#Table-Driven-Test)

## Introduction

## Style

## Test

Testing guideline has 2 important rules for the coding quality and readability
1. Use Table-Driven-Test
2. Keep code coverage over 85%
   - test coverage != high testing quality, but low coverage means bad testing quality
   - check with the following commands `go test -cover ./...`
3. Test all use cases and error cases

### Table-Driven-Test

Use table-driven tests with subtests to avoid duplicating code. 

<table>
<thead><tr><th>Bad</th><th>Good</th><th>In Vald</th></tr></thead>
<tbody>
<tr>
<td>
	
```go
## case 1
host, port, err := net.SplitHostPort("192.0.2.0:8000")
if err != nil {
    t.Errorf("error is not nil: %v", err)
}

if want, got := "192.0.2.0", host; want != got {
    t.Errorf("host is not equals. want: %s, but got: %s", want: %s, got)
}

if want, got := "8000", port; want != got {
    t.Errorf("port is not equals. want: %s, but got: %s", want: %s, got)
}

## case2
host, port, err = net.SplitHostPort("192.0.2.0:http")
if err != nil {
    t.Errorf("error is not nil: %v", err)
}

if want, got := "192.0.2.0", host; want != got {
    t.Errorf("host is not equals. want: %s, but got: %s", want: %s, got)
}

if want, got := "http", port; want != got {
    t.Errorf("port is not equals. want: %s, but got: %s", want: %s, got)
}
```

</td>
<td>

```go
tests := []struct {
    str string
    wantHost string
    wantPort string
} {
    ## case 1
    {
        str: "192.0.2.0:8000",
        wantHost: "192.0.2.0",
        wantPort: "8000",
    },
    ## case 2
    {
    	str: "192.0.2.0:http",
    	wantHost: "192.0.2.0",
    	wantPort: "http",
	},
}

for _, tt := range tests {
    t.Run(tt.str, func(tt *testing.T) {
        host, port, err := net.SplitHostPort(tt.str)
        if err != nil {
            t.Errorf("error is not nil: %v", err)
        }
        if want, got := tt.wantHost, host; want != got {
            t.Errorf("host is not equals. want: %s, but got: %s", want: %s, got)
        }
        if want, got := tt.wantPort, port; want != got {
            t.Errorf("port is not equals. want: %s, but got: %s", want: %s, got)
        }
    })
}
```

</td>
<td>

```go
test.New(
    test.WithCase(
        caser.New(
            WithName("case_1"),
            WithArg("192.0.2.0:8000")
            WithWant("192.0.2.0", "8000", nil),
        ),

        caser.New(
            WithName("case_2"),
            WithArg("192.0.2.0:http")
            WithWant("192.0.2.0", "http", nil),
        ),
    ),
    test.WithDriverFunc(
        func(ctx context.Context, d test.DataProvider) []interface{} {
            host, port, err := net.SplitHostPort(d.Args()[0].(string))
            return []interface{} {
                host, port, err
            }
        }
    ),
).Run(context.Background(), t)
```

</td>


</tr>
</tbody>
</table>

Table-Driven-Test makes it easy to add new test case.


We define the test case table as `map[string]func(*testing.T)test`, which is referred as the test case name and the test case implementation `tt`. 
```go
tests := map[string]func(t *testing.T) test {
    "test case name": func(tt *testing.T) test {
        return test {
            args: args {
                host: "host",
                port: "port",
            },
            field: field {
                timeout: 1 * time.Second,
            },
        }
    }
}
```

### The steps to create a Table-Driven-Test.

1. `args` structure

If there are two or more arguments to be passed to the method, create a `args` structure. If there is only one argument, do not create an `args` structure.
```go
type args struct {
    host string
    port string
}
```


2. `field` structure

If you create an object and test its methods, create a `field` struct if the object has two or more fields to initialize. If there is only one field, do not create `field` structure.

```go
type field struct {
    host string
    port string
}
```

3. `test` structure

`test` structure has `args` and `field` structure and `checkFunc` function. If you need one of `args` and `field` structure, create `field` and `args` structure.
The `checkFunc` function is used to check the return value of the function being tested.

```go
type test struct {
    args args
    field field
    checkFunc func(t *testing.T, err error)
}
```


Example:

```go

type args struct {
    addr string
    txt string
}

type field struct {
    timeout time.Duration
}

type test struct {
    args args
    field field
    checkFunc func(t *testing.T, err error)
}

tests := map[string]func(*testing.T) test {
    "test name": func(tt *testing.T) test {
        tt.Helper()

        return test {
            args: args {
                host: "host",
                port: "port",
            },
            field: field {
                host: "host",
                port: "port",
            },
            checkFunc func(tt *testing.T, err error) {
                t.Helper()
                if err != nil {
                    tt.Errorf("error is not nil: %v", err)
                }
            },
        }
    } 
}

for name, fn := range tests {
    t.Run(name, func(tt *tesint.T) {
        test := fn(tt)

        c := client {
            timeout: test.field.timeout,
        }

        err := c.Send(test.args.addr, test.args.txt)
        test.checkFunc(tt, err)
    })
}

```
