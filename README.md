<div align="center">
<img src="./assets/image/svg/logo.svg" width="50%">
</div>

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![CLA](https://cla-assistant.io/readme/badge/vdaas/vald?&style=flat-square)](https://cla-assistant.io/vdaas/vald)
[![release](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Codacy Badge](https://img.shields.io/codacy/grade/a6e544eee7bc49e08a000bb10ba3deed?style=flat-square)](https://www.codacy.com/app/i.can.feel.gravity/vald?utm_source=github.com&utm_medium=referral&utm_content=vdaas/vald&utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdaas/vald?style=flat-square)](https://goreportcard.com/report/github.com/vdaas/vald)
[![GolangCI](https://golangci.com/badges/github.com/vdaas/vald.svg?style=flat-square)](https://golangci.com/r/github.com/vdaas/vald)
[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/vdaas/vald)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/vdaas/vald)
[![DepShield Badge](https://depshield.sonatype.org/badges/vdaas/vald/depshield.svg?style=flat-square)](https://depshield.github.io)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvdaas%2Fvald.svg?type=small&style=flat-square)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvdaas%2Fvald?ref=badge_small)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/vdaas/vald/?ref=repository-badge)
[![Contributors](https://img.shields.io/github/contributors/vdaas/vald?style=flat-square)](https://github.com/vdaas/vald/graphs/contributors)
[![Slack](https://img.shields.io/badge/slack-join-brightgreen?logo=slack&style=flat-square)](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)
<!--[![codecov](https://img.shields.io/codecov/c/github/vdaas/vald.svg?style=flat-square&logo=codecov)](https://codecov.io/gh/vdaas/vald) -->

## What is Vald?

Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine.

Vald is designed and implemented based on Cloud-Native architecture.

It uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors.

Vald has automatic vector indexing and index backup, and horizontal scaling which made for searching from billions of feature vector data.

Vald is easy to use, feature-rich and highly customizable as you needed.

Go to [Get Started](./docs/user/get-started.md) page to try out Vald :)

(If you are interested in ANN benchmarks, please refer to [the official website](http://ann-benchmarks.com/).)<br>

### Main Features

- Asynchronous Auto Indexing
    - Usually the graph requires locking during indexing, which causes stop-the-world. But Vald uses distributed index graphs so it continues to work during indexing.

- Customizable Ingress/Egress Filtering
    - Vald implements it's own highly customizable Ingress/Egress filter.
    - Which can be configured to fit the gRPC interface.
        - Ingress Filter: Ability to Vectorize through filter on request.
        - Egress Filter: rerank or filter the searching result with your own algorithm.

- Cloud-native based vector searching engine
    - Horizontal scalable on memory and CPU for your demand.

- Auto Backup for Index data
    - Vald has a feature to store the backup of the index data using MySQL or Cassandra which enables disaster recovery.

- Distributed Indexing
    - Vald distribute vector index to multiple agents, each agent stores different index.

- Index Replication
    - Vald stores each index in multiple agents which enables index replicas.
    - Automatically rebalance the replica when some Vald agent goes down.

- Easy to use
    - Vald can be easily installed in a few steps.

- Highly customizable
    - You can configure the number of vector dimensions, the number of replica and etc.

- Multi language supported
    - Golang, Java, Clojure, Node.js, and Python client library are supported.
    - gRPC APIs can be triggered by any programming languages which support gRPC.
    - REST API is also supported.```

## Requirement

kubernetes 1.17~

## Get Started

Please refer to [Get Started](./docs/user/get-started.md).

## Installation

### Using Helm

```shell
helm repo add vald https://vald.vdaas.org/charts
helm install vald-cluster vald/vald
```

If you use the default values.yaml, the `nightly` images will be installed.

#### Docker image tagging policy

- `nightly` ... latest build of master branch
- `vX.X.X` ... released versions
- `latest` ... latest build of release versions
- `stable` ... latest long-term supported version

### Using Helm-operator

[vald-helm-operator](https://github.com/vdaas/vald/blob/master/charts/vald-helm-operator)

## Example

```shell
Write example here
```

## Architecture Overview

<div align="center">
<img src="./assets/image/svg/Vald Architecture Overview.svg" width="100%">
</div>

Please refer [here](./docs/overview/architecture.md) for more details of the architecture overview in the future.

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
- [vankichi](https://github.com/vankichi)
- [kevindiu](https://github.com/kevindiu)


<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://kpango.com"><img src="https://avatars1.githubusercontent.com/u/9798091?v=4" width="100px;" alt=""/><br /><sub><b>Yusuke Kato</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kpango" title="Code">💻</a> <a href="#design-kpango" title="Design">🎨</a> <a href="https://github.com/vdaas/vald/commits?author=kpango" title="Documentation">📖</a> <a href="#example-kpango" title="Examples">💡</a> <a href="#fundingFinding-kpango" title="Funding Finding">🔍</a> <a href="#ideas-kpango" title="Ideas, Planning, & Feedback">🤔</a> <a href="#infra-kpango" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-kpango" title="Maintenance">🚧</a> <a href="#platform-kpango" title="Packaging/porting to new platform">📦</a> <a href="#plugin-kpango" title="Plugin/utility libraries">🔌</a> <a href="#projectManagement-kpango" title="Project Management">📆</a> <a href="#question-kpango" title="Answering Questions">💬</a> <a href="https://github.com/vdaas/vald/pulls?q=is%3Apr+reviewed-by%3Akpango" title="Reviewed Pull Requests">👀</a> <a href="#security-kpango" title="Security">🛡️</a> <a href="#tool-kpango" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=kpango" title="Tests">⚠️</a> <a href="#tutorial-kpango" title="Tutorials">✅</a> <a href="#talk-kpango" title="Talks">📢</a></td>
    <td align="center"><a href="https://github.com/rinx"><img src="https://avatars3.githubusercontent.com/u/1588935?v=4" width="100px;" alt=""/><br /><sub><b>Rintaro Okamura</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=rinx" title="Code">💻</a> <a href="https://github.com/vdaas/vald/commits?author=rinx" title="Documentation">📖</a> <a href="#example-rinx" title="Examples">💡</a> <a href="#ideas-rinx" title="Ideas, Planning, & Feedback">🤔</a> <a href="#infra-rinx" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-rinx" title="Maintenance">🚧</a> <a href="#platform-rinx" title="Packaging/porting to new platform">📦</a> <a href="#plugin-rinx" title="Plugin/utility libraries">🔌</a> <a href="#question-rinx" title="Answering Questions">💬</a> <a href="https://github.com/vdaas/vald/pulls?q=is%3Apr+reviewed-by%3Arinx" title="Reviewed Pull Requests">👀</a> <a href="#tool-rinx" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=rinx" title="Tests">⚠️</a> <a href="#tutorial-rinx" title="Tutorials">✅</a> <a href="#talk-rinx" title="Talks">📢</a></td>
    <td align="center"><a href="https://morimoto.dev/"><img src="https://avatars2.githubusercontent.com/u/413873?v=4" width="100px;" alt=""/><br /><sub><b>Kosuke Morimoto</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kmrmt" title="Code">💻</a> <a href="#example-kmrmt" title="Examples">💡</a> <a href="https://github.com/vdaas/vald/pulls?q=is%3Apr+reviewed-by%3Akmrmt" title="Reviewed Pull Requests">👀</a> <a href="#tool-kmrmt" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=kmrmt" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/vankichi"><img src="https://avatars3.githubusercontent.com/u/13959763?v=4" width="100px;" alt=""/><br /><sub><b>Kiichiro YUKAWA</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=vankichi" title="Documentation">📖</a> <a href="#example-vankichi" title="Examples">💡</a> <a href="#maintenance-vankichi" title="Maintenance">🚧</a> <a href="https://github.com/vdaas/vald/pulls?q=is%3Apr+reviewed-by%3Avankichi" title="Reviewed Pull Requests">👀</a> <a href="#tool-vankichi" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=vankichi" title="Tests">⚠️</a> <a href="#tutorial-vankichi" title="Tutorials">✅</a></td>
    <td align="center"><a href="https://github.com/datelier"><img src="https://avatars3.githubusercontent.com/u/57349093?v=4" width="100px;" alt=""/><br /><sub><b>datelier</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=datelier" title="Code">💻</a> <a href="#ideas-datelier" title="Ideas, Planning, & Feedback">🤔</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

## LICENSE

vald released under Apache 2.0 license, refer [LICENSE](https://github.com/vdaas/vald/blob/master/LICENSE) file.

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

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
