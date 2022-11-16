# Contributing to Vald

Thank you for your interest in Vald, and thank you for investing your time in contributing to Vald!

Read our [Code Of Conduct](https://github.com/vdaas/vald/blob/main/CODE_OF_CONDUCT.md) to keep our community approachable and respectable.

In this guide, you will get an idea of how to contribute to Vald.

If you are not a developer, don't worry, some contributions don't require writing a single line of code.

## Contributor guide

Please read the [README](https://github.com/vdaas/vald/blob/main/README.md) to get an overview of Vald.

We welcome you to contribute to Vald to make Vald better. We accept the following types of contributions:

- Issue
  - Bug report
  - Feature request
  - Proposal
  - Security issue report
- Pull request

Please also feel free to ask anything on [Vald Slack channel](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA) :)

If you decided to contribute source code changes, please follow the following guideline to install the required tools.

We recommend using Linux to build and test Vald.
For MacOS / Windows users, please install and use `docker` to create a Linux container to build and test Vald. You may need to install the packages and execute the commands listed below inside the docker container rather than in your local environment.
If you want to start development on your host environment, please also consider mounting the Vald repository path from the host environment to the container environment after cloning the Vald repository described below.

```bash
docker run -v '{vald repo}':'{folder mount}' {container name}
```

Please install the following tools on your environment.

- [make](https://www.gnu.org/software/make/)
- [cmake](https://cmake.org/)
- [Protobuf](https://grpc.io/docs/protoc-installation/)
- [npm](https://www.npmjs.com/)
- [unzip](https://linux.die.net/man/1/unzip)
- [Git](https://git-scm.com/)
- [Go](https://go.dev/) (v1.19.2 is recommended)

For Debian-based Linux distribution users, you can install these required tools using `apt`.

```bash
sudo apt install make cmake protobuf-compiler npm unzip git golang
```

Please [fork Vald repository](https://github.com/vdaas/vald/fork) to your repository and clone your Vald repository to your Go path.

```bash
# clone vdaas repo
mkdir -p $(go env GOPATH)/src/github.com/vdaas/
cd $(go env GOPATH)/src/github.com/vdaas/
git clone https://github.com/vdaas/vald.git
cd vald

# rename origin repo to upstream and set origin to remote folked repo
git remote rename origin upstream
git remote add origin {your forked repo}
git fetch origin
```

Please also run the following command to initialize the development environment and install the necessary packages and tools.

```bash
make init # initialize development environment, and install NGT
make tools/install # install development tools like helm, kind, etc.
make gotests/install # install gotests tools to generate test stubs.
```

## Issue

[Issues](https://github.com/vdaas/vald/issues) are used to track tasks that contributors can help with.

An issue can be a bug report, feature request, or security issue report.

We welcome you to give us ideas to improve Vald or report any issue that exists in Vald by creating an issue on the Vald repository.

Please find [here](./docs/contributing/issue.md) for more details about the issue, and please also feel free to ask anything on [Vald Slack channel](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA).

### Propose a new issue

If you have found something that should be updated in Vald, search open issues to see if someone else has already talked about it.

If it is something new, please open a new issue using the following template:

- [Bug report](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=)
- [Feature request](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Ffeature%2C+priority%2Flow%2C+team%2Fcore&template=feature_request.md&title=)
- [Security issue report](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fsecurity%2C+priority%2Fmedium%2C+team%2Fcore%2C+team%2Fsre&template=security_issue_report.md&title=)

Please fill in the information written on the template to help us to understand what you think.

### Solve an issue

Please find our [opening issues](https://github.com/vdaas/vald/issues) and find one that interested you.

You can find the issue using the label `type/*`. For example, you can find all the bug issues using the label `type/bug`.

## Make changes

Before making changes on Vald, please follow these steps to contribute to any of our open-source repositories:

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald)
2. Set your name and email (these should match the information on your submitted CLA)

    ```bash
    git config --global user.name "Firstname Lastname"
    git config --global user.email "your_email@example.com"
    ```

    Please also refer [here](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup) for more details on setting up Git.
3. Setup signing key on your development environment
    Please refer [here](https://docs.github.com/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key) to configure the signing key.
    Vald recommends signing the commit to prove that the commit actually came from you, as it is easy to add anyone as an author of the commit, which can be used in hiding the author of malicious code.

### How to make changes

1. Comment on the issue and say you are working on it to avoid conflict with others also working on the same issue.
2. Fork the Vald repository.
    - You need to fork the Vald repository only once, please see [contributor guideline](#contributor-guide) for more details
3. Create your feature branch on your forked repository.
    - Please refer to [this section](#Branch-naming-convention) for the branch naming convention.

    ```bash
    git checkout -b [type]/[area]/[description]
    ```

4. Make code changes.
    - Please follow the design or requirement discussed on the issue.
    - Make sure you understand our [coding guideline](./docs/contributing/coding-style.md) to follow our coding style to keep the coding style consistent
    - After making code changes, please refer [here](#after-making-code-changes) to format the source code and generate the test code
5. Verify your changes.
    - If you are making code changes, please refer to [this section](#test-your-changes)
6. Add the updated files to the branch.
    - Please only add the files related to your changes

    ```bash
    git add [files]
    ```

7. Commit your changes to the branch.
    - Please write a brief description of the changes to this commit.

    ```bash
    git commit --signoff -m '[commit message]'
    ```

8. Push to the forked branch.

    ```bash
    git push origin [type]/[area]/[description]
    ```

9.  Create a new pull request against the Vald repository.
     - Please also mention the issue on the pull request if needed.
10.  Wait for the code review.
     - Resolve any issue/questions raised by reviewers until it is merged.

### Branch naming convention

Before working on changes, you need to create a development branch on your forked branch.

Name the development branch  `[type]/[area]/[description]`.

| Field       | Explanation                           | Naming Rule                                                                                                               |
| :---------- | :------------------------------------ | :------------------------------------------------------------------------------------------------------------------------ |
| type        | The PR type                           | The type of PR can be a feature, bug, refactoring, benchmark, security, documentation, dependencies, ci, test, etc...    |
| area        | Area of context                       | The area of PR can be gateway, agent, agent-sidecar, lb-gateway, etc...                                                |
| description | Summarized description of your branch | The description must be hyphenated. Please use [a-zA-Z0-9] and a hyphen as characters, and do not use any other characters. |

(\*) If you changed multiple areas, please list each area with "-".

For example, when you add a new feature for internal/servers, the name of the branch will be `feature/internal/add-newfeature-for-servers`.

### After making code changes

After making code changes, we suggest you execute the following command to generate the necessary test stubs and format code.

```bash
make gotests/gen # execute gotests tools to generate unit test code stubs
make format # format go and yaml files
```

The command `make gotests/gen` generate unit test code stubs to easier to implement unit test code. Please see [this section](#unit-test) for more details.

The command `make format` is used to generate the license header on the source code file, and execute the code formatter to format Go and YAML files.

It will also install the following tools to format the source code.

- [golines](https://github.com/segmentio/golines)
- [gofumpt](https://github.com/mvdan/gofumpt)
- [strictgoimports](https://github.com/momotaro98/strictgoimports)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [prettier](https://prettier.io/)

These tools are required to format your source code to keep the coding style consistent.

### Test your changes

Testing your changes is very important to ensure your implementation is working as expected.

#### Unit test

Unit test is used to check whether the unit is implemented correctly in various cases.

We suggest you implement or update the unit test code when making logical changes or implementing new functionality in Vald.

Before implementing the unit test code, we suggest you read our [unit test guideline](./docs/contributing/unit-test-guideline.md) to guide you to create good unit tests and [coding guideline for unit test](./docs/contributing/coding-style.md#test) to guide you to implement unit tests.

If you want to execute the unit test on only part of the code, you can use `go test` command to execute the unit test on the specific package/function.
For example, if you want to execute the unit test on a specific package, use the following command.

```bash
go test -race [package]
```

This command will execute the unit test on the package, and also enable the race detector to check if any race occurs in the implementation.

If you want to execute the unit test on the whole Vald implementation, Vald provides the following command to do that.

```bash
make test
```

This command will execute all unit tests of `*target*_test.go` files on `cmd`, `internal` and `pkg` packages. It is useful to ensure that your changes will not affect the behavior of other components and packages.

#### End-To-End (E2E) testing

End-To-End (E2E) testing is used to test the application flow of Vald is working as expected from beginning to end.

If you want to execute E2E test on Vald, Vald provides the following commands to test the implementation.

```bash
make e2e
```

The command `make e2e` execute E2E tests to ensure whether the functionality is working as expected. It will perform the actual CRUD action on a cluster and verify the result.

E2E tests require deploying Vald on a Kubernetes cluster beforehand. You can deploy Vald on your Kubernetes cluster, or you can create a Kubernetes cluster on your local machine easily by using the tools like [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/).

Please refer to our [get started](./docs/tutorial/get-started.md) to create the cluster and deploy Vald on a Kubernetes cluster.

If you want to execute E2E test on your Kubernetes cluster, you may need to modify the configuration on [Makefile](https://github.com/vdaas/vald/blob/main/Makefile) before executing the E2E test.

| Config name                        | Description |
| :--------------------------------- | :---------- |
| E2E_BIND_HOST                      | The target host of Kubernetes cluster |
| E2E_BIND_PORT                      | The target port of Kubernetes cluster |
| E2E_TIMEOUT                        | The timeout of E2E test |
| E2E_DATASET_NAME                   | The dataset name of the E2E test|
| E2E_INSERT_COUNT                   | The number of index insert in E2E test |
| E2E_SEARCH_COUNT                   | The number of search request in E2E test |
| E2E_SEARCH_BY_ID_COUNT             | The number of search by ID request in E2E test |
| E2E_GET_OBJECT_COUNT               | The number of get object request in E2E test |
| E2E_UPDATE_COUNT                   | The number of update request in E2E test |
| E2E_UPSERT_COUNT                   | The number of upsert request in E2E test |
| E2E_REMOVE_COUNT                   | The number of remove request in E2E test |
| E2E_WAIT_FOR_CREATE_INDEX_DURATION | The wait time of create index operation after insert is completed in E2E test |
| E2E_TARGET_NAME                    | The target pod name in the Vald cluster |
| E2E_TARGET_NAMESPACE               | The target namespace of the Vald cluster |
| E2E_TARGET_PORT                    | The pod forward port of the target pod to the local host |
| E2E_PORTFORWARD_ENABLED            | Enable/Disable port forwarding |

### Pull request

After making code changes and testing your changes, you may create a pull request to ask for accepting the changes.

Each pull request and commit should be small enough to contain only one purpose, for easier review and tracking.
Please fill in the description on the pull request and write down the overview of the changes.

Please also choose the correct type label on the pull request, we provide the following type label in Vald:

| Label            | Description                          |
| :--------------- | :----------------------------------- |
| type/bug         | For bug fixes pull request           |
| type/dependency  | For dependency update pull request   |
| type/feature     | For new feature pull request         |
| type/refactoring | For code refactoring pull request    |
| type/security    | For security fix pull request        |
| type/test        | For test implementation pull request |

We also provide the following label to execute specific actions on the [GitHub Actions](https://github.co.jp/features/actions).

| Label                | Description                                           |
| :------------------- | :---------------------------------------------------- |
| action/e2e-chaos     | Execute E2E chaos test                                |
| action/e2e-deploy    | Execute E2E deployment test                           |
| action/e2e-max-dim   | Execute maximum dimension E2E test                    |
| action/e2e-profiling | Execute E2E test with profiling                       |
| action/fossa         | Execute [fossa](https://fossa.com/) security checking |

If you are solving an issue, please also link the pull request to the issue.

Vald team will review the pull request. We may ask for changes or questions we have before the pull request is merged.

After the pull request is merged, your changes will be applied to Vald, and the changes will be included in the next Vald release.
