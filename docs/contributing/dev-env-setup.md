# Development environment setup

Before making any changes to Vald, you need to install the required tools to make changes to Vald.

In this section, we will describe how to install the required tools and the required steps before making changes to Vald.

Please note that if you have already setup the development environment, you do not need to do it again.

## Setup GitHub

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

## Windows / Mac

For Mac / Windows users, please download and install [docker](https://www.docker.com/) to create a Linux container to build and test Vald.

You need to install the packages and execute the commands listed in the [Linux section](#linux-users) inside the docker container rather than in your local environment.

If you want to start development on your host environment rather than in the docker container, please consider mounting the Vald repository path from the host environment to the container environment after cloning the Vald repository described below.

```bash
docker run -v '{vald repo}':'{folder mount}' {container name}
```

For more details about docker, please refer to the [docker documentation](https://docs.docker.com/get-started/overview/).

## Linux

For Linux users, please install the following tools on your environment.

- [curl](https://curl.se/)
- [make](https://www.gnu.org/software/make/)
- [cmake](https://cmake.org/)
- [Protobuf](https://grpc.io/docs/protoc-installation/)
- [npm](https://www.npmjs.com/)
- [unzip](https://linux.die.net/man/1/unzip)
- [Go](https://go.dev/) (v1.19.5 or later is recommended)

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