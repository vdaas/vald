# Vald coding guideline

## Introduction
This guideline includes the coding style for all Vald contributors and reviewers. Everyone should follow this guideline to keep the style consistent so everyone can understand and contribute to Vald easier once they learn this guideline. You should have the basic knowledge of how to write Golang before contributing to Vald. If you found any bug please create a GitHub issue and we will work on it.
This guideline is based on [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md), [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments), [Gitlab Go standard and style guideline](https://docs.gitlab.com/ee/development/go_guide/) and [Effective Go](https://golang.org/doc/effective_go.html).
For the guideline to write test code please refer to [here](xxxxxxxxx).

## Code Formatting and Naming Convension
Code formatting and naming conventions affect coding readability and maintainability. Every developer has a different coding style, luckily Golang provides tools to format source code and checking for the potential issue in the source code. We suggest using [gofmt](https://golang.org/cmd/gofmt/) to format the source code in Vald, and [golint](https://github.com/golang/lint). We suggest everyone install the plugin for your editor to automatically format the code once you edit the code.
But having a tools to format source code doesn't mean you do not need to care the formatting of the code, for example:
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
The project layout includes the folder and the file structure in the project. Basically we follow the [prject-layout](https://github.com/golang-standards/project-layout) example in Vald.

### Packages
The package defines the context of the objects in the package, for example the corresponding methods and structs belongs to corresponding package. Unlike other languages like Java, in Golang we use the package name to declar which context of the object we are going to use. For example in [time](https://golang.org/pkg/time/) package, it defines all the objects about time like `time.Now()` method to get the current time.

Here is the naming conventions of the package:
- Package name should be the same as the folder name.
- Package name should keep as simple as it should, and should contain only one specific context in the package.
- Package name should not be too general, for example `util` or `helper`, which will cause all the objects from different contexts to be store in one package. If you really want to name the package as `util`,  please define the more specific package  name more  `ioutil` or `httputil`.

All packages should contains `doc.go` file under the package to describe what is the package is. For example, under the folder name called `cache` should contains a file named `doc.go`, which contains the package documentation. For example

```golang
// Package cache provides implementation of cache
package cache
````

### Interfaces
Interface defines the program interface for usability and future extendability.
Unlike other languages like Java, golang support implicit interface implementation. The type implements do not need to specify the interface name; to "implments" the interface the structs only needs to defined the methods same as the interface, so please be careful to define the method name inside the interface.

Here is the naming conventions of the interface:
- Use MixedCaps

```golang
type RoundTripper interface {
    // interface definition
}
```

### Structs
Structs in golang is the object definition, we can attach any fields and methods to the struct.

Here is the naming conventions of the struct:
- Use MixedCaps

#### Struct initialization
There are many ways to initialize structs in Golang, base on the use case we can decide which way to initialize objects in Golang.

### Variables and Constant

### Methods

#### Getter and Setter

#### Defer functions

## Program comments
Program comments makes the code more easier to understand. Basically we suggest not to write many comments inside the source code, unless the source code is very complicated and confusing; otherwise we should divide the source code into methods to keep the readability and usability of the source code.

Everyone should write the comments to all the public objects on your source code, like public packages, interface, structs, methods, and even public constant and variable. The godoc will be generated base on the comment of source code.

## Documentation
Documentation is generated based on the program comments. Please refer to [godoc](https://godoc.org/github.com/vdaas/vald) for the program documentation.

## Internal packages
Vald implement its own internal package to extend the functionality of the standard library and third-party library. Please refer to [godoc](https://godoc.org/github.com/vdaas/vald/internal) for the internal package document.

## Dependency management and Build
We should use `go mod tidy` to manage the `go.mod` file in the project.
