# Vald CI Container

<!-- introduction sentence -->

`vald-ci-container` is designed for running CI workflows on GitHub Actions.

This image includes the basic libraries for running some workflows on the [vdaas/vald](https://github.com/vdaas/vald) repository.

<div align="center">
    <img src="https://github.com/vdaas/vald/blob/main/assets/image/readme.svg" width="50%" />
</div>

[![latest Image](https://img.shields.io/docker/v/vdaas/vald-ci-container/latest?label=vald-ci-container)](https://hub.docker.com/r/vdaas/vald-ci-container/tags?page=1&name=latest)
[![License: Apache 2.0](https://img.shields.io/github/license/vdaas/vald.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![latest ver.](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![Twitter](https://img.shields.io/badge/twitter-follow-blue?logo=twitter&style=flat-square)](https://twitter.com/vdaas_vald)

## Requirement

<!-- FIXME: If image has some requirements, describe here with :warning: emoji -->

<details><summary>linux/amd64</summary><br>

- CPU instruction: requires `AVX2` or `AVX512`

</details>

<details><summary>linux/arm64</summary><br>

- CPU instruction: NOT Apple Silicon

</details>

## Get Started

<!-- Get Started -->

`vald-ci-container` is used for running the workflow on GitHub Actions.

```yaml
name: Name of workflow
on:
  push:
    branches:
      - main
jobs:
  {job_titile}:
    name: {job_name}
    runs-on: ubuntu-latest
    container:
      image: ghcr.io/vdaas/vald/vald-ci-container:latest
    steps:
      - name: {step_name}
    ...
```

The sample workflows are [here](https://github.com/vdaas/vald/search?l=YAML&q=vald-ci-container).

## Versions

| tag     |    linux/amd64     |    linux/arm64     | description                                                                                                                 |
| :------ | :----------------: | :----------------: | :-------------------------------------------------------------------------------------------------------------------------- |
| latest  | :white_check_mark: | :white_check_mark: | the latest image is the same as the latest version of [vdaas/vald](https://github.com/vdaas/vald) repository version.       |
| nightly | :white_check_mark: | :white_check_mark: | the nightly applies the main branch's source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.            |
| vX.Y.Z  | :white_check_mark: | :white_check_mark: | the vX.Y.Z image applies the source code of the [vdaas/vald](https://github.com/vdaas/vald) repository.                     |
| pr-XXX  | :white_check_mark: |        :x:         | the pr-X image applies the source code of the pull request X of the [vdaas/vald](https://github.com/vdaas/vald) repository. |

## Dockerfile

<!-- FIXME -->

The `Dockerfile` of this image is [here](https://github.com/vdaas/vald/blob/main/dockers/ci/base/Dockerfile).

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
