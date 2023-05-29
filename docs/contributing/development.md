# Development documentation

This document describes how to set up the development environment and how to develop Vald.

## Set up environment

### Prerequisites

#### OS

- When using Docker related environment, you can use any OS that supports Docker.
- When using native environment, `Linux` is required.

#### Architecture

`amd64` is recommended because `NGT`, the vector search library we use, is optimized for `amd64`.
But you can also build and test `Vald` on `arm64` with the same way as described below.

### Devcontainer

This is the easiest way to start developing `Vald`. You can just open our [devcontainer.json](https://github.com/vdaas/vald/blob/main/.devcontainer/devcontainer.json) with `VS Code` and go.

### Other

We don't officially have a setup documentation for now, but you can take a look at the [`Dockerfile`](https://github.com/vdaas/vald/blob/main/dockers/dev/Dockerfile).
That's everything you need to build and test `Vald`, so you can use it as a reference.

> If you would like to use the `Dockerfile` directlly, please note that `docker-in-docker` environment is required to run our E2E tests.
> In devcontainer, [`VS Code` handles it for us](https://github.com/devcontainers/features/tree/main/src/docker-in-docker).

## Run tests

### Unit tests

The command below will run all the unit tests.

```bash
make test
```

### E2E tests

The steps below will deploy `Vald` to local `k3d` cluster and run the E2E tests.

1. Change `example/helm/values.yaml` to `dimensions: 784` and `distance_type: l2`.
2. Run the commands below.

```bash
# Download the dataset
make hack/benchmark/assets/dataset/fashion-mnist-784-euclidean.hdf5

# Start k3d
make k3d/start

# Wait for a while until the cluster is ready
# You might want to use k9s for this

# Deploy Vald
make k8s/vald/deploy

# Wait for a while until the deployment is ready

# Run E2E tests
make e2e E2E_WAIT_FOR_CREATE_INDEX_DURATION=3m

# The result will be shown in three minutes or so

# Delete the cluster
make k8s/vald/delete
```
