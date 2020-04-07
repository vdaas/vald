# Go Style Guide in Vald

## Table of Contents

- [Go Style Guide in Vald](#go-style-guide-in-vald)
  - [Table of Contents](#table-of-contents)
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
  - [Style](#style)
  - [Test](#test)
    - [Table-Driven-Test](#table-driven-test)
    - [The steps to create a Table-Driven-Test.](#the-steps-to-create-a-table-driven-test)

## Introduction
This guideline includes the coding style for all Vald contributors and reviewers. Everyone should follow this guideline to keep the style consistent so everyone can understand and contribute to Vald easier once they learn this guideline. You should have the basic knowledge of how to write Golang before contributing to Vald. If you found any bug please create [a GitHub issue](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=) and we will work on it.

Please also read the [Contribution guideline](../../CONTRIBUTING.md) before you start contributing to Vald.

## Code Formatting and Naming Convension
Code formatting and naming conventions affect coding readability and maintainability. Every developer has a different coding style, luckily Golang provides tools to format source code and checking for the potential issue in the source code. We suggest using [gofmt](https://golang.org/cmd/gofmt/) to format the source code in Vald, and [golint](https://github.com/golang/lint). We suggest everyone install the plugin for your editor to automatically format the code once you edit the code.
But having tools to format source code doesn't mean you do not need to care the formatting of the code, for example:
```golang
yamlStr := "apiVersion: v1\n" +
   "kind: Service\n" +
   "metadata:\n" +
   "  name: grafana\n" +
   "spec:\n" +
   "  ports:\n" +
   "    - name: grafana\n" +
   "      port: 3000\n" +
   "      targetPort: 3000\n" +
   "      protocol: TCP\n"
```

It is very hard to read and maintain, in this case, we should follow this style.
```golang
yamlStr1 := `apiVersion: v1
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

The project layout includes the folder and the file structure in the project. We follow the [prject-layout](https://github.com/golang-standards/project-layout) example in Vald.

### Packages

The package defines the context of the objects in the package, for example, the corresponding methods and structs belong to the corresponding package. Unlike other languages like Java, in Golang we use the package name to declare which context of the object we are going to use. For example in [time](https://golang.org/pkg/time/) package, it defines all the objects about time like `time.Now()` method to get the current time.

Here is the naming conventions of the package:

- All lower case.
- No plurals.
- Should be the same as the folder name.
- Should keep as simple as it should, and should contain only one specific context in the package.
- Should not be too general, for example `util` or `helper`, which will cause all the objects from different contexts to be store in one package. If you want to name the package as `util`,  please define the more specific package name more  `ioutil` or `httputil`.

All packages should contain `doc.go` file under the package to describe what is the package is. For example, under the folder name called `cache`, should contains a file named `doc.go`, which contains the package documentation. Here is the example `doc.go` of the cache package.

```golang
// Package cache provides implementation of cache
package cache
````

### Interfaces

Interface defines the program interface for usability and future extendability.
Unlike other languages like Java, Golang supports implicit interface implementation. The type implements do not need to specify the interface name; to "implements" the interface the structs only need to defined the methods the same as the interface, so please be careful to define the method name inside the interface.

The interface should be named as:

- Use MixedCaps
- Do not use short form unless it is a common terms.

```golang
type RoundTripper interface {
    // interface definition
}
```

### Structs

Structs in Golang is the object definition, we can attach any fields and methods to the struct. The naming convention is the same as the interface one.
If the structs are implementing the interface, the structs name should be related to the interface, for example:

```golang
type Listener interface {
   // Interface definition
}

// Listener instance for file
type FileListener struct {

}

// Listener instance for HTTP
type HttpListener struct {

}
```

#### Struct initialization

There are many ways to initialize structs in Golang, base on the use case we can decide which way to initialize objects in Golang.
In Vald we use [functional option pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) to initialize complex structs. Please read [server.go](../../../internal/servers/servers.go) and [option.go](../../../internal/servers/option.go) for the reference implementation.
The functional options should be separated as another file to improve the readability of the source code, and the method name should start with `With` word to differentiate with other methods.

Also, you can use `&T{}` to initialize the struct. Do not use the `new(T)` method to initialize the struct.

### Variables and Constant

The variable and the constant should be named as:

- Use MixedCaps
- Global variables and constants should not use short form unless it is a common terms.
- Private variables and contants should use the short form to improve readability.

The variable and the constant name may lead to misunderstanding or confusion, so if the variable and constant name are different to understand, please write some comment even if it is a private member.
If the multiple variables and the constants have the same grouping, please use the grouping name as the prefix of the variable and constant name.

Here is some example of the declaration of variables and constants:

```
/* Global variables */
// Same group of variable (error), so add a prefix `Err` to each error variables
ErrInvalidCacherType = New("invalid cacher type")
// ErrXXXXXXX

/* Private variables */
// This variable needs comment in order to understand
// sds represent the shut down strategy
sds     []string

// This variable name is common so no comment requires for this variable
eg      errgroup.Group
```

### Methods

The method name should be named as:

- Use MixedCaps.
- Should not use short form unless it is a common term.
- It should be understandable for everyone even if it is a private method.

#### Getter and Setter

The Getter and Setter are almost the same as other languages, but the naming convention of the Getter method is different from other languages. Instead of `GetVar1()`, the getter of `Var1` should be the same as the variable name itself `Var1()`.

### Error handling

### Logging

We define our own logging interface in [internal/log package](../../internal/log). By default we use [glg](https://github.com/kpango/glg) to do the logging internally.
We defined the following logging levels.

| Log level | Description                                                                                                                                                                                                               | Example                                                                                                                               |
|-----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| DEBUG     | Fine-grained information for debugging and diagnostic the application. Enable this logging level in production environment may cause performance issue.                                                                   | Entry that will insert into the database with details. HTTP request that will send to the server with details.                        |
| INFO      | Normal application behavior, to trace what is happening inside the application.                                                                                                                                           | Inserted entry into the database without much details. HTTP request sent to the server without details. Server is started or stopped. |
| WARN      | Message that indicate the application may have issue or occurring unusual situation,  but does not affect the application behavior. Someone should investigate the warning later.                                         | Failed to insert entry into the the database, but success with the retry. Failed to update the cache, and the cache is not important. |
| ERROR     | Message that indicate the application is having serious issue or  represent the failure of some important going on in the application. It do not cause the application go down. Someone must investigate the error later. | Failed to insert entry into the the database even with retry. Failed to update the cache, and the cache is very important.            |
| FATAL     | Message that indicate the application is corrupting or having serious issue. The application will go down after logging the fatal error. Someone must investigate and resolve the fatal as soon as possible.              | Failed to init the required cache during the application start.                                                                       |

## Program comments

Program comments make easier to understand the source code. We suggest not to write many comments inside the source code unless the source code is very complicated and confusing; otherwise we should divide the source code into methods to keep the readability and usability of the source code.

Everyone should write the comments to all the public objects on your source code, like public packages, interface, structs, methods, and even public constant and variable. The godoc will be generated base on the comment of source code.

## Documentation

Documentation is generated based on the program comments. Please refer to [Godoc](https://godoc.org/github.com/vdaas/vald) for the program documentation.

## Internal packages

Vald implements its internal package to extend the functionality of the standard library and third-party library. Please refer to [godoc](https://godoc.org/github.com/vdaas/vald/internal) for the internal package document.

## Dependency management and Build

We should use `go mod tidy` to manage the `go.mod` file in the project.

## Style

## Test

The testing guideline has 2 important rules for the coding quality and readability

1. Use Table-Driven-Test
2. Keep code coverage over 85%
   - test coverage != high testing quality, but low coverage means bad testing quality
   - check with the following commands `go test -cover ./...`
3. Test all use cases and error cases

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

### The steps to create a Table-Driven-Test.

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
