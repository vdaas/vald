# Contributing to Vald

Thank you for your interest in Vald, and thank you for investing your time in contributing to Vald!

Please read our [Code Of Conduct](https://github.com/vdaas/vald/blob/main/CODE_OF_CONDUCT.md) to keep our community approachable and respectable.

In this guide, you will get an idea of how to contribute to Vald.

If you are not a developer, don't worry, some contributions don't require writing a single line of code.

## Contributions

Please read the [About Vald](https://vald.vdaas.org/docs/overview/about-vald) to get an overview of Vald.

We welcome you to contribute to Vald to make Vald better.
We accept the following types of contributions:

- Issue
  - Bug report
  - Feature request / Proposal
  - Security issue report
- Pull request
  - Source code implementation
    - Business logic implementation
    - Test implementation
  - Documentation

Please also feel free to ask anything on [Vald Slack channel](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA) :)

## Issue contribution

[Issues](https://github.com/vdaas/vald/issues) are used to track tasks that contributors can help with.

Please find [here](./docs/contributing/issue.md) for more details about how to contribute to the issue.

## Pull request contribution

Pull request is also called merge request to let other knows about the changes you have made, to review and discuss the changes and finally merge it to the main branch.

In Vald, you need to create a pull request to ask for the review and actually make changes to Vald.

You need to create pull request to make changes on:

1. Source code changes
   1. Business logic implementation
   2. Test implementation
2. Documentation creation / update

### Install and configure GitHub

Before creating a pull request, you need to install and configure GitHub to create development branch, make changes and finally create a pull request.

Please install [Git](https://git-scm.com/) and configure it first.

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald).
1. Set your name and email (these should match the information on your submitted CLA).

   ```bash
   git config --global user.name "GitHub user name"
   git config --global user.email "your_email@example.com"
   ```

   Please also refer [here](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup) for more details on setting up Git.

1. Setup signing key on your development environment.

   Please refer [here](https://docs.github.com/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key) to configure the signing key.

   Vald recommends signing the commit to prove that the commit actually came from you, as it is easy to add anyone as an author of the commit, which can be used in hiding the author of malicious code.

After install and configure Git, please clone Vald repository to your Go path, [fork Vald repository](https://github.com/vdaas/vald/fork) and setup the remote branch to your forked Vald repository.

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

### Document contribution

In Vald, the document is written as markdown files mainly under [docs](https://github.com/vdaas/vald/tree/main/docs) directory, and the asset files like image files is placed on [assets/docs](https://github.com/vdaas/vald/tree/main/assets/docs).

For the image files, we use [diagrams.net](https://app.diagrams.net/) to draw and modify the image files.

#### Installation guide for document contributor

For making changes on document (Markdown files), we recommand using [Visual Studio Code](https://code.visualstudio.com/) to make changes on Markdown files.
Please follow the link above to install it to your development environment.

We also recommend you to install a plugin for Visual Studio Code called [Markdown All in One](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one) to enable some extra features for easier modification on Markdown file.

### Source code contribution

If you decided to contribute source code to Vald, include business logic changes and test case implementation, you need to understand the package structure of Vald and know which part you need to make changes.

| Package name | Description                                                      |
| :----------- | :--------------------------------------------------------------- |
| apis         | Vald API definitation                                            |
| cmd          | Entry point of the Vald components                               |
| example      | Example code of Vald                                             |
| internal     | Internal package to extend and customize libraries functionality |
| pkg          | Contains business logic implementation                           |

Please make sure you understand our [coding guideline](./docs/contributing/coding-style.md) to follow our coding style to keep the coding style consistent.

#### Installation guide for source code contributor

We recommend using Linux environment to contribute code.
Please find the below sections to install the required tools on your environment.

##### Mac/Windows users

For Mac / Windows users, please install and use [docker](https://www.docker.com/) to create a Linux container to build and test Vald.

You need to install the packages and execute the commands listed in the [Linux section](#linux-users) inside the docker container rather than in your local environment.

If you want to start development on your host environment rather than in the docker container, please consider mounting the Vald repository path from the host environment to the container environment after cloning the Vald repository described below.

```bash
docker run -v '{vald repo}':'{folder mount}' {container name}
```

For more details about docker, please refer to the [docker documentation](https://docs.docker.com/get-started/overview/).

##### Linux users

For Linux users, please install the following tools on your environment.

- [curl](https://curl.se/)
- [make](https://www.gnu.org/software/make/)
- [cmake](https://cmake.org/)
- [Protobuf](https://grpc.io/docs/protoc-installation/)
- [npm](https://www.npmjs.com/)
- [unzip](https://linux.die.net/man/1/unzip)
- [Go](https://go.dev/) (v1.19.2 or later is recommended)

For Debian-based Linux distribution users, you can install these required tools using `apt`.

```bash
sudo apt install curl make cmake protobuf-compiler npm unzip git golang
```

Please also run the following command under your Vald repository to initialize the development environment and install the necessary packages and tools.

```bash
make init # initialize development environment, and install NGT
make tools/install # install development tools like helm, kind, etc.
make gotests/install # install gotests tools to generate test stubs.
```

### Make changes

1. Make sure no one is working on the same issue/feature

   Before making any changes, you need to check if anyone is working on the same feature in the pull request list.

   If you are solving an issue, check if anyone is working on the issue and comment on the issue and say you are working on it to avoid conflict with others also working on the same issue.

2. Create your feature branch on your forked repository

   Before working on changes, you need to create a development branch on your forked branch.

   Name the development branch `[type]/[area]/[description]`.

   | Field       | Explanation                           | Naming Rule                                                                                                                 |
   | :---------- | :------------------------------------ | :-------------------------------------------------------------------------------------------------------------------------- |
   | type        | The PR type                           | The type of PR can be a feature, bug, refactoring, benchmark, security, documentation, dependencies, ci, test, etc...       |
   | area        | Area of context                       | The area of PR can be gateway, agent, agent-sidecar, lb-gateway, etc...                                                     |
   | description | Summarized description of your branch | The description must be hyphenated. Please use [a-zA-Z0-9] and a hyphen as characters, and do not use any other characters. |

   (\*) If you changed multiple areas, please list each area with "-".

   For example, when you add a new feature for internal/servers, the name of the branch will be `feature/internal/add-newfeature-for-servers`.

   ```bash
   git checkout -b [type]/[area]/[description]
   ```

3. Make changes on Vald

   If you have discussed the design and the requirement of the changes with Vald members, please follow the design and the requirement discussed.

4. After making source code changes

   If you are making source code changes, we suggest you execute the following command to generate the necessary test stubs and format code.

   ```bash
   make gotests/gen # execute gotests tools to generate unit test code stubs
   make format # format go and yaml files
   ```

   The command `make gotests/gen` generate unit test code stubs to easier to implement unit test code.

   The command `make format` is used to generate the license header on the source code file, and execute the code formatter to format Go and YAML files.

5. Verify the changes

   If you are making logical changes on Vald, please refer to [this document](./docs/contributing/testing-guideline.md) for more detail about how to test your changes.

6. Commit and push your changes to the branch

   After verifing the changes, you may want to push the changes to your development branch.

   Please add the files that related to your changes only.

   ```bash
   git add [files]
   ```

   Please write a brief description of the changes to the commit, and push it to your forked repository.

   ```bash
   git commit --signoff -m '[commit message]'
   git push origin [type]/[area]/[description]
   ```

7. Create a new pull request against the Vald repository

   After committing your changes, you may create a pull request to ask for accepting the changes.

   Please create the pull request to the Vald repository under `vdaas` orginization.

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

   Please add the action label to the pull request to execute specific action if needed.

   If you are solving an issue, please also link the pull request to the issue.

8. Review and merge pull request

   Vald team will review the pull request.

   We may ask for changes or questions we have during the review process.

   We will add a mention to your GitHub account on each comment and reply to make the communication smooth.

   After the review is done, we will merge the pull request to Vald. Your changes will be applied to Vald, and the changes will be included in the next Vald release.
