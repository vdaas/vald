# Vald Agent Sidecar

<!-- introduction sentence -->

`vald-agent-sidecar` is the docker image for the vald-agent-sidecar component.

This image saves the index metadata files to external storage like Amazon S3 or Google Cloud Storage.

It has 2 main features:

1. Backup
   - When `Agent` completes creating the index metadata files, `Sidecar` hooks to store them in the external storage.
1. Restore
   - When the Vald Agent Pod restarts, the index structure on the `vald-agent` component is restored from the external backup files.

<!-- FIXME: document URL -->

For more details, please refer to the [component document](https://vald.vdaas.org/docs/overview/component/agent/#sidecar).

<div align="center">
    <img src="https://github.com/vdaas/vald/blob/main/assets/image/readme.svg?raw=true" width="50%" />
</div>

[![latest Image](https://img.shields.io/docker/v/vdaas/vald-agent-sidecar/latest?label=vald-agent-sidecar)](https://hub.docker.com/r/vdaas/vald-agent-sidecar/tags?page=1&name=latest)
[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![latest ver.](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

## Requirement

<!-- FIXME: If image has some requirements, describe here with :warning: emoji -->

### linux/amd64

- Image: `vald-agent-ngt`
- External components: Amazon S3 or Google Cloud Storage

### linux/arm64

- Image: `vald-agent-ngt`
- External components: Amazon S3 or Google Cloud Storage

## Get Started

<!-- Get Started -->
<!-- Vald Agent NGT requires more chapter Agent Standalone -->

`vald-agent-sidecar` is used for one of the components of the Vald cluster, which means it should be used on the Kubernetes cluster, not the local environment or Docker.

Please refer to the [Get Started](https://vald.vdaas.org/docs/tutorial/get-started) for deploying the Vald cluster and [Backup configuration](https://vald.vdaas.org/docs/user-guides/backup-configuration/) to enable the backup feature.

## Versions

| tag     | linux/amd64 | linux/arm64 | description                                                                                                                     |
| :------ | :---------: | :---------: | :------------------------------------------------------------------------------------------------------------------------------ |
| latest  |     ✅      |     ✅      | the latest image is the same as the latest version of [vdaas/vald](https://github.com/vdaas/vald) repository version.           |
| nightly |     ✅      |     ✅      | the nightly applies the main branch's source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                |
| vX.Y.Z  |     ✅      |     ✅      | the vX.Y.Z image applies the source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                         |
| pr-XXX  |     ✅      |     ❌      | the pr-XXX image applies the source code of the pull request XXX of the [vdaas/vald](https://github.com/vdaas/vald) repository. |

## Dockerfile

<!-- FIXME -->

The `Dockerfile` of this image is [here](https://github.com/vdaas/vald/blob/main/dockers/agent/sidecar/Dockerfile).

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
