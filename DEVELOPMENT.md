# Development documentation

This document describes how to setup the development environment and how to develop Vald.

## Set up environment

### Prerequisites

#### OS

- When using Docker related environment, you can use any OS that supports Docker.
- When using native environment, Linux is required.

#### Architecture

`amd64` is recommended because `NGT`, the vector search library we use, is optimized for `amd64`.
But you can also build and test `Vald` on `arm64` with the same way as described below.

### Devcontainer

This is the easiest way to start developing `Vald`. You can just open our [devcontainer.json](./.devcontainer/devcontainer.json) with `VS Code` and go.

### Other

We don't officially have a setup documentation for now, but you can take a look at the [`Dockerfile`](./dockers/dev/Dockerfile).
That's everything you need to build and test `Vald`, so you can use it as a reference.

> If you would like to use the `Dockerfile` directlly, please note that `docker-in-docker` environment is required to run our E2E tests.
> In devcontainer, [`VS Code` handles it for us](https://github.com/devcontainers/features/tree/main/src/docker-in-docker).

## Run tests

### Unit tests

### E2E tests
