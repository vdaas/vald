<div align="center">
<a href="https://vald.vdaas.org/">
    <img src="./assets/image/readme.svg" width="50%" />
</a>
</div>

[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/vdaas/vald.svg)](https://pkg.go.dev/github.com/vdaas/vald)
[![Codacy Badge](https://img.shields.io/codacy/grade/a6e544eee7bc49e08a000bb10ba3deed?style=flat-square)](https://www.codacy.com/app/i.can.feel.gravity/vald?utm_source=github.com&utm_medium=referral&utm_content=vdaas/vald&utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdaas/vald?style=flat-square)](https://goreportcard.com/report/github.com/vdaas/vald)
[![DepShield Badge](https://depshield.sonatype.org/badges/vdaas/vald/depshield.svg?style=flat-square)](https://depshield.github.io)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B21465%2Fvald.svg?type=small)](https://app.fossa.com/projects/custom%2B21465%2Fvald?ref=badge_small)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/vdaas/vald/?ref=repository-badge)
[![DeepSource](https://deepsource.io/gh/vdaas/vald.svg/?label=resolved+issues&show_trend=true&token=UpNEsc0zsAfGw-MPPa6O05Lb)](https://deepsource.io/gh/vdaas/vald/?ref=repository-badge)
[![CLA](https://cla-assistant.io/readme/badge/vdaas/vald?&style=flat-square)](https://cla-assistant.io/vdaas/vald)
[![Artifact Hub](https://img.shields.io/badge/chart-ArtifactHub-informational?logo=helm&style=flat-square)](https://artifacthub.io/packages/chart/vald/vald)
[![Slack](https://img.shields.io/badge/slack-join-brightgreen?logo=slack&style=flat-square)](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

<!--[![codecov](https://img.shields.io/codecov/c/github/vdaas/vald.svg?style=flat-square&logo=codecov)](https://codecov.io/gh/vdaas/vald) -->

## What is Vald?

Vald is a highly scalable distributed fast approximate nearest neighbor (ANN) dense vector search engine.

Vald is designed and implemented based on Cloud-Native architecture.

Vald has automatic vector indexing and index backup, and horizontal scaling which made for searching from billions of feature vector data.

Vald is easy to use, feature-rich and highly customizable as you needed.

It uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors.

(If you are interested in ANN benchmarks, please refer to [ann-benchmarks.com](https://ann-benchmarks.com/).)

For more information, please refer to [Official Web Site](https://vald.vdaas.org).

<div align="center">
  <img src="./assets/image/svg/vald_architecture_overview.svg" width="100%" />
</div>

Vald can handle any object data, image, audio processing, video, text, binary, or etc., if converting to the vector, and be used for:

- Recognition
- Recommendation
- Detecting
- Grammar checker
- Real-time translator
- also you want to do!

## Requirements

- Kubernetes 1.19~
- AVX2 instructions (required by Vald Agent NGT)

## Get Started

Go to [Get Started](https://vald.vdaas.org/docs/tutorial/get-started) page to try out Vald !

## Installation

### Using Helm

```shell
helm repo add vald https://vald.vdaas.org/charts
helm install vald-cluster vald/vald
```

If you use the default values.yaml, the `nightly` images will be installed.

### Using Helm-operator

Please refer to [vald-helm-operator](https://github.com/vdaas/vald/blob/main/charts/vald-helm-operator).

## Components

<table>
  <tr>
    <th>Component</th>
    <th>Docker image</th>
    <th>latest image</th>
    <th>nightly image</th>
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
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-agent-ngt/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-agent-ngt/latest?label=vald-agent-ngt" />
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-agent-ngt/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-agent-ngt/nightly?label=vald-agent-ngt" />
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
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-agent-sidecar/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-agent-sidecar/latest?label=vald-agent-sidecar" />
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-agent-sidecar/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-agent-sidecar/nightly?label=vald-agent-sidecar" />
      </a>
    </td>
  </tr>
  <tr>
    <td>Discoverer</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-discoverer-k8s">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-discoverer-k8s?label=vdaas%2Fvald-discoverer-k8s&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-discoverer-k8s">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--discoverer--k8s-brightgreen?logo=docker&style=flat-square"/>
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-discoverer-k8s/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-discoverer-k8s/latest?label=vald-discoverer-k8s" />
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-discoverer-k8s/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-discoverer-k8s/nightly?label=vald-discoverer-k8s" />
      </a>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-lb-gateway">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-lb-gateway?label=vdaas%2Fvald-lb-gateway&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-lb-gateway">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--lb--gateway-brightgreen?logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://hub.docker.com/r/vdaas/vald-filter-gateway">
        <img src="https://img.shields.io/docker/pulls/vdaas/vald-filter-gateway?label=vdaas%2Fvald-filter-gateway&logo=docker&style=flat-square"/>
      </a><br/>
      <a href="https://github.com/orgs/vdaas/packages/container/package/vald/vald-filter-gateway">
        <img src="https://img.shields.io/badge/ghcr.io-vdaas%2Fvald%2Fvald--filter--gateway-brightgreen?logo=docker&style=flat-square"/>
      </a><br/>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-lb-gateway/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-lb-gateway/latest?label=vald-lb-gateway" />
      </a><br />
      <a href="https://hub.docker.com/r/vdaas/vald-filter-gateway/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-filter-gateway/latest?label=vald-filter-gateway" />
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-lb-gateway/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-lb-gateway/nightly?label=vald-lb-gateway" />
      </a><br>
      <a href="https://hub.docker.com/r/vdaas/vald-filter-gateway/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-filter-gateway/nightly?label=vald-filter-gateway" />
      </a><br />
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
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-manager-index/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-manager-index/latest?label=vald-index-manager" />
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-manager-index/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-manager-index/nightly?label=vald-index-manager" />
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
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-helm-operator/tags?page=1&name=latest">
        <img src="https://img.shields.io/docker/v/vdaas/vald-helm-operator/latest?label=vald-helm-operator" />
      </a>
    </td>
    <td>
      <a href="https://hub.docker.com/r/vdaas/vald-helm-operator/tags?page=1&name=nightly">
        <img src="https://img.shields.io/docker/v/vdaas/vald-helm-operator/nightly?label=vald-helm-operator" />
      </a>
    </td>
  </tr>
</table>

Docker images tagging policy:

- `nightly` ... latest build of main branch
- `vX.X.X` ... released versions
- `latest` ... latest build of release versions
- `stable` ... latest long-term supported version

## Tools

- [SDK](https://vald.vdaas.org/docs/user-guides/sdks/): Official client libraries
- [Demo](https://github.com/vdaas/vald-demo): Demo repository using sample data

## Vald Users

<p align="center">
    <img src="./assets/image/vald-users/yahoojapan.svg" alt="yahoojapan" width="30%" />
&nbsp; &nbsp; &nbsp; &nbsp;
    <img src="./assets/image/vald-users/japansearch_color.png.webp" alt="jpsearch" />
  </ul>
</p>

## Contribution

Please read the [contribution guide](https://vald.vdaas.org/docs/contributing/contributing-guide).

Before your first commit to this repository, it is strongly recommended to run the commands below.

```shell
git clone https://github.com/vdaas/vald && cd vald
make init
```

## Contributors

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->

[![All Contributors](https://img.shields.io/badge/all_contributors-16-orange.svg?style=flat-square)](#contributors)

<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="http://kpango.com"><img src="https://avatars1.githubusercontent.com/u/9798091?v=4?s=100" width="100px;" alt="Yusuke Kato"/><br /><sub><b>Yusuke Kato</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kpango" title="Code">💻</a> <a href="#design-kpango" title="Design">🎨</a> <a href="#maintenance-kpango" title="Maintenance">🚧</a> <a href="#projectManagement-kpango" title="Project Management">📆</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/rinx"><img src="https://avatars3.githubusercontent.com/u/1588935?v=4?s=100" width="100px;" alt="Rintaro Okamura"/><br /><sub><b>Rintaro Okamura</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=rinx" title="Code">💻</a> <a href="https://github.com/vdaas/vald/commits?author=rinx" title="Documentation">📖</a> <a href="#maintenance-rinx" title="Maintenance">🚧</a> <a href="#platform-rinx" title="Packaging/porting to new platform">📦</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://morimoto.dev/"><img src="https://avatars2.githubusercontent.com/u/413873?v=4?s=100" width="100px;" alt="Kosuke Morimoto"/><br /><sub><b>Kosuke Morimoto</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kmrmt" title="Code">💻</a> <a href="#example-kmrmt" title="Examples">💡</a> <a href="#tool-kmrmt" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=kmrmt" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/vankichi"><img src="https://avatars3.githubusercontent.com/u/13959763?v=4?s=100" width="100px;" alt="Kiichiro YUKAWA"/><br /><sub><b>Kiichiro YUKAWA</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=vankichi" title="Documentation">📖</a> <a href="#maintenance-vankichi" title="Maintenance">🚧</a> <a href="https://github.com/vdaas/vald/commits?author=vankichi" title="Tests">⚠️</a> <a href="#tutorial-vankichi" title="Tutorials">✅</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/datelier"><img src="https://avatars3.githubusercontent.com/u/57349093?v=4?s=100" width="100px;" alt="datelier"/><br /><sub><b>datelier</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=datelier" title="Code">💻</a> <a href="#ideas-datelier" title="Ideas, Planning, & Feedback">🤔</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/kevindiu"><img src="https://avatars1.githubusercontent.com/u/1985382?v=4?s=100" width="100px;" alt="Kevin Diu"/><br /><sub><b>Kevin Diu</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=kevindiu" title="Documentation">📖</a> <a href="#example-kevindiu" title="Examples">💡</a> <a href="https://github.com/vdaas/vald/commits?author=kevindiu" title="Tests">⚠️</a> <a href="#tutorial-kevindiu" title="Tutorials">✅</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://twitter.com/hiroto_hlts2"><img src="https://avatars0.githubusercontent.com/u/25459661?v=4?s=100" width="100px;" alt="Hiroto Funakoshi"/><br /><sub><b>Hiroto Funakoshi</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=hlts2" title="Documentation">📖</a> <a href="#tool-hlts2" title="Tools">🔧</a> <a href="https://github.com/vdaas/vald/commits?author=hlts2" title="Tests">⚠️</a> <a href="#tutorial-hlts2" title="Tutorials">✅</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/taisuou"><img src="https://avatars0.githubusercontent.com/u/21119375?v=4?s=100" width="100px;" alt="taisho"/><br /><sub><b>taisho</b></sub></a><br /><a href="#design-taisuou" title="Design">🎨</a> <a href="https://github.com/vdaas/vald/commits?author=taisuou" title="Documentation">📖</a> <a href="#example-taisuou" title="Examples">💡</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/pgrimaud"><img src="https://avatars1.githubusercontent.com/u/1866496?v=4?s=100" width="100px;" alt="Pierre Grimaud"/><br /><sub><b>Pierre Grimaud</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=pgrimaud" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://omerkatz.com"><img src="https://avatars.githubusercontent.com/u/48936?v=4?s=100" width="100px;" alt="Omer Katz"/><br /><sub><b>Omer Katz</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=thedrow" title="Documentation">📖</a> <a href="#tutorial-thedrow" title="Tutorials">✅</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/zchee"><img src="https://avatars.githubusercontent.com/u/6366270?v=4?s=100" width="100px;" alt="Koichi Shiraishi"/><br /><sub><b>Koichi Shiraishi</b></sub></a><br /><a href="#a11y-zchee" title="Accessibility">️️️️♿️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/liusy182"><img src="https://avatars.githubusercontent.com/u/3293332?v=4?s=100" width="100px;" alt="Siyuan Liu"/><br /><sub><b>Siyuan Liu</b></sub></a><br /><a href="#a11y-liusy182" title="Accessibility">️️️️♿️</a> <a href="#example-liusy182" title="Examples">💡</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/dotdc"><img src="https://avatars.githubusercontent.com/u/12827900?v=4?s=100" width="100px;" alt="David Calvert"/><br /><sub><b>David Calvert</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=dotdc" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/takuyaymd"><img src="https://avatars.githubusercontent.com/u/49614391?v=4?s=100" width="100px;" alt="takuyaymd"/><br /><sub><b>takuyaymd</b></sub></a><br /><a href="https://github.com/vdaas/vald/issues?q=author%3Atakuyaymd" title="Bug reports">🐛</a> <a href="https://github.com/vdaas/vald/commits?author=takuyaymd" title="Code">💻</a> <a href="#maintenance-takuyaymd" title="Maintenance">🚧</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/junsei-ando"><img src="https://avatars.githubusercontent.com/u/1892077?v=4?s=100" width="100px;" alt="junsei-ando"/><br /><sub><b>junsei-ando</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=junsei-ando" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ykadowak"><img src="https://avatars.githubusercontent.com/u/60080334?v=4?s=100" width="100px;" alt="Yusuke Kadowaki"/><br /><sub><b>Yusuke Kadowaki</b></sub></a><br /><a href="https://github.com/vdaas/vald/commits?author=ykadowak" title="Code">💻</a> <a href="https://github.com/vdaas/vald/commits?author=ykadowak" title="Tests">⚠️</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## LICENSE

vald released under Apache 2.0 license, refer [LICENSE](https://github.com/vdaas/vald/blob/main/LICENSE) file.

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B21465%2Fvald.svg?type=large)](https://app.fossa.com/projects/custom%2B21465%2Fvald?ref=badge_large)
