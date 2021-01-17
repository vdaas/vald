<div align="center">
<img src="./assets/image/svg/logo.svg" width="50%">
</div>

[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/vdaas/vald.svg)](https://pkg.go.dev/github.com/vdaas/vald)
[![Codacy Badge](https://img.shields.io/codacy/grade/a6e544eee7bc49e08a000bb10ba3deed?style=flat-square)](https://www.codacy.com/app/i.can.feel.gravity/vald?utm_source=github.com&utm_medium=referral&utm_content=vdaas/vald&utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdaas/vald?style=flat-square)](https://goreportcard.com/report/github.com/vdaas/vald)
[![DepShield Badge](https://depshield.sonatype.org/badges/vdaas/vald/depshield.svg?style=flat-square)](https://depshield.github.io)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B21465%2Fvald.svg?type=small)](https://app.fossa.com/projects/custom%2B21465%2Fvald?ref=badge_small)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/vdaas/vald/?ref=repository-badge)
[![CLA](https://cla-assistant.io/readme/badge/vdaas/vald?&style=flat-square)](https://cla-assistant.io/vdaas/vald)
[![Artifact Hub](https://img.shields.io/badge/chart-ArtifactHub-informational?logo=helm&style=flat-square)](https://artifacthub.io/packages/chart/vald/vald)
[![Slack](https://img.shields.io/badge/slack-join-brightgreen?logo=slack&style=flat-square)](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)
<!--[![codecov](https://img.shields.io/codecov/c/github/vdaas/vald.svg?style=flat-square&logo=codecov)](https://codecov.io/gh/vdaas/vald) -->

## What is Vald?

Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine.

Vald is designed and implemented based on Cloud-Native architecture.

It uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors.

Vald has automatic vector indexing and index backup, and horizontal scaling which made for searching from billions of feature vector data.

Vald is easy to use, feature-rich and highly customizable as you needed.

Go to [Get Started](./docs/tutorial/get-started.md) page to try out Vald :)

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
    - Go, Java, Clojure, Node.js, and Python client library are supported.
    - gRPC APIs can be triggered by any programming languages which support gRPC.
    - REST API is also supported.```

## Requirements

- Kubernetes 1.17~
- AVX2 instructions (required by Vald Agent NGT)

## Get Started

Please refer to [Get Started](./docs/tutorial/get-started.md).

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

## Components
<table>
  <tr>
    <th>Component</th>
    <th>Docker image</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-agent-ngt">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-agent-ngt?label=vdaas%2Fvald-agent-ngt&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-agent-ngt">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--agent--ngt-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Agent Sidecar</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-agent-sidecar">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-agent-sidecar?label=vdaas%2Fvald-agent-sidecar&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-agent-sidecar">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--agent--sidecar-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Discoverer K8s</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-discoverer-k8s">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-discoverer-k8s?label=vdaas%2Fvald-discoverer-k8s&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-discoverer-k8s">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--discoverer--k8s-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-gateway">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-gateway?label=vdaas%2Fvald-gateway&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-gateway">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--gateway-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Backup Manager (MySQL)</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-manager-backup-mysql">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-manager-backup-mysql?label=vdaas%2Fvald-manager-backup-mysql&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-manager-backup-mysql">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--manager--backup--mysql-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Backup Manager (Cassandra)</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-manager-backup-cassandra">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-manager-backup-cassandra?label=vdaas%2Fvald-manager-backup-cassandra&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-manager-backup-cassandra">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--manager--backup--cassandra-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-manager-compressor">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-manager-compressor?label=vdaas%2Fvald-manager-compressor&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-manager-compressor">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--manager--compressor-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Meta (Redis)</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-meta-redis">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-meta-redis?label=vdaas%2Fvald-meta-redis&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-meta-redis">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--meta--redis-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Meta (Cassandra)</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-meta-cassandra">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-meta-cassandra?label=vdaas%2Fvald-meta-cassandra&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-meta-cassandra">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--meta--cassandra-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-manager-index">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-manager-index?label=vdaas%2Fvald-manager-index&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-manager-index">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--manager--index-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-helm-operator">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-helm-operator?label=vdaas%2Fvald-helm-operator&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-helm-operator">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--helm--operator-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
  <tr>
    <td>Load Test</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-loadtest">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-loadtest?label=vdaas%2Fvald-loadtest&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-loadtest">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--loadtest-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
  </tr>
</table>

## Contribution

Please read the [contribution guide](https://github.com/vdaas/vald/blob/master/CONTRIBUTING.md)

## Contributors
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-9-orange.svg?style=flat-square)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):
<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://kpango.com"><img src="https://avatars1.githubusercontent.com/u/9798091?v=4" width="100px;" alt=""/><br /><sub><b>Yusuke Kato</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kpango" title="Code">💻</a> <a href="#ideas-kpango" title="Ideas, Planning, & Feedback">🤔</a> <a href="#maintenance-kpango" title="Maintenance">🚧</a> <a href="#projectManagement-kpango" title="Project Management">📆</a></td>
    <td align="center"><a href="https://github.com/rinx"><img src="https://avatars3.githubusercontent.com/u/1588935?v=4" width="100px;" alt=""/><br /><sub><b>Rintaro Okamura</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=rinx" title="Code">💻</a> <a href="https://github.com/vdaas/vald/commits?author=rinx" title="Documentation">📖</a> <a href="#maintenance-rinx" title="Maintenance">🚧</a> <a href="#platform-rinx" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://morimoto.dev/"><img src="https://avatars2.githubusercontent.com/u/413873?v=4" width="100px;" alt=""/><br /><sub><b>Kosuke Morimoto</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kmrmt" title="Code">💻</a> <a href="#example-kmrmt" title="Examples">💡</a> <a href="#tool-kmrmt" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=kmrmt" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/vankichi"><img src="https://avatars3.githubusercontent.com/u/13959763?v=4" width="100px;" alt=""/><br /><sub><b>Kiichiro YUKAWA</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=vankichi" title="Documentation">📖</a> <a href="#maintenance-vankichi" title="Maintenance">🚧</a> <a href="https://github.com/vdaas/vald/commits?author=vankichi" title="Tests">⚠️</a> <a href="#tutorial-vankichi" title="Tutorials">✅</a></td>
    <td align="center"><a href="https://github.com/datelier"><img src="https://avatars3.githubusercontent.com/u/57349093?v=4" width="100px;" alt=""/><br /><sub><b>datelier</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=datelier" title="Code">💻</a> <a href="#ideas-datelier" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/kevindiu"><img src="https://avatars1.githubusercontent.com/u/1985382?v=4" width="100px;" alt=""/><br /><sub><b>Kevin Diu</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kevindiu" title="Documentation">📖</a> <a href="#example-kevindiu" title="Examples">💡</a> <a href="https://github.com/vdaas/vald/commits?author=kevindiu" title="Tests">⚠️</a> <a href="#tutorial-kevindiu" title="Tutorials">✅</a></td>
    <td align="center"><a href="https://twitter.com/hiroto_hlts2"><img src="https://avatars0.githubusercontent.com/u/25459661?v=4" width="100px;" alt=""/><br /><sub><b>Hiroto Funakoshi</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=hlts2" title="Documentation">📖</a> <a href="#tool-hlts2" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=hlts2" title="Tests">⚠️</a> <a href="#tutorial-hlts2" title="Tutorials">✅</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/taisuou"><img src="https://avatars0.githubusercontent.com/u/21119375?v=4" width="100px;" alt=""/><br /><sub><b>taisho</b></sub></a><br /><a href="#design-taisuou" title="Design">🎨</a> <a href="https://github.com/vdaas/vald/commits?author=taisuou" title="Documentation">📖</a> <a href="#example-taisuou" title="Examples">💡</a></td>
    <td align="center"><a href="https://github.com/pgrimaud"><img src="https://avatars1.githubusercontent.com/u/1866496?v=4" width="100px;" alt=""/><br /><sub><b>Pierre Grimaud</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=pgrimaud" title="Documentation">📖</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

## LICENSE

vald released under Apache 2.0 license, refer [LICENSE](https://github.com/vdaas/vald/blob/master/LICENSE) file.

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B21465%2Fvald.svg?type=large)](https://app.fossa.com/projects/custom%2B21465%2Fvald?ref=badge_large)
