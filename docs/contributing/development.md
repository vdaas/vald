# Development documentation

This document describes how to set up the development environment and how to develop Vald.

## Set up environment

### Prerequisites

#### OS

- When using Docker related environment, you can use any OS that supports Docker.
- When using native environment, `Linux` or `macOS` (Darwin; `arm64`/Apple Silicon and `amd64`) are supported.

#### Architecture

`amd64` is recommended because `NGT`, the vector search library we use, is optimized for `amd64`.
But you can also build and test `Vald` on `arm64` with the same way as described below.

### Devcontainer

This is the easiest way to start developing `Vald`. You can just open our [devcontainer.json](https://github.com/vdaas/vald/blob/main/.devcontainer/devcontainer.json) with `VS Code` and go.

### macOS (Darwin) native build

Vald can be built and tested natively on macOS (Apple Silicon and Intel). It uses
**Apple clang** (from the Xcode Command Line Tools) as the C/C++ compiler,
[Homebrew](https://brew.sh) for the remaining native dependencies, and Apple's
`Accelerate.framework` for BLAS/LAPACK. The C/C++ libraries (NGT, Faiss, HDF5,
zlib) are built from source into `/usr/local` by the Makefile, which requires
privileged install steps.

```bash
# 1. Xcode Command Line Tools — REQUIRED: provides Apple clang and the macOS SDK
xcode-select --install
clang --version    # verify: should print "Apple clang version ..."

# 2. Homebrew dependencies
brew install go protobuf buf cmake hdf5 zlib libomp

# 3. Resolve Homebrew prefixes as the current user, then let Make elevate only
#    the filesystem install commands. Do not run the entire Makefile as root.
OPENMP_PREFIX="$(brew --prefix libomp)" \
HDF5_PREFIX="$(brew --prefix hdf5)" \
ZLIB_PREFIX="$(brew --prefix zlib)" \
SUDO=sudo make ngt/install hdf5/install faiss/install

# 4. Run the unit tests
make test
```

> **Compiler:** the build uses Apple clang from step 1 — `brew install llvm` is
> **not** required. (Homebrew `llvm` is only needed if you explicitly want LLVM's
> `clang`/`llvm-ar`; in that case `brew install llvm` and add
> `$(brew --prefix llvm)/bin` to your `PATH` first.) On Darwin the Makefile
> resolves OpenMP and HDF5/ZLIB headers through the prefixes above, uses Apple
> thin-LTO (`-flto=thin`), and `Accelerate.framework` for BLAS/LAPACK — so
> `OpenBLAS`/`lapack`/`gcc` (gfortran) are **not** required. On Linux the original
> GNU toolchain flags are unchanged.

### Other

We don't officially have a setup documentation for now, but you can take a look at the [`Dockerfile`](https://github.com/vdaas/vald/blob/main/dockers/dev/Dockerfile).
That's everything you need to build and test `Vald`, so you can use it as a reference.

> If you would like to use the `Dockerfile` directly, please note that `docker-in-docker` environment is required to run our E2E tests.
> In devcontainer, [`VS Code` handles it for us](https://github.com/devcontainers/features/tree/main/src/docker-in-docker).

## Run tests

### Unit tests

The command below will run all the unit tests.

```bash
make test
```

### E2E tests

The steps below will deploy a Vald cluster to the local `k3d` cluster and run the E2E tests.

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
make k8s/vald/deploy HELM_VALUES=example/helm/values.yaml

# Wait for a while until the deployment is ready

# Run E2E tests
make e2e E2E_WAIT_FOR_CREATE_INDEX_DURATION=3m

# The result will be shown in three minutes or so

# Delete the cluster
make k8s/vald/delete
```
