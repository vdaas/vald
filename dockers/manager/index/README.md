# Vald Manager Index

<!-- introduction sentence -->

`vald-manager-index` is the image for vald-index-manager component.

`vald-manager-index` has the unique feature to control the index timing for all Vald Agent pods in the Vald cluster.

The main features are:

- Syncing data from Vald Discoverer
- Controlling indexing process for each Vald Agent pods

For more details, please refer to the [component document](https://vald.vdaas.org/docs/overview/component/index-manager).

<div align="center">
    <img src="https://github.com/vdaas/vald/blob/main/assets/image/readme.svg" width="50%" />
</div>

[![latest Image](https://img.shields.io/docker/v/vdaas/vald-agent-ngt/latest?label=vald-agent-ngt)](https://hub.docker.com/r/vdaas/vald-agent-ngt/tags?page=1&name=latest)
[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![latest ver.](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

## Requirement

<!-- FIXME: If image has some requirements, describe here with :warning: emoji -->

<details><summary>linux/amd64</summary><br>

- Image: vald-discoverer-k8s, vald-agent-ngt

</details>

<details><summary>linux/arm64</summary><br>

- Image: vald-discoverer-k8s, vald-agent-ngt

</details>

## Get Started

<!-- Get Started -->
<!-- Vald Agent NGT requires more chapter Agent Standalone -->

`vald-manager-index` is used for one of the components of the Vald cluster, which means it should be used on the Kubernetes cluster, not the local environment or Docker.

Please refer to the [Get Started](https://vald.vdaas.org/docs/tutorial/get-started) for deploy Vald cluster.

## Versions

| tag     |    linux/amd64     |    linux/arm64     | description                                                                                                                 |
| :------ | :----------------: | :----------------: | :-------------------------------------------------------------------------------------------------------------------------- |
| latest  | :white_check_mark: | :white_check_mark: | the latest image is the same as the latest version of [vdaas/vald](https://github.com/vdaas/vald) repository version.       |
| nightly | :white_check_mark: | :white_check_mark: | the nightly applies the main branch's source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.            |
| vX.Y.Z  | :white_check_mark: | :white_check_mark: | the vX.Y.Z image applies the source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                     |
| pr-XXX  | :white_check_mark: |        :x:         | the pr-X image applies the source code of the pull request X of the [vdaas/vald](https://github.com/vdaas/vald) repository. |

## Dockerfile

<!-- FIXME -->

The `Dockerfile` of this image is [here](https://github.com/vdaas/vald/blob/main/dockers/manager/index/Dockerfile).

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
