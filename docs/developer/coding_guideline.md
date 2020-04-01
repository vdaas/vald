# Vald coding guideline

## Introduction

This guideline includes the coding style for all Vald contributors and reviewers. Everyone should following this guideline to keep the style consistent so everyone can understand and contribute to Vald easier once they learn this guideline. You should have the basic knowledge of how to write Golang before contributing to Vald. If you found any bug please create a github issue and we will working on it.
This guideline is based on [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md), [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments), [Gitlab Go standard and style guideline](https://docs.gitlab.com/ee/development/go_guide/) and [Effective Go](https://golang.org/doc/effective_go.html).
For the guideline to write test code please refer to [here](xxxxxxxxx).

## Code Formatting
Code formatting is very important, it affects the code readability and maintainability. Every developers have different coding style. Luckily Golang provide tools to format the code, and we suggest using [gofmt](https://golang.org/cmd/gofmt/) to format the source code in Vald. We suggest everyone install the plugin for your editor to automatically format the code once you edit the code.
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

It is very hard to read and maintain, in this case we should use the following style.
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

## Naming convension
This section describe the naming rules of every objects of the source code in Vald. Please refer the sub-section for the detail naming convension of each objects in Vald.

### Source code file name

### Package name

### Interface name

### Struct name

### Function name

#### Getter and Setter

### Variable name

## Program comments
Program comments makes the code more easier to understand. Basically we suggest not to write many comments inside the source code, unless the source code is very complicated and confusing; otherwise we should divide the source code into functions to keep the readability of the source code.

Everyone should write the comments to all the public objects on your source code, like public packages or interface.

## Initialization
There are many ways to initialize objects in Golang. Base on the use case we can decide which style to follow to initialize objects in Golang.
