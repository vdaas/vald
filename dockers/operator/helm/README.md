# Vald Helm Operator

<!-- introduction sentence -->

`vald-helm-operator` is the Helm operator image for operating the Vald cluster.

This image is responsible for the following:

- monitoring 2 custom resource definitions, [`ValdRelease`](https://github.com/vdaas/vald/blob/main/charts/vald-helm-operator/crds/valdrelease.yaml) and [`ValdHelmOperatorRelease`](https://github.com/vdaas/vald/blob/main/charts/vald-helm-operator/crds/valdhelmoperatorrelease.yaml)
- Operating Vald Helm Operator and Vald cluster according to applied CRD.

Vald recommends using Vald Helm Operator for production usage.

<div align="center">
    <img src="https://github.com/vdaas/vald/blob/main/assets/image/readme.svg?raw=true" width="50%" />
</div>

[![latest Image](https://img.shields.io/docker/v/vdaas/vald-helm-operator/latest?label=vald-helm-operator)](https://hub.docker.com/r/vdaas/vald-helm-operator/tags?page=1&name=latest)
[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![latest ver.](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

## Requirement

<!-- FIXME: If image has some requirements, describe here with :warning: emoji -->

### linux/amd64

- Others: apply CRDs

### linux/arm64

- Others: apply CRDs

## Get Started

<!-- Get Started -->
<!-- Vald Agent NGT requires more chapter Agent Standalone -->

`vald-helm-operator` is used for monitoring, managing and operating the Vald cluster.
Please refer to the [Deployment](https://vald.vdaas.org/docs/user-guides/deployment/#using-with-vald-helm-operator) document for deployment.

## Versions

| tag     | linux/amd64 | linux/arm64 | description                                                                                                                     |
| :------ | :---------: | :---------: | :------------------------------------------------------------------------------------------------------------------------------ |
| latest  |     ✅      |     ✅      | the latest image is the same as the latest version of [vdaas/vald](https://github.com/vdaas/vald) repository version.           |
| nightly |     ✅      |     ✅      | the nightly applies the main branch's source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                |
| vX.Y.Z  |     ✅      |     ✅      | the vX.Y.Z image applies the source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                         |
| pr-XXX  |     ✅      |     ❌      | the pr-XXX image applies the source code of the pull request XXX of the [vdaas/vald](https://github.com/vdaas/vald) repository. |

## Dockerfile

<!-- FIXME -->

The `Dockerfile` of this image is [here](https://github.com/vdaas/vald/blob/main/dockers/operator/helm/Dockerfile).

## About Vald Project

<!-- About Vald Project -->
<!-- This chapter is static -->

The information about the Vald project, please refer to the following:

- [Official website](https://vald.vdaas.org)
- [GitHub](https://github.com/vdaas/vald)

## Contacts

We're love to support you!
Please feel free to contact us anytime with your questions or issue reports.

- [Official Slack WS](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA)
- [GitHub Issue](https://github.com/vdaas/vald/issues)

## License

This product is under the terms of the Apache License v2.0; refer [LICENSE](https://github.com/vdaas/vald/blob/main/LICENSE) file.
