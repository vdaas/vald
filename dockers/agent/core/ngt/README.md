# Vald Agent NGT

<!-- introduction sentence -->

`vald-agent-ngt` is the docker image for the vald-agent-ngt component.

This image is responsible for the following:

- Store index data along with the user request.
  - The store destination is In-Memory, Volume Mounts, Persistent Volume, or External Storage(⚠).
- Search the nearest neighbor vectors of the request vector and return the search result.

⚠ When you'd like to use the external storage, it requires [vald-agent-sidecar](https://hub.docker.com/r/vdaas/vald-agent-sidecar/tags?page=1&name=latest) on the Kubernetes cluster.

For more details, please refer to the [component document](https://vald.vdaas.org/docs/overview/component/agent).

<div align="center">
    <img src="https://github.com/vdaas/vald/blob/main/assets/image/readme.svg?raw=true" width="50%" />
</div>

[![latest Image](https://img.shields.io/docker/v/vdaas/vald-agent-ngt/latest?label=vald-agent-ngt)](https://hub.docker.com/r/vdaas/vald-agent-ngt/tags?page=1&name=latest)
[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![latest ver.](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

## Requirement

<!-- FIXME: If image has some requirements, describe here with :warning: emoji -->

### linux/amd64

- CPU instruction: requires `AVX2` or `AVX512`

### linux/arm64

⚠ This image does NOT support running on M1/M2 Mac.

## Get Started

<!-- Get Started -->
<!-- Vald Agent NGT requires more chapter Agent Standalone -->

You can use `vald-agent-ngt` in 3 ways.

1. One of the components of the Vald cluster
   - Refer to [Get Started](https://vald.vdaas.org/docs/tutotial/get-started/).
1. Standalone on the Kubernetes cluster
   - Refer to [Vald Agent Standalone on Kubernetes](https://vald.vdaas.org/docs/tutorial/vald-agent-standalone-on-k8s/)
1. Standalone on Docker
   - Refer to [Vald Agent Standalone on Docker](https://vald.vdaas.org/docs/tutorial/vald-agent-standalone-on-docker/)

## Versions

| tag     | linux/amd64 | linux/arm64 | description                                                                                                                     |
| :------ | :---------: | :---------: | :------------------------------------------------------------------------------------------------------------------------------ |
| latest  |     ✅      |     ✅      | the latest image is the same as the latest version of [vdaas/vald](https://github.com/vdaas/vald) repository version.           |
| nightly |     ✅      |     ✅      | the nightly applies the main branch's source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                |
| vX.Y.Z  |     ✅      |     ✅      | the vX.Y.Z image applies the source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                         |
| pr-XXX  |     ✅      |     ❌      | the pr-XXX image applies the source code of the pull request XXX of the [vdaas/vald](https://github.com/vdaas/vald) repository. |

## Dockerfile

<!-- FIXME -->

The `Dockerfile` of this image is [here](https://github.com/vdaas/vald/blob/main/dockers/agent/core/ngt/Dockerfile).

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
