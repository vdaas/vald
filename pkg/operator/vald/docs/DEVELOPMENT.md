# Development Guide

This directory contains the controller only: it reconciles `ValdOperatorRelease` resources and
generates `ValdRelease` (VRS) objects. Provisioning a target cluster and installing the
Vald Helm Operator (VHO) are out of scope.

## Prerequisites

- Go 1.23+
- Docker 17.03+
- `kubectl` 1.11.3+ and access to a Kubernetes 1.11.3+ cluster
- A working Vald Helm Operator (VHO) in the target cluster to turn generated
  `ValdRelease` objects into running Vald pods

## What this controller does

It reconciles `ValdOperatorRelease` (`vor`, `vald.vdaas.org/v1`) resources and generates
`ValdRelease` (VRS) objects. See [explain.md](./explain.md) for the design, the CRD spec,
the reconcile lifecycle, and the full environment-variable reference.

## Common tasks

```sh
# Build and run unit tests
make build
make test

# Regenerate CRDs / RBAC / deepcopy after changing api/v1 types
make manifests generate

# Lint
make lint

# Build and push the controller image
make docker-build docker-push IMG=<registry>/valdoperatorrelease:tag

# Install CRDs and deploy the controller to the current kube-context
make install
make deploy IMG=<registry>/valdoperatorrelease:tag
```

Run `make help` for the full list of targets.

> `config/` is kubebuilder-generated. Do not hand-edit generated manifests; change the
> kubebuilder markers in the Go sources and regenerate, then use `make build-installer`
> to produce the bundled `install.yaml`.

## Remaining OSS-readiness work

Tracked in the repository's issue tracker. Highlights:

- Rename the `vor` / `valdoperatorrelease` identifiers for OSS.
- Reuse `vdaas/vald` packages (types, logging, errors, k8s) instead of the current
  self-contained implementations.
- Relicense file headers to the vald header format.
