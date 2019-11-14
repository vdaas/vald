<div align="center">
<img src="./assets/image/logo.svg" width="50%">
</div>

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/vdaas/vald.svg?style=flat-square)](https://github.com/vdaas/vald/releases/latest)
[![codecov](https://codecov.io/gh/vdaas/vald/branch/master/graph/badge.svg?token=2CzooNJtUu&style=flat-square)](https://codecov.io/gh/vdaas/vald)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a6e544eee7bc49e08a000bb10ba3deed)](https://www.codacy.com/app/i.can.feel.gravity/vald?utm_source=github.com&utm_medium=referral&utm_content=vdaas/vald&utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdaas/vald)](https://goreportcard.com/report/github.com/vdaas/vald)
[![GolangCI](https://golangci.com/badges/github.com/vdaas/vald.svg?style=flat-square)](https://golangci.com/r/github.com/vdaas/vald)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/vdaas/vald)
[![GoDoc](http://godoc.org/github.com/vdaas/vald?status.svg)](http://godoc.org/github.com/vdaas/vald)

[![Build docker image: agent-ngt](https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20agent-ngt/badge.svg)](https://hub.docker.com/r/vdaas/vald-agent-ngt)
[![Build docker image: backup-manager](https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20backup-manager/badge.svg)](https://hub.docker.com/r/vdaas/vald-manager-backup-mysql)
[![Build docker image: discoverer-k8s](https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20discoverer-k8s/badge.svg)](https://hub.docker.com/r/vdaas/vald-discoverer-k8s)
[![Build docker image: gateway-vald](https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20gateway-vald/badge.svg)](https://hub.docker.com/r/vdaas/vald-gateway)
[![Build docker image: meta-redis](https://github.com/vdaas/vald/workflows/Build%20docker%20image:%20meta-redis/badge.svg)](https://hub.docker.com/r/vdaas/vald-meta-redis)

vald is high scalable distributed high-speed approximate nearest neighbor search engine

## Requirement

kubernetes 1.12~

## Installation

```shell
helm install vdaas/vald
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
make git/config/init
make git/hooks/init
```

## Contribution

1. Fork it ( https://github.com/vdaas/vald/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Author

- [kpango](https://github.com/kpango)
- [kou-m](https://github.com/kou-m)
- [rinx](https://github.com/rinx)

## Contributor


## LICENSE

vald released under Apache 2.0 license, refer [LICENSE](https://github.com/vdaas/vald/blob/master/LICENSE) file.
