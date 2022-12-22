# Go Style Guide in Vald

## Introduction

This guideline includes the coding style for all Vald contributors and reviewers.
Everyone should follow this guideline to keep the style consistent, so everyone can understand and contribute to Vald easier once they learn it.
You should have the basic knowledge of how to write Go before contributing to Vald.
If you find any bugs, please create [a GitHub issue](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=), and we will work on it.

Please also read the [Contribution guideline](../contributing/contributing-guide.md) before you start contributing to Vald.

## Code Formatting and Naming Convention

Code formatting and naming conventions affect coding readability and maintainability.
Every developer has a different coding style.
Luckily Go provides tools to format source code and check for potential issues in the source code.
We recommend using [golines](https://github.com/segmentio/golines), [gofumpt](https://github.com/mvdan/gofumpt), and [goimports](https://github.com/golang/tools/tree/master/cmd/goimports) to format the source code in Vald and [golangci-lint](https://github.com/golangci/golangci-lint) with the `--enable-all` option.
We suggest everyone install the plugin for your editor to format the code once you edit the code automatically and use `make format/go` command if you want to format the source code manually.

But having tools to format source code does not mean you do not need to care about the formatting of the code, for example:

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
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
```

</td><td>

```go
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

</td></tr></tbody></table>

### Project Layout

The project layout includes the folder and the file structure in the project. We follow the [project-layout](https://github.com/golang-standards/project-layout) example in Vald.

### Packages

The package defines the context of the objects in the package; for example, the corresponding methods and structs belong to the corresponding package.
Unlike other languages like Java, in Go, we use the package name to declare which context of the object we are going to use.
For example, the [time package](https://golang.org/pkg/time/) defines all the objects about time, like `time.Now()` method to get the current time.

Here are the naming conventions of the package:

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
- Should keep as simple as it should and contain only one specific context in the package.

  ```go
  // bad
  package encodebase64

  // good
  package base64 // inside the encoding/base64 folder
  ```

- Should not be too general, for example, `util` or `helper`, which will cause all the objects from different contexts to be stored in one package.
  If you want to name the package as `util`, please define the more specific package name more `ioutil` or `httputil`.

All packages should contain a `doc.go` file under the package to describe what is the package is. For example, under the folder name called `cache`, should contain a file named `doc.go`, which contains the package documentation. Here is the example `doc.go` of the cache package.

    ```go
    // Package cache provides an implementation of cache
    package cache
    ````

### General style

This section describes the general guideline for the Vald programming style, every Vald contributor should keep these general guidelines in mind while working on the implementation of Vald.

#### Order of declaration

Put the higher priority or frequently used declaration on top of other declarations.
It makes it easier for Vald to read and search its target source code in Vald.

For example, the interface declaration should have higher priority than the struct or function declaration; hence it should be put above other declarations.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
type S struct {}

func (s *S) fn() {}

type I interface {}
```

</td><td>

```go
type I interface {}

type S struct {}

func (s *S) fn() {}
```

</td></tr></tbody></table>

#### Group similar definition

Group similar definitions such as a struct or an interface declaration.
We should not group an interface and a struct declaration in the same block, for example:

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
type (
    I interface {}
    I2 interface {}

    s struct {}
    s2 struct {}
)
```

</td><td>

```go
type (
    I interface {}
    I2 interface {}
)

type (
    s struct {}
    s2 struct {}
)
```

</td></tr></tbody></table>

### Interfaces

The interface defines the program interface for usability and future extendability.
Unlike other languages like Java, Go supports implicit interface implementation.
The type implements do not need to specify the interface name; to "implement" the interface, the structs only need to define the methods the same as the interface, so please be careful to define the method name inside the interface.

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

Structs in Go is the object definition; we can attach any fields and methods to the struct. The naming convention is the same as the interface one.
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

There are many ways to initialize structs in Go; based on the use case, we can decide which way to initialize structs in Go.
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

To initialize complex structs, we can use a [functional option pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis). Please read [server.go](https://github.com/vdaas/vald/blob/main/internal/servers/servers.go) and [option.go](https://github.com/vdaas/vald/blob/main/internal/servers/option.go) for the reference implementation.
The options implementation should be separated as another file called `option.go` to improve the readability of the source code, and the method name should start with `With` word to differentiate it from other methods.

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

The variable and the constant name may lead to misunderstanding or confusion, so if the variable and constant name are different to understand, please write some comments, even if it is a private member.

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

In this section, rules also apply to the `function` (without receiver). The method name should be named as:

- Use MixedCaps.

  ```go
  // bad
  func (s *something) somemethod() {}

  // bad
  func (s *something) some_method() {}

  // good
  func (s *something) someMethod() {}
  ```

- Avoid using long function name.
- Do not use short form unless the function name is too long.

  ```go
  // bad
  func (s *something) genereateRandomNumber() int {}

  // good
  func (s *something) genRandNum() int {}
  ```

- It should be understandable for everyone even if it is a private method.

#### Getter and Setter

The Getter and Setter are almost the same as other languages, but the naming convention of the Getter method is different from other languages. Instead of `GetVar1()`, the getter of `Var1` should be the same as the variable named `Var1()`.

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

### Unused Variables

An unused variable may increase the complexity of the source code.
It may confuse the developer hence introducing a new bug.
So please delete the unused variable.

Generally, the unused variable should be reported during compilation, but in some cases, the compiler may not report an error.
This is an example of an unused variable declaration that does not cause a compilation error.

```go
// In this case, this example are not using `port` field, but dose not cause a compilation error.
// So please delete `port` field of `server`.

type server struct {
    addr string
    port int  // you have to delete this field.
}

// The `port` field of `server` is not used.
srv := &server {
    addr: "192.168.33.10:1234",
}

if err := srv.Run(); err != nil {
    log.Fatal(err)
}
```

### Error handling

All errors should define in [internal/errors package](https://github.com/vdaas/vald/blob/main/internal/errors). All errors should be start with `Err` prefix, and all errors should be handle if possible.

Please use [internal/errgroup](https://github.com/vdaas/vald/blob/main/internal/errgroup) for synchronized error handling on multi-goroutine processing.

### Error checking

All functions return an `error` if the function can fail.
It is very important to ensure that error checking is performed.
To reduce human mistakes that miss the error checking, please check the error using the following style:

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
err := fn()
if err != nil {
    // handle error
}
```

</td><td>

```go
if err := fn(); err != nil {
    // handle error
}
```

</td></tr></tbody></table>

If you need the value outside the if statement, please use the following style:

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
if conn, err := net.Dial("tcp", "localhost:80");  err != nil {
    // handle error
} else {
    // use the conn
}
```

</td><td>

```go
conn, err := net.Dial("tcp", "localhost:80")
if err != nil {
    // handle error
}
// use the conn
```

</td></tr></tbody></table>

### Logging

We define our own logging interface in [internal/log package](https://github.com/vdaas/vald/blob/main/internal/log).
By default, we use [glg](https://github.com/kpango/glg) to do the logging internally.
We defined the following logging levels.

| Log level | Description                                                                                                                                                                                                                                | Example situation                                                                                                                                 | Example message                                                                                                                                                                                                        |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| DEBUG     | Fine-grained information for debugging and diagnosing the application.<br>Enabling this logging level in the production environment may cause performance issues.                                                                          | An entry that will insert into the database with details.<br>An HTTP request and response will send and receive from the server with details.     | User 1 will insert into the database, name: Mary, Age: 19.<br>An HTTP request will send to http://example.com, with the body:<br>key1:value1, key2:value2.<br>The response of the HTTP request: body: HTTPResponseBody |
| INFO      | Normal application behavior to trace what is happening inside the application.                                                                                                                                                             | Inserted entry into the database without details. <br>HTTP requests are sent to the server without details. <br>The server is started or stopped. | User 1 is inserted into the database.<br>The HTTP request is sent to XXX server successfully.                                                                                                                          |
| WARN      | The message that indicates the application may have an issue or occurring unusual situation,<br>but does not affect the application behavior.<br>Someone should investigate the warning later.                                             | Failed to insert an entry into the database, but successful with the retry.<br>Failed to update the cache, and the cache is not important.        | User 1 is successfully inserted into the database with retry,<br>retry count: 1, error: ErrMsg1, retry count: 2, error: ErrMsg2                                                                                        |
| ERROR     | The message indicates the application has a serious issue or,<br>represents the failure of some important thing going on in the application.<br>It does not cause the application to go down.<br>Someone must investigate the error later. | Failed to insert an entry into the database, with retry count exceeded.<br>Failed to update the cache, and the cache is not important.            | User 1 is failed to insert in the database with errors:<br>retry count: 1, error: ErrMsg1, retry count: 2, error: ErrMsg2, ....                                                                                        |
| FATAL     | The message indicates the application is corrupting or has a serious issue.<br>The application will go down after logging the fatal error. <br>Someone must investigate and resolve the fatality as soon as possible.                      | Failed to init the required cache during the application start.                                                                                   |                                                                                                                                                                                                                        |

## Implementation

This section includes some examples of general implementation widely used in Vald.
The implementation may differ based on your use case.

### Functional Option

The functional option pattern is widely used in Vald.
You can refer to [this section](#Struct-initialization) for more details of the use case of this pattern.

We provide the following errors to describe the error to apply the option.

| Error                    | Description                                                                               |
| ------------------------ | ----------------------------------------------------------------------------------------- |
| errors.ErrInvalidOption  | Error to apply the option, and the error is ignorable and the warning log will be printed |
| errors.ErrCriticalOption | Critical error to apply the option, the error cannot be ignored and should be handled     |
| errors.ErrIgnoredOption  | The option is ignored, and the debug log will be printed                                  |

We strongly recommend the following implementation to set the value using the functional option.

If an invalid value is set to the functional option, the `ErrInvalidOption` error defined in the [internal/errors/option.go](https://github.com/vdaas/vald/blob/main/internal/errors/option.go) should be returned.

The name argument (the first argument) of the `ErrInvalidOption` error should be the same as the functional option name without the `With` prefix.

For example, the functional option name `WithVersion` should return the error with the argument `name` as `version`.

```go
func WithVersion(version string) Option {
    return func(c *client) error {
        if len(version) == 0 {
            return errors.NewErrInvalidOption("version", version)
        }
        c.version = version
        return nil
    }
}
```

We recommend the following implementation to parse the time string and set the time to the target struct.

```go
func WithTimeout(dur string) Option {
    func(c *client) error {
        if dur == "" {
            return errors.NewErrInvalidOption("timeout", dur)
        }
        d, err := timeutil.Parse(dur)
        if err != nil {
            return errors.NewErrInvalidOption("timeout", dur, err)
        }
        c.timeout = d
        return nil
    }
}
```

We recommend the following implementation to append the value to the slice if the value is not nil.

```go
func WithHosts(hosts ...string) Option {
    return func(c *client) error {
        if len(hosts) == 0 {
            return errors.NewErrInvalidOption("hosts", hosts)
        }
        if c.hosts == nil {
            c.hosts = hosts
        } else {
            c.hosts = append(c.hosts, hosts...)
        }
        return nil
    }
}
```

If the functional option error is critical, we should return `ErrCriticalOption` error instead of `ErrInvalidOption`.

```go
func WithConnectTimeout(dur string) Option {
    return func(c *client) error {
        if dur == "" {
            return errors.NewErrInvalidOption("connectTimeout", dur)
        }
        d, err := timeutil.Parse(dur)
        if err != nil {
            return errors.NewErrCriticalOption("connectTimeout", dur, err)
        }

        c.connectTimeout = d
        return nil
    }
}
```

We need to handle the error returned from the functional option on the caller side.

If the option fails to apply, an error wrapped with `ErrOptionFailed` defined in the [internal/errors/errors.go](https://github.com/vdaas/vald/blob/main/internal/errors/errors.go) should be returned.

We recommend the following implementation to apply the options.

```go
func New(opts ...Option) (Server, error) {
    srv := new(server)
    for _, opt := range opts {
        if err := opt(srv); err != nil {
            werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

            e := new(errors.ErrCriticalOption)
            if errors.As(err, &e) {
                log.Error(werr)
                return nil, werr
            }
            log.Warn(werr)
        }
    }
}
```

### Constructor

In Vald, the functional option pattern is widely used when we create an object.

When setting the value with the functional option, the value is validated inside the option method.

However, we may forget to set the required fields when creating the object; hence the target object will remain nil.
Therefore, we strongly suggest validating the object during initialization.

If we forget to set the options method, an error will be returned so we can handle it properly.

```go
func func New(opts ...Option) (Server, error) {
    srv := new(server)
    for _, opt := range append(defaultOptions, opts...) {
        if err := opt(srv); err != nil {
            werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

            e := new(errors.ErrCriticalOption)
            if errors.As(err, &e) {
                log.Error(werr)
                return nil, werr
            }
            ue := new(errors.ErrIgnoredOption)
            if errors.As(err, &ue) {
                log.Debug(werr)
            } else {
                log.Warn(werr)
            }
        }
    }

    if srv.eg == nil {
        return nil, errors.NewErrInvalidOption("eg", srv.eg)
    }

    return srv, nil
}
```

We also recommend that you use the default options and the unexported functional option to set the objects so that we cannot use them externally.

```go
var defaultOptions = []Option {
    func(s *server) error {
        s.ctxio = io.New()
        return nil
    },
}
```

## Program comments

Program comments make it easier to understand the source code.
We suggest not writing many comments inside the source code unless the source code is very complicated and confusing; otherwise, we should divide the source code into methods to keep readability and usability of the source code.

Everyone should comment on all the public objects on your source code, like public packages, interfaces, structs, methods, and even public constants and variables.
The `Godoc` will be generated based on the comment of the source code.

## Documentation

Documentation is generated from the program comments. Please refer to [Godoc](https://pkg.go.dev/github.com/vdaas/vald) for the program documentation.

## Internal packages

Vald implements its internal package to extend and customize the standard library and third-party library functionality.
We should use the internal package instead of the standard library to implement Vald.
Please refer to [Godoc](https://pkg.go.dev/github.com/vdaas/vald/internal) for the internal package document.

## Dependency management and Build

We should use `go mod tidy` to manage the `go.mod` file in the project.

## Test

The testing guideline has three essential rules for coding quality and readability:

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

We define the test case table as `map[string]func(*testing.T)test`, which is referred to as the test case name and the test case implementation `tt`.

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

### Create a Table-Driven-Test

#### Structures

1. `args` structure

   If two or more arguments are passed to the method, create an `args` structure. If there is only one argument, do not create an `args` structure.

   ```go
   type args struct {
       host string
       port string
   }
   ```

1. `field` structure

   If you create an object and test its methods, create a `field` struct if the object has two or more fields to initialize. If there is only one field, do not create a `field` structure.

   ```go
   type field struct {
       host string
       port string
   }
   ```

1. `test` structure

   `test` structure has `args` and `field` structure and `checkFunc` function.
   Create a field or an args structure if you need it.
   The `checkFunc` function is used to check the return value of the function being tested.

   ```go
   type test struct {
       args args
       field field
       checkFunc func(t *testing.T, err error)
   }
   ```

#### Inputs

1. test case name

   Test case names should be readable, meaningful, and understandable easily.
   If you create a new test, the test name should be named based on the below-naming templates

   - Success cases:
     - Start with `success` or `{verb} success` or `success {verb}`
     - End with the condition `when {condition}` or `with {condition}
     - e.g.: `set success when the value is default value`
   - Fail/Error cases:
     - Start with `fail` or `{verb} fail` or `fail {verb}`
     - End with the condition `when {condition}` or `with {condition}
     - e.g.: `fail option setting when value is invalid value(string)`
   - Return cases:
     - If the test case do not match with `Success` or `Fails`, please use `Return` pattern.
     - Start with `Returns {Object}`
     - End with the condition `when {condition}` or `with {condition}
     - e.g.: `return invalid error when the input is invalid`

1. testing arguments

   Input arguments for testing should be a meaningful value.
   We should test with more realistic values to produce a more realistic testing result.

   For example, to test the function with `host` argument, you should set your hostname (e.g., `vald.vdaas.com`) as the input value to the `host` argument.

Example:

```go
type args struct {
    host string
    port string
}

type field struct {
    timeout time.Duration
}

type want struct {
    err error
}

type test struct {
    name       string
    args       args
    field      field
    want       want
    checkFunc func(err error) error
}

defaultCheckFunc := func(w want, err error) error {
    if err != w.err {
        return errors.Errorf("got error: %v, want: %v", err, w.err)
    }
}

tests := []test {
    {
        name: "success send when host and port are correct value",
        args: args {
            host: "vdaas.vald.org",
            port: "80",
        },
        field: field {
            host: "vdaas.vald.org",
            port: "80",
        },
        want: want {
            err: nil,
        },
    },
}

for _, tc := range tests {
    test := tc
    t.Run(test.name, func(tt *tesint.T) {
        checkFunc = defaultCheckFunc
        if test.checkFunc != nil {
            checkFunc = test.checkFunc
        }

        c := client {
            timeout: test.field.timeout,
        }
        err := c.Send(test.args.addr, test.args.txt)
        if err := checkFunc(err); err != nil {
            tt.Errorf("error = %v", err)
        }
    })
}

```

### Generate test code

We implement our own [gotests](https://github.com/cweill/gotests) template to generate test code.
If you want to install `gotest` tools, please execute the following command under the project root directory.

```bash
make gotests/install
```

If you use the following command to generate the missing test code.

```bash
make gotests/gen
```

After executing the command above, the file `*target*_test.go` will be generated for each Go source file.
The test code generated follows the table-driven test format.
You can implement your test code under the `tests` variable generated following the table-driven test format.

### Customize test case

We do not suggest modifying the generated code other than the `tests` variable
Still, in some cases, you may need to change the generated code to meet your requirement, for example:

1. init() function

   init() function is executed automatically before the test is started.
   You may need to initialize some singleton before running your test cases.
   For example, Vald uses [glg](https://github.com/kpango/glg) library for logging by default; if the logger is not initialized before the test, the nil pointer error may be thrown during the test is running.
   You may need to implement `init()` function like:

   ```go
   func init() {
       log.Init()
   }
   ```

   And place it on the header of the test file.

1. goleak option

   The generated test code will default use the [goleak](https://github.com/uber-go/goleak) library to test if there is any goroutine leak.
   Sometimes you may want to skip the detection; for example, Vald uses the [fastime](https://github.com/kpango/fastime) library, but the internal goroutine is not closed due to the needs of the library.
   To skip the goleak detection, we need to create the following variable to store the ignore function.

   ```go
   var (
       // goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
       goleakIgnoreOptions = []goleak.Option{
           goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
       }
   )
   ```

   And use it in `VerifyNone()` in the test case.

   ```go
   for _, tc := range tests {
       test := tc
       t.Run(test.name, func(tt *testing.T) {
           // modify the following line
           defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
   ```

   In Vald, we implemented [internal goleak package](https://github.com/vdaas/vald/blob/main/internal/test/goleak/goleak.go) and wrap the goleak validation logic to ignore the common goleak functions.
   In test implementation, we can easily import the internal goleak package and ignore all the necessary goleak functions by default.

   ```go
   // import internal goleak package
   import (
       "github.com/vdaas/vald/internal/test/goleak"
   )

   ...
   for _, tc := range tests {
       test := tc
       t.Run(test.name, func(tt *testing.T) {
           // use internal goleak to ignore necessary function by default
           defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
   ```

1. goleak usage

   There are two methods for valid goroutine leak:

   1. Use `goleak.VerifyNone()` to validate it on each test cases.
   1. Use `goleak.VerifyTestMain()` to validate it on each package.

   By default, the goroutine leak validation is executed in each test case.

   ```go
   for _, tc := range tests {
        test := tc
        t.Run(test.name, func(tt *testing.T) {
            defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
   ```

   In some cases, it may not work as some cleanup or asynchronous processes remaining in the background, and validating it on the test case using `goleak.VerifyNone()` may cause a false alarm.

   To resolve it, we can consider using `goleak.VerifyTestMain()` to validate goleak when all test cases are passed in each package.

   ```go
   // implement TestMain and verify it at package level
   func TestMain(m *testing.M) {
        goleak.VerifyTestMain(m)
   }

   ...

   for _, tc := range tests {
        test := tc
        t.Run(test.name, func(tt *testing.T) {
            // use VerifyTestMain() above to detect goroutine leak
            // defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
   ```

   Please note that the `TestMain()` can only implement once in each package, you may need to find the right place to implement this function on the package.

1. Defer function

   By default, the template provides `beforeFunc()` and `afterFunc()` to initialize and finalize the test case, but in some cases, it may not support your use case.
   For example, `recover()` function only works in `defer()` function.
   If you need to use the `recover()` function to handle the panic in your test code, you may need to implement your custom `defer()` function and change the generated test code.

   For example:

   ```go
   for _, tc := range tests {
        test := tc
        t.Run(test.name, func(tt *testing.T) {
            defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

            // insert your defer function here
            defer func(w want, tt *testing.T) {
                // implement your defer func logic
                if err:= recover(); err != nil {
                    // check the panic
                }
            }(test.want, tt)

            if test.beforeFunc != nil {
                test.beforeFunc(test.args)
            }
            // generated test code
   ```

1. Unused fields

   By default, the template provides the `fields` structure to initialize an object of the test target.
   But in some cases, not all `fields` are needed, so please delete the unnecessary fields.
   For example, the following struct and the corresponding function:

   ```go
   type server struct {
       addr string
       port int
   }
   func (s *server) Addr() string {
       return s.addr
   }
   ```

   And the generated test code is:

   ```go
   func Test_server_Addr(t *testing.T) {
       type fields struct {
           addr string
           port int
       }
       type want struct {
           // generated test code
   ```

   Since the `port` variable is not used in this test case, you can delete the `port` definition in the test case.

   ```go
   func Test_server_Addr(t *testing.T) {
       type fields struct {
           addr string
           // port int   <-- this line should be deleted
       }
       type want struct {
           // generated test code
   ```

1. Context

   If the target function accepts context as an input argument, the test code generated will include it in `args`.

   ```go
    type args struct {
        ctx context.Context
    }
    ...
   ```

   But in test implementation, it is hard to manage the context lifecycle.
   We want to guarantee the context is closed after each test is executed to avoid any missing termination of the process.

   In Vald, we suggest customizing the test implementation from the generated test code to create the context and manage the lifecycle in every test.

   We can remove the context from the `args` struct and create the context in the test execution code.

   ```go
    // we can remove the ctx from args list
    type args struct {
        // ctx context.Context
    }
    ...

    for _, tc := range tests {
        test := tc
        t.Run(test.name, func(tt *testing.T) {
            // create context here
            ctx, cancel := context.WithCancel(context.Background())
            // handle the context lifecycle by canceling the context after the test is executed
            defer cancel()

            ...

            // use the context created above
            got, gotErr := someFunc(ctx)

            ...
   ```

1. Struct initialization

   By default, when testing the function of the struct, the target struct initialization is implemented by setting the data from the `fields` defined in the test case.
   This initialization method has a few disadvantages:

   1. When there are many fields in the struct, it is hard to set them all.
   1. The default value is the zero value of the type, not the struct default value from the struct initialization function.

   To resolve these problems, we can modify the test implementation to use the struct initialization function instead of setting the struct fields on the test cases.

   For example, the current implementation may look like this:

   ```go
   func Test_server_SaveIndex(t *testing.T) {
        ...

        // all the fields to initialize the target struct
        type fields struct {
            name              string
            ip                string
            ngt               service.NGT
            eg                errgroup.Group
            streamConcurrency int
        }

        ...

        for _, tc := range tests {
            test := tc
            t.Run(test.name, func(tt *testing.T) {
                tt.Parallel()
                ...

                // we initialize the target struct using fields defined in test case
                s := &server{
                    name:              test.fields.name,
                    ip:                test.fields.ip,
                    ngt:               test.fields.ngt,
                    eg:                test.fields.eg,
                    streamConcurrency: test.fields.streamConcurrency,
                    // we may have more fields defined in the struct,
                    // and we need to set them all here
                }

                gotRes, err := s.SaveIndex(test.args.ctx, test.args.in1)
                ...
   ```

   We can customize the implementation to use the struct initialization function.

   ```go
   func Test_server_SaveIndex(t *testing.T) {
        ...

        // we customize the fields to initialize the target struct
        type fields struct {
            // options to configure server struct
            srvOpts []Option

            // options to configure ngt
            svcCfg  *config.NGT
            svcOpts []service.Option
        }

        ...

        for _, tc := range tests {
            test := tc
            t.Run(test.name, func(tt *testing.T) {
                tt.Parallel()
                ...

                // initialize required fields of server struct
                eg, _ := errgroup.New(ctx)
                ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts, service.WithErrGroup(eg))...)
                if err != nil {
                    tt.Errorf("failed to init ngt service, error = %v", err)
                }

                // we use the struct initialization function instead of setting the fields,
                // by combining the use of the functional option, we can configure the struct,
                // and the default values of the struct will be set automatically if we are not defining it.
                s, err := New(append(test.fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
                if err != nil {
                    tt.Errorf("failed to init server, error= %v", err)
                }

                gotRes, err := s.SaveIndex(test.args.ctx, test.args.in1)
                ...
   ```

### Parallel test

In Vald, we use parallel tests to accelerate the execution of tests by default.
There are two layers of enabling parallel tests.

1. Parallel for the test function
1. Parallel for the subtests in the test function

The generated test case will enable these two parallel modes by default. It is implemented by:

```go
func Test_server_CreateIndex(t *testing.T) {
    t.Parallel() // parallel for the test function
    type args struct {

    ...

    for _, tc := range tests {
        test := tc
        t.Run(test.name, func(tt *testing.T) {
            tt.Parallel() // parallel for subtests
    ...
```

Be careful of using the parallel test, avoid using share object used in the test, to avoid race detector to detect the race in test and fail the test.

### Helper function

In golang testing package, they defined the helper functions in testing types to indicate the current function is a test helper function and skip printing the file and line information to the output.

In Vald, we can apply it to the different helper functions like `beforeFunc()` or `afterFunc()` defined in the test case struct.

```go
type test struct {
    ...
    beforeFunc func(*testing.T) // helper function to initialze testing
    afterFunc  func(*testing.T) // helper function to cleanup
}

tests := []test {
    {
        name: "send success when host and port are correct value",
        args: args {
            host: "vdaas.vald.org",
            port: "80",
        },
        field: field {
            host: "vdaas.vald.org",
            port: "80",
        },
        beforeFunc: func(t *testing.T) {
            t.Helper() // indicate the current function is helper function
            /* initialization logic of testing environment */
        },
        want: want {
            err: nil,
        },
        afterFunc: func(t *testing.T) {
            t.Helper() // indicate the current function is helper function
            /* cleanup logic */
        }
    },
}
```

### Using Mock

In Vald, we use a lot of external libraries, and there are a lot of dependencies between libraries.

As a result, due to the complexity, it has become hard to determine whether or not to mock dependencies.

#### Condition

When dependencies have the following factor, you can decide to mock the dependencies.

- Incomplete implementation
- I/O
  - e.g., Network access, disk operation, etc.
- Hardware dependent
  - e.g., CPU, Memory usage, disk I/O, etc.
- Hard to create the error of dependencies
- Difficult to initialize
  - e.g., Random number and time, file I/O initialization, environment dependent, etc.
- The test result may change in each runtime
  - e.g., the only test result may change in each runtime, System call inside implementation, etc.

#### Risk

Before applying mock to the object, you should understand the following risks.

- We **do not** know whether the dependencies are correctly implemented or not.
- We cannot notice the changes in dependencies.

#### Implementation

The implementation of the mock object should be:

- Same package as the mock target.
- File name is `xxx_mock.go`
- Struct name is `Mock{Interface name}`

For example, we decided to mock the following implementation `Encoder`.

```go
package json

type Encoder interface {
    Encode(interface{}) ([]byte, error)
}
```

```go
type encoder struct {
    encoder json.Encoder
}

func (e *encoder) Encode(obj interface{}) ([]byte, error) {
    return e.encoder.Encode(obj)
}
```

The following is an example of mock implementation:

```go
package json

type MockEncoder struct {
    EncoderFunc func(interface{}) ([]byte, error)
}

func (m *MockEncoder) Encode(obj interface{}) ([]byte, error) {
    return m.EncodeFunc(obj)
}
```

The following is an example implementation of test code to create the mock object and mock the implementation.

```go
tests := []test {
    {
        name: "returns (byte{}, nil) when encode success"
        fields: fields {
            encoding: &json.MockEncoder {
                EncoderFunc: func(interface{}) ([]byte, error) {
                    return []byte{}, nil
                },
            },
        }
        ......
    }
}
```
