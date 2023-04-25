# Image Title

<!-- introduction sentence -->

`image-name` is the XXX for vald-XXX-YYY.

The responsibility of this image is XXX.

<!-- FIXME: document URL -->

For more details, please refer to the [component document](https://vald.vdaas.org/docs/overview/component).

<div align="center">
    <img src="https://github.com/vdaas/vald/blob/main/assets/image/readme.svg" width="50%" />
</div>

[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![latest ver.](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

## Vald Project

<!-- About Vald Project -->
<!-- This chapter is static -->

Vald is a highly scalable distributed search engine.

Vald is designed and implemented based on the Cloud-Native architecture.

Vald has automatic vector indexing, index backup, horizontal scalability, and more features for handling billion-scale vectors.

### Client Libraries

Vald provides the official client libraries:

- [Go](https://github.com/vdaas/vald-client-go)
- [Java](https://github.com/vdaas/vald-client-java)
- [Python](https://github.com/vdaas/vald-client-python)
- [Node.js](https://github.com/vdaas/vald-client-node)
- [Clojure](https://github.com/vdaas/vald-client-clj)

Also, everyone can build a client library with any programming language by using API proto.
Please refer to the [Building gRPC proto](https://vald.vdaas.org/docs/api/build_proto/).

### Links

For more information about the Vald project, please refer to the following:

- [Official website](https://vald.vdaas.org)
- [GitHub](https://github.com/vdaas/vald)

### Contacts

We're love to support you!
Please feel free to contact us anytime with your questions or issue reports.

- [Official Slack WS](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA)
- [GitHub Issue](https://github.com/vdaas/vald/issues)

## Get Started

<!-- Get Started -->
<!-- Vald Agent NGT requires more chapter Agent Standalone -->

`image-name` is used for one of the components of the Vald cluster, which means it should be used on the Kubernetes cluster, not the local environment or Docker.

Please refer to the [Get Started](https://vald.vdaas.org/docs/get-started) for deploy Vald cluster.

## Versions

| tag     |       x86_64       |        Arm         | description                                                                                                           |
| :------ | :----------------: | :----------------: | :-------------------------------------------------------------------------------------------------------------------- |
| latest  | :white_check_mark: | :white_check_mark: | the latest image is the same as the latest version of [vdaas/vald](https://github.com/vdaas/vald) repository version. |
| nightly | :white_check_mark: |        :x:         | the nightly applies the main branch's source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.      |
| vX.Y.Z  | :white_check_mark: | :white_check_mark: | the vX.Y.Z image applies the source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.               |

## Dockerfile

<!-- FIXME -->

The `Dockerfile` of this image is [here](https://github.com/vdaas/vald/blob/main/dockers/agent/core/ngt/Dockerfile).

## License

This product is under the terms of the Apache License v2.0; refer [LICENSE](https://github.com/vdaas/vald/blob/main/LICENSE) file.

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
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/takuyaymd"><img src="https://avatars.githubusercontent.com/u/49614391?v=4?s=100" width="100px;" alt="takuyaymd"/><br /><sub><b>takuyaymd</b></sub></a><br /><a href="https://github.com/vdaas/vald/issues?q=author%3Atakuyaymd" title="Bug reports">🐛</a> <a href="https://github.com/vdaas/vald/commits?author=takuyaymd" title="Code">💻</a></td>
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
