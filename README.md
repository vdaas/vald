<div align="center">
<img src="./assets/image/logo.svg" width="50%">
</div>

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![CLA](https://cla-assistant.io/readme/badge/vdaas/vald?&style=flat-square)](https://cla-assistant.io/vdaas/vald)
[![release](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![codecov](https://img.shields.io/codecov/c/github/vdaas/vald.svg?style=flat-square)](https://codecov.io/gh/vdaas/vald)
[![Codacy Badge](https://img.shields.io/codacy/grade/a6e544eee7bc49e08a000bb10ba3deed?style=flat-square)](https://www.codacy.com/app/i.can.feel.gravity/vald?utm_source=github.com&utm_medium=referral&utm_content=vdaas/vald&utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdaas/vald?style=flat-square)](https://goreportcard.com/report/github.com/vdaas/vald)
[![GolangCI](https://golangci.com/badges/github.com/vdaas/vald.svg?style=flat-square)](https://golangci.com/r/github.com/vdaas/vald)
[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/vdaas/vald)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/vdaas/vald)
[![DepShield Badge](https://depshield.sonatype.org/badges/vdaas/vald/depshield.svg?style=flat-square)](https://depshield.github.io)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvdaas%2Fvald.svg?type=small&style=flat-square)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvdaas%2Fvald?ref=badge_small)
[![Total visitor](https://visitor-count-badge.herokuapp.com/total.svg?repo_id=vald)](https://github.com/vdaas/vald/graphs/traffic)
[![Visitors in today](https://visitor-count-badge.herokuapp.com/today.svg?repo_id=vald)](https://github.com/vdaas/vald/graphs/traffic)

vald is high scalable distributed high-speed approximate nearest neighbor search engine

## Requirement

kubernetes 1.12~

## Installation

### Using Helm

```shell
helm repo add vald https://vald.vdaas.org/charts
helm install --generate-name vald/vald
```

If you use the default values.yaml, the `nightly` images will be installed.

#### Docker image tagging policy

- `nightly` ... latest build of master branch
- `vX.X.X` ... released versions
- `latest` ... latest build of release versions
- `stable` ... latest long-term supported version

### Using Helm-operator

```shell
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/master/k8s/helm-operator/serviceaccount.yaml
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/master/k8s/helm-operator/clusterrole.yaml
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/master/k8s/helm-operator/clusterrolebinding.yaml
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/master/k8s/helm-operator/operator.yaml
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/master/k8s/helm-operator/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/master/k8s/helm-operator/cr.yaml
```

## Example

```shell
Write example here
```

## Architecture Overview

<div align="center">
<img src="./assets/image/Vald Architecture Overview.svg" width="100%">
</div>

## Development

Before your first commit to this repository, it is strongly recommended to run the commands below.

```shell
make init
```

## Contribution

Please read the [contribution guide](https://github.com/vdaas/vald/blob/master/CONTRIBUTING.md)

## Author

- [kpango](https://github.com/kpango)
- [kmrmt](https://github.com/kmrmt)
- [rinx](https://github.com/rinx)

## Contributor

- [hlts2](https://github.com/hlts2)

## LICENSE

vald released under Apache 2.0 license, refer [LICENSE](https://github.com/vdaas/vald/blob/master/LICENSE) file.

[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/vdaas/vald/?ref=repository-badge)

<table>
  <tr>
    <th>component</th>
    <th>implementation</th>
    <th>Docker name</th>
    <th>Docker build status</th>
  </tr>
  <tr>
    <td>agent</td>
    <td>NGT</td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-agent-ngt">vdaas/vald-agent-ngt</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20agent-ngt/badge.svg"></td>
  </tr>
  <tr>
    <td>discoverer</td>
    <td>K8s</td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-discoverer-k8s">vdaas/vald-discoverer-k8s</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20discoverer-k8s/badge.svg"></td>
  </tr>
  <tr>
    <td>gateway</td>
    <td></td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-gateway">vdaas/vald-gateway</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20gateway-vald/badge.svg"></td>
  </tr>
  <tr>
    <td rowspan=2>backup manager</td>
    <td>MySQL</td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-manager-backup-mysql">vdaas/vald-manager-backup-mysql</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20backup-manager-mysql/badge.svg"></td>
  </tr>
  <tr>
    <td>Cassandra</td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-manager-backup-cassandra">vdaas/vald-manager-backup-cassandra</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20backup-manager-cassandra/badge.svg"></td>
  </tr>
  <tr>
    <td>compressor</td>
    <td></td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-manager-compressor">vdaas/vald-manager-compressor</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20manager-compressor/badge.svg"></td>
  </tr>
  <tr>
    <td rowspan=2>meta</td>
    <td>Redis</td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-meta-redis">vdaas/vald-meta-redis</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20meta-redis/badge.svg"></td>
  </tr>
  <tr>
    <td>Cassandra</td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-meta-cassandra">vdaas/vald-meta-cassandra</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20meta-cassandra/badge.svg"></td>
  </tr>
  <tr>
    <td>index manager</td>
    <td></td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-manager-index">vdaas/vald-manager-index</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20manager-index/badge.svg"></td>
  </tr>
  <tr>
    <td>helm-operator</td>
    <td></td>
    <td><a href="https://hub.docker.com/r/vdaas/vald-helm-operator">vdaas/vald-helm-operator</a></td>
    <td><img src="https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20helm-operator/badge.svg"></td>
  </tr>
</table>

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvdaas%2Fvald.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvdaas%2Fvald?ref=badge_large)
