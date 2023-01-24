# Contributing guide

Thank you for your interest in Vald, and thank you for investing your time in contributing to Vald!

In this guide, you will get an idea of how to contribute to Vald.

This guide is for everyone who wants to contribute to Vald. Even if you are not a developer, don't worry, some contributions don't require writing a single line of code.

Please read our [Code Of Conduct](https://github.com/vdaas/vald/blob/main/CODE_OF_CONDUCT.md) to keep our community approachable and respectable.

Before you make any contribution to Vald, please read the [About Vald](https://vald.vdaas.org/docs/overview/about-vald) to get an overview of Vald.

## Type of contribution

Vald is an open-source project, everyone can contribute to Vald.

We accept any kind of contribution, including the following types of contributions:

- Issue

  - Bug report
  - Security issue report
  - Feature request / Proposal

- Pull request
  - Source code implementation
    - Business logic implementation
    - Test implementation
  - Documentation

Please note that you can make a contribution not only listed above, any voice is meaningful and important to us.

Please feel free to contact us on [Slack](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA). We are waiting for you :)

## Change flow

Before making any changes, you may need to understand the standard flow of making changes.

If you found any bugs or security issues in Vald, or you want to request a new feature, please follow [this section](#issue-contribution) to contribute to the issue.

If you decided to make a pull request contribution, including source code changes, or documentation changes, you need to understand the overall flow to make changes, including:

1. Create a development branch
2. Make code changes
3. Test your changes
4. Commit your changes
5. Create a pull request

We will describe the details in [this section](#pull-request-contribution).

## Issue contribution

[Issues](https://github.com/vdaas/vald/issues) are used to track tasks that contributors can help with.

Please find [here](./docs/contributing/issue.md) for more details about how to contribute to the issue.

## Pull request contribution

Pull request is also called merge request to let others know about the changes you have made, to review and discuss the changes made and finally merge it to the main branch.

In Vald, you need to create a pull request to ask for the review and actually make changes to Vald.

In this section, we will describe what you need to do to make a pull request contribution to Vald.

### Development environment setup

Before making any changes to Vald, you need to install the required tools to make changes to Vald.

In this section, we will describe how to install the required tools and the required steps before making changes to Vald.

Please note that if you have already setup the development environment, you do not need to do it again.

#### Setup GitHub

Please install [Git](https://git-scm.com/) and configure it first.

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald).
1. Set your name and email (these should match the information on your submitted CLA).

   ```bash
   git config --global user.name "GitHub user name"
   git config --global user.email "your_email@example.com"
   ```

   Please also refer [here](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup) for more details on setting up Git.

1. Setup the signing key on your development environment.

   Please refer to [here](https://docs.github.com/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key) to configure the signing key.

   Vald recommends signing the commit to prove that the commit actually came from you. The reason is that it is easy to add anyone as an author of the commit, which can be used in hiding the author of malicious code.

1. Fork Vald repository.

   Please [fork Vald repository](https://github.com/vdaas/vald/fork) to copy Vald repository to your own GitHub organization. It allows you to make changes to it without affecting Vald repository.

#### Windows / Mac

For Mac / Windows users, please download and install [docker](https://www.docker.com/) to create a Linux container to build and test Vald.

You need to install the packages and execute the commands listed in the [Linux section](#linux-users) inside the docker container rather than in your local environment.

If you want to start development on your host environment rather than in the docker container, please consider mounting the Vald repository path from the host environment to the container environment after cloning the Vald repository described below.

```bash
docker run -v '{vald repo}':'{folder mount}' {container name}
```

For more details about docker, please refer to the [docker documentation](https://docs.docker.com/get-started/overview/).

#### Linux

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

Please clone Vald repository to the GOPATH, and setup the remote branch to your forked Vald repository.

```bash
# clone vdaas repo
mkdir -p $(go env GOPATH)/src/github.com/vdaas/
cd $(go env GOPATH)/src/github.com/vdaas/
git clone https://github.com/vdaas/vald.git
cd vald

# rename origin repo to upstream and set origin to remote forked repo
git remote rename origin upstream
git remote add origin {your forked repo}
git fetch origin
```

Please also execute the following command to initialize the development environment and install the necessary packages and tools.

```bash
make init # initialize development environment, and install NGT
make tools/install # install development tools like helm, kind, etc.
make gotests/install # install gotests tools to generate test stubs.
```

### Vald code structure

Before making changes, you need to understand Vald code structure in order to know which part of the code you needed to make changes.

Vald is mainly written in Golang. If you are going to contribute to the source code, please find below the description of the packages.

| Package name | Description                                                      |
| :----------- | :--------------------------------------------------------------- |
| apis         | Vald API definition                                              |
| cmd          | Entry point of the Vald components                               |
| example      | Example code of Vald                                             |
| internal     | Internal package to extend and customize libraries functionality |
| pkg          | Contains business logic implementation                           |
| hack         | Contains implementation to change the behavior of tools          |
| tests        | Contains test implementation except for unit tests like e2e test |

Other than source code, you may need to make changes on manifest files of Vald such as documentation, helm configuration files, proto files, GitHub CI/CD manifest, etc.

Please find below the description of the folder architecture with manifest files.

| Folder name   | Description                                                    |
| :------------ | :------------------------------------------------------------- |
| .devcontainer | Contains development container manifest                        |
| .github       | GitHub CI/CD settings                                          |
| Makefile.d    | Make command definition                                        |
| assets        | Contains assertion files                                       |
| charts        | Contains Helm charts files and configuration                   |
| design        | Contains Vald design documentation                             |
| dockers       | Contains docker files of all components in Vald                |
| docs          | Contains Vald documentation                                    |
| k8s           | Contains all example k8s manifest of all components in Vald    |
| versions      | Contains version definition of Vald and third-party components |

### Make command

Vald provides a different `make` command to help you make changes to Vald. Different actions in the `make` command provide different functionality for the user.

Use the following command. `make [action]` to execute an action of the make command.

Here is some useful action of the make command.

| Action           | Description                                                                                        |
| :--------------- | :------------------------------------------------------------------------------------------------- |
| format           | Format go, yaml, markdown, and json files, and generate the license header on the source code file |
| binary/build     | Build all Vald components into an executable                                                       |
| docker/build     | Build all Vald components into docker images                                                       |
| proto/all        | Rebuild all proto files and generate source code files                                             |
| helm/schema/vald | Generate json schema for Vald Helm Chart                                                           |
| helm/docs/vald   | Generate Helm documentation                                                                        |
| gotests/gen      | Execute gotests tools to generate unit test code stubs                                             |
| test             | Execute unit test on cmd, pkg, and internal packages                                               |
| bench            | Execute the benchmarking on NGT, Vald agent, and Vald LB gateway                                   |

For more actions and details about our make commands, please find our [Makefile](https://github.com/vdaas/vald/blob/main/Makefile) and [Makefile.d](https://github.com/vdaas/vald/tree/main/Makefile.d).

### CI/CD

Currently, Vald contains mainly 3 types of CI/CD pipelines running on GitHub action:

- Build and deploy pipelines
- Testing pipelines
- Linter pipelines

And these pipelines are executed when a pull request is created.

| Pipeline name                            | Description                                           |
| :--------------------------------------- | :---------------------------------------------------- |
| Build docker image: {image name}         | Pipeline to build docker image of specific image      |
| Coverage                                 | Execute code test coverage report                     |
| Run tests / Run tests for {package name} | Run unit test on the specific package                 |
| Run e2e {target}                         | Execute End-to-end testing                            |
| Run Helm lint / {chart name}             | Execute linter for Helm chart files                   |
| reviewdog - {target}                     | Execute reviewdog to review the changes on the target |
| DeepSource - {target}                    | Execute DeepSource to analysis target files           |

About the `Build docker image` pipeline, it will build the docker image for the specific tag and upload it to the [DockerHub](https://hub.docker.com/u/vdaas).

Specifically, it builds the docker image on every pull request and tags the image as `pr-{pull request number}`.

Whenever a pull request is merged into the main branch, a nightly build will be built and uploaded to the DockerHub.

When Vald release, this pipeline will build the docker image with a specific version (e.g. `v1.6.1`) and also update the `latest` tag image.

About the `Run e2e` pipeline, it will execute the end-to-end testing and fails if any error occurred.

This pipeline will be executed only if the label `e2e/{target}` is added to the pull request.

### Create development branch

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

### Make changes

In this section, we will briefly describe how to make changes to Vald.

Before making changes, please make sure you know what you are doing and know about which part of the source code is required to change, please do not change the code which is unnecessary.

#### Code changes

If you are making source code changes, please follow and understand our [coding guideline](./docs/contributing/coding-style.md) to follow our coding style to keep the coding style consistent.

If you have discussed the design and the requirement of the changes with Vald members, please follow the design and the requirement discussed.

After your changes are made, we suggest you execute the following command to generate the necessary test stubs and format the source code.

```bash
make gotests/gen # execute gotests tools to generate unit test code stubs
make format # format go and yaml files
```

#### Document & image changes

In Vald, the documents are written in Markdown format and stored in the `docs` folder. Eventually, documents will be deployed to the [Vald official website](https://vald.vdaas.org/).

In Vald, we support some styles specific to the official website to display different types of content.

We provide the following CSS class for the website:

| CSS class name | Description                              |
| :------------- | :--------------------------------------- |
| caution        | Display the content as a caution message |
| warning        | Display the content as a warning message |
| note           | Display the content as a note message    |

To apply the above CSS class, quote the sentence to the div with the CSS class name.

For example, to quote the sentence as a warning message:

```html
<div class="warning">This is warning message!!</div>
```

About the images in the document, we are using [diagrams.net](https://www.diagrams.net/) to draw the image.

All the images of the document are stored in the `assets/docs` directory, and each of the images should contain a `.drawio` extension file and the `svg` or `png` file.

This `.drawio` file is the source file and can only be opened on [diagrams.net](https://www.diagrams.net/).

Please find [here](https://www.diagrams.net/features) for more about how to use diagrams.net.

After creating or modifying the image on diagrams.net, please store and update the `.drawio`, and export the image to the `.svg` file to the Vald repository.

#### Test your changes

Please make sure to test and validate your changes before adding and committing your changes.

For code changes, please refer to [this document](./docs/contributing/testing-guideline.md) for more detail about how to test your changes.

### Add and commit changes

After verifying the changes, you need to add your changes and push them to your development branch.

Please add the files that are related to your changes only.

```bash
git add [file1] [file2] ...
```

And please write a brief description of the commit message, and push it to your forked repository.

```bash
git commit --signoff -m '[commit message]'
git push origin [type]/[area]/[description]
```

### Create Pull Request

After committing your changes, you may create a pull request to ask for accepting the changes.

Please create the pull request to the Vald repository under `vdaas` organization.

Each pull request and commit should be small enough to contain only one purpose, for easier review and tracking.
Please fill in the description on the pull request and write down the overview of the changes.

#### Pull Request labels

Labels in the pull request indicate what kind of the changes is, and provide some extra feature to the pull request.

Please choose the correct type label on the pull request, we provide the following type label in Vald:

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

Please add the action label to the pull request to execute specific actions if you want to execute them.

#### Confirmation item of Pull Request

Before asking for a review, please make sure to check the following things:

- CI/CD pipeline works properly
- Linter warnings and comments

About CI/CD, please confirm CI/CD works properly, especially the build pipelines and test pipelines.

It helps to ensure Vald is working properly to ensure the quality of deliverables.

About Linter warnings and comments, Vald is integrated with different linters to check the quality of Vald.

Please check all of the warnings and comments from linters, review them one by one, and fix them if necessary.

#### Review

Vald team will review the pull request.

We may ask for changes or questions we have during the review process.

We will add a mention to your GitHub account on each comment and reply to make the communication smooth.

After the review is done, we will merge the pull request to Vald. Your changes will be applied to Vald, and the changes will be included in the next Vald release.
