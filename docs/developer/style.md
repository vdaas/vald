# Go Style Guide in Vald

## Table of Contents

- [Introduction](#Introduction)
- [Style](#Style)
- [Test](#Test)
	- [Table-Driven-Test](#Table-Driven-Test)

## Introduction

## Style

## Test

Testing guideline has 3 important rules for the coding quality and readability
1. Use Table-Driven-Test
2. Keep code coverage over 85%
   - test coverage != high testing quality, but low coverage means bad testing quality
   - check with the following commands `go test -cover ./...`
3. Test all use cases and error cases

### Table-Driven-Test

Table-Driven-Test makes it easy to add new test case. Also it can avoid duplicating code. 

This table presents 3 test styles. We need to follow `Vald Style`.

<table>
<thead><tr><th>Bad</th><th>Good</th><th>Vald Style</th></tr></thead>
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

To create tests based on `Vald Style`, We provide the original test framework.

Please refer to [this page](../../internal/test).

### The steps to create a Table-Driven-Test in Vald

#### Basic Style

This case can automatically check by the arguments of `WithWant` and the object returned by the function of arguments of `WithFunc`.

```go
test.New(
    test.WithCase(
        caser.New(
            WithName("case_1"),
            WithArg("192.0.2.0:8000")
            WithWant("192.0.2.0", "8000", nil), // want objects
        ),
    ),
    test.WithDriverFunc(
        func(ctx context.Context, d test.DataProvider) []interface{} {
            host, port, err := net.SplitHostPort(d.Args()[0].(string))

            // got objects
            return []interface{} {
                host, port, err
            }
        }
    ),
).Run(context.Background(), t)
```

#### Customize the Assertion

1. In-Out Assertion

This case can overwrite the default assert function to custom assert function.

```go

type Person struct {
    name string
    createdAt time.Duration
}

test.New(
    test.WithCase(
        caser.New(
            WithName("case_1"),
            WithArg("vald"),
            WithWant(&Person {
                name: "vald",
            }, nil),
            WithAssertFunc(func(gots, wants []interface{}) error {
                if len(gots) != len(wants) {
                    return errors.Errorf("length not equals. want: %d, but got: %v", len(wants), len(gots))
                }

                //  Because `createdAt` is difficult to compare, so check other fields.
                if got, want := got[0].(*Person).name, want[0].(*Person).name; got != want {
                    return errors.Errorf("person name is not equals. want: %s, but got: %s", want, got)
                }

                if got, want := got[1], want[1]; !reflect.DeepEquals(got, want) {
                    return errors.Errorf("error is not equals. want: %v, but got: %v", want, got)
                }

                return nil
            }),
        ),
    ),
    test.WithDriverFunc(
        func(ctx context.Context, d test.DataProvider) []interface{} {
            p, err := Register(d.Args[0].(string))
            
            
            return []interface {
                p, err
            }
        }
    ),
).Run(context.Background(), t)
```

2. Field Assertion

This case can validate the arguments of the mock object.

```go
test.New(
    test.WithCase(
        caser.New(
            WithName("case_1"),
            WithArg(100, "data"),
            WithFieldFunc(fund(t *testing.T) []interface{} {
                t.Helper()

                mockR := &MockRepository {
                    RegisterFunc: func(id int, data string) {

                        // Validate argument.
                        // If id is invalid, An error is output
                        if id < 0 {
                            t.Errorf("invalid id: %d", id)
                        }

                        if len(data) == "" {
                            t.Error("data is empty")
                        }
                    }
                }

                return []interface{} {
                    mockR,
                }
			}),
            WithWant(nil),
        ),
    ),
    test.WithDriverFunc(
        func(ctx context.Context, d test.DataProvider) []interface{} {
            c := &client {
                repo: d.Field()[0].(*Repository),
            }

            err := c.Redister(d.Args[0].(int), d.Args[1].(string))

            return []interface{} {
                err,
            }
        }
    ),
).Run(context.Background(), t)
```
