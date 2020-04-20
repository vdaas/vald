# Go Style Guide in Vald

## Table of Contents

- [Introduction](#introduction)
- [Code Formatting and Naming Convension](#code-formatting-and-naming-convension)
  - [Project Layout](#project-layout)
  - [Packages](#packages)
  - [Interfaces](#interfaces)
  - [Structs](#structs)
    - [Struct initialization](#struct-initialization)
  - [Variables and Constant](#variables-and-constant)
  - [Methods](#methods)
    - [Getter and Setter](#getter-and-setter)
  - [Error handling](#error-handling)
  - [Logging](#logging)
- [Program comments](#program-comments)
- [Documentation](#documentation)
- [Internal packages](#internal-packages)
- [Dependency management and Build](#dependency-management-and-build)
- [Test](#test)
  - [Table-Driven-Test](#table-driven-test)
  - [The steps to create a Table-Driven-Test](#the-steps-to-create-a-table-driven-test)

## Introduction

This guideline includes the coding style for all Vald contributors and reviewers. Everyone should follow this guideline to keep the style consistent so everyone can understand and contribute to Vald easier once they learn this guideline. You should have the basic knowledge of how to write Golang before contributing to Vald. If you found any bug please create [a GitHub issue](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=) and we will work on it.

Please also read the [Contribution guideline](../../CONTRIBUTING.md) before you start contributing to Vald.

## Code Formatting and Naming Convension

Code formatting and naming conventions affect coding readability and maintainability. Every developer has a different coding style, luckily Golang provides tools to format source code and checking for the potential issue in the source code. We recommend using [goimports](https://github.com/golang/tools/tree/master/cmd/goimports) to format the source code in Vald, and [golangci-lint](https://github.com/golangci/golangci-lint) with `--enable-all` option. We suggest everyone install the plugin for your editor to format the code once you edit the code automatically, and  we suggest using `make update/goimports` command if you want to format the source code manually.

But having tools to format source code doesn't mean you do not need to care the formatting of the code, for example:

```go
// bad
badStr := "apiVersion: v1\n" +
   "kind: Service\n" +
   "metadata:\n" +
   "  name: grafana\n" +
   "spec:\n" +
   "  ports:\n" +
   "    - name: grafana\n" +
   "      port: 3000\n" +
   "      targetPort: 3000\n" +
   "      protocol: TCP\n"

// good
goodStr := `apiVersion: v1
kind: Service
metadata:
  name: grafana
spec:
  ports:
    - name: grafana
      port: 3000
      targetPort: 3000
      protocol: TCP
`
```

### Project Layout

The project layout includes the folder and the file structure in the project. We follow the [project-layout](https://github.com/golang-standards/project-layout) example in Vald.

### Packages

The package defines the context of the objects in the package, for example, the corresponding methods and structs belong to the corresponding package. Unlike other languages like Java, in Golang we use the package name to declare which context of the object we are going to use. For example in [time](https://golang.org/pkg/time/) package, it defines all the objects about time like `time.Now()` method to get the current time.

Here is the naming conventions of the package:

- All lower case.

```go
// bad
package ioUtil

// good
package ioutil
```

- No plurals.

```go
// bad
package times

// good
package time

```

- Should be the same as the folder name.
- Should keep as simple as it should, and should contain only one specific context in the package.

```go
// bad
package encodebase64

// good
package base64 // inside the encoding/base64 folder
```

- Should not be too general, for example `util` or `helper`, which will cause all the objects from different contexts to be store in one package. If you want to name the package as `util`,  please define the more specific package name more  `ioutil` or `httputil`.

All packages should contain `doc.go` file under the package to describe what is the package is. For example, under the folder name called `cache`, should contains a file named `doc.go`, which contains the package documentation. Here is the example `doc.go` of the cache package.

```go
// Package cache provides implementation of cache
package cache
````

### Interfaces

Interface defines the program interface for usability and future extendability.
Unlike other languages like Java, Golang supports implicit interface implementation. The type implements do not need to specify the interface name; to "implements" the interface the structs only need to defined the methods the same as the interface, so please be careful to define the method name inside the interface.

The interface should be named as:

- Use MixedCaps

```go
// bad
type Roundtripper interface {
  // interface definition
}

// good
type RoundTripper interface {
  // interface definition
}
```

- Use understandable common short form.

```go
// bad
type ATSigner interface {
  // interface definition
}

// good
type AccessTokenSigner interface {
  // interface definition
}

// good
type HTTPServer interface {
  // interface definition
}
```

### Structs

Structs in Golang is the object definition, we can attach any fields and methods to the struct. The naming convention is the same as the interface one.
If the structs are implementing the interface, the structs name should be related to the interface, for example:

```go
type Listener interface {
   // interface definition
}

// Listener instance for file
type FileListener struct {
   // implement listener interface
}

// Listener instance for HTTP
type HTTPListener struct {
   // implement listener interface
}
```

#### Struct initialization

There are many ways to initialize structs in Golang, base on the use case we can decide which way to initialize structs in Golang.
To initialize struct, it is suggested to use `new(T)` instead of `&T{}` unless you need to initialize with values. For example:

```go
type Something struct {
    name string
}

// good
a := new(Something)

// bad
b := new(Something)
b.name = "Mary"

// use this syntax instead of b
c := &Something{
    name: "Mary",
}
```

To initialize complex structs, we can use [functional option pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis). Please read [server.go](../../../internal/servers/servers.go) and [option.go](../../../internal/servers/option.go) for the reference implementation.
The options implementation should be separated as another file called `option.go` to improve the readability of the source code, and the method name should start with `With` word to differentiate with other methods.

### Variables and Constant

The variable and the constant should be named as:

- Use MixedCaps

```go
// bad
yamlprocessor := new(something)

// good
yamlProcessor := new(something)
```

- Use short form.

```go
// bad
yamlString := "something"

// good
yamlStr := "something"

// in some case it is acceptable and actually if it is easier to read
s := new(something)
```

The variable and the constant name may lead to misunderstanding or confusion, so if the variable and constant name are different to understand, please write some comment even if it is a private member.

```go
// somebody may not understand this variable, so write a simple comment to the variable definition
sac := new(something) // signed access token
```

If the multiple variables and the constants have the same grouping, please use the grouping name as the prefix of the variable and constant name.

```go
// Same group of variable (error), so add a prefix `Err` to each error variables
ErrInvalidCacherType = New("invalid cacher type")
// ErrXXXXXXX
```

### Methods

The method name should be named as:

- Use MixedCaps.

```go
// bad
func (s *something) somemethod() {}

// bad
func (s *something) some_method() {}

// good
func (s *something) someMethod() {}
```

- Do not use short form unless the function name is too long.

```go
// bad
func (s *something) genereateRandomNumber() int {}

// good
func (s *something) genRandNum() int {}
```

- It should be understandable for everyone even if it is a private method.

#### Getter and Setter

The Getter and Setter are almost the same as other languages, but the naming convention of the Getter method is different from other languages. Instead of `GetVar1()`, the getter of `Var1` should be the same as the variable name itself `Var1()`.

```go
// getter of the `signedTok`
func (s *something) SignedTok() string{
    return s.signedTok
}

// setter of the `signedTok`
func (s *something) SetSignedTok(st string) {
    s.signedTok = st
}
```

### Error handling

All errors should define in [internal/errors package](../../internal/errors). All errors should be start with `Err` prefix, and all errors should be handle if possible.

Please use [internal/errgroup](../../internal/errgroup) for synchronized error handling on multi-goroutine processing.

### Logging

We define our own logging interface in [internal/log package](../../internal/log). By default we use [glg](https://github.com/kpango/glg) to do the logging internally.
We defined the following logging levels.

| Log level | Description                                                                                                                                                                                                                                    | Example situation                                                                                                                                  | Example message                                                                                                                                                                                                        |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| DEBUG     | Fine-grained information for debugging and diagnostic the application.<br>Enable this logging level in production environment may cause performance issue.                                                                                     | An entry that will insert into the database with details.<br>An HTTP request and response that will send and receive from the server with details. | User 1 will insert into the database, name: Mary, Age: 19.<br>An HTTP request will send to http://example.com, with the body:<br>key1:value1, key2:value2.<br>The response of the HTTP request: body: HTTPResponseBody |
| INFO      | Normal application behavior, to trace what is happening inside the application.                                                                                                                                                                | Inserted entry into the database without details. <br>HTTP requests sent to the server without details. <br>The server is started or stopped.      | User 1 is inserted into the database.<br>The HTTP request is sent to XXX server successfully.                                                                                                                          |
| WARN      | The message that indicates the application may having the issue or occurring unusual situation,<br>but does not affect the application behavior.<br>Someone should investigate the warning later.                                              | Failed to insert entry into the the database, but success with the retry.<br>Failed to update the cache, and the cache is not important.           | User 1 is successfully inserted into the database with retry,<br>retry count: 1, error: ErrMsg1, retry count: 2, error: ErrMsg2                                                                                        |
| ERROR     | The message that indicates the application is having a serious issue or,<br>represent the failure of some important going on in the application.<br>It does not cause the application to go down.<br>Someone must investigate the error later. | Failed to insert an entry into the database, with retry count exceeded.<br>Failed to update the cache, and the cache is not important.             | User 1 is failed to insert in the database, errors:<br>retry count: 1, error: ErrMsg1, retry count: 2, error: ErrMsg2, ....                                                                                            |
| FATAL     | Message that indicate the application is corrupting or having serious issue.<br>The application will go down after logging the fatal error. <br>Someone must investigate and resolve the fatal as soon as possible.                            | Failed to init the required cache during the application start.                                                                                    |                                                                                                                                                                                                                        |

## Program comments

Program comments make easier to understand the source code. We suggest not to write many comments inside the source code unless the source code is very complicated and confusing; otherwise we should divide the source code into methods to keep the readability and usability of the source code.

Everyone should write the comments to all the public objects on your source code, like public packages, interface, structs, methods, and even public constant and variable. The godoc will be generated base on the comment of source code.

## Documentation

Documentation is generated from the program comments. Please refer to [Godoc](https://godoc.org/github.com/vdaas/vald) for the program documentation.

## Internal packages

Vald implements its internal package to extend and customize the functionality of the standard library and third-party library.
We should use the internal package instead of standard libray to implement Vald.
Please refer to [godoc](https://godoc.org/github.com/vdaas/vald/internal) for the internal package document.

## Dependency management and Build

We should use `go mod tidy` to manage the `go.mod` file in the project.

## Test

The testing guideline has 3 important rules for the coding quality and readability:

1. Use Table-Driven-Test
1. Keep code coverage over 85%
   - test coverage != high testing quality, but low coverage means bad testing quality
   - check with the following commands `go test -cover ./...`
1. Test all use cases and error cases

### Table-Driven-Test

Use table-driven tests with subtests to avoid duplicating code.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
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
        str: "192.0.2.0:8000",
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
</tr>
</tbody>
</table>

Table-Driven-Test makes it easy to add a new test case.

We define the test case table as `map[string]func(*testing.T)test`, which is referred to the test case name and the test case implementation `tt`.

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

### The steps to create a Table-Driven-Test

1. `args` structure

    If there are two or more arguments to be passed to the method, create a `args` structure. If there is only one argument, do not create an `args` structure.

    ```go
    type args struct {
        host string
        port string
    }
    ```

1. `field` structure

    If you create an object and test its methods, create a `field` struct if the object has two or more fields to initialize. If there is only one field, do not create `field` structure.

    ```go
    type field struct {
        host string
        port string
    }
    ```

1. `test` structure

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
