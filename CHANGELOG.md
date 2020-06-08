# CHANGELOG

<!-- NEW ENTRY -->

## v0.0.38

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.38`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.38`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.38`
gateway | `docker pull vdaas/vald-gateway:v0.0.38`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.38`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.38`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.38`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.38`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.38`
index manager | `docker pull vdaas/vald-manager-index:v0.0.38`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.38`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.38)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.38/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.38/charts/vald-helm-operator/README.md)

### Changes
- send PR when K8s manifests are updated (#435)
- Implementation of agent-sidecar storage backup logic (#409)
- Fix structure of grpc java package (#431)
- Remove AUTHORS/CONTRIBUTORS file (#428)
- Contribute document (#390)
- Separate the component section from architecture doc (#430)
- fix: delete fch channel because fch causes channel blocking (#429)
- Update Operator SDK version (#412)
- Documentation: About vald (#374)
- Vald architecture document (#366)
- Add JSON schema for Vald Helm Chart (#365)
- Revise ChatOps not to add go.mod & go.sum when /format runs (#406)
- add agent sidecar flame for implementation (#404)
- Upgrade Operator SDK version / Remove useless GO111MODULE=off (#402)
- Upgrade tools version (#399)
- Add app.kubernetes.io/xxx labels to all resources (#397)
- Vald contacts document (#373)
- add trace spans and metrics for agent-ngt and index-manager (#389)
- Add gen-test command for chatops (#379)
- Add internal/db/storage/blob (#388)

## v0.0.37

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.37`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.37`
gateway | `docker pull vdaas/vald-gateway:v0.0.37`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.37`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.37`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.37`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.37`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.37`
index manager | `docker pull vdaas/vald-manager-index:v0.0.37`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.37`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.37)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.37/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.37/charts/vald-helm-operator/README.md)

### Changes
- add agent auto save indexing feature (#385)
- :bug: fix ngt `distance_type` (#384)
- Add topology spread constraints (#383)

## v0.0.36

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.36`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.36`
gateway | `docker pull vdaas/vald-gateway:v0.0.36`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.36`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.36`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.36`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.36`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.36`
index manager | `docker pull vdaas/vald-manager-index:v0.0.36`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.36`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.36)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.36/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.36/charts/vald-helm-operator/README.md)

### Changes
- update dependencies version (#381)
- Fix missing value on compressor health servers (#377)
- Fix compressor readiness shutdown_duration / Fix cassandra … (#376)
- Bump gopkg.in/yaml.v2 from 2.2.8 to 2.3.0 (#375)
- Fix`internal/log/format` to match the test template (#369)
- Fix `internal/log/logger` to match the test template (#371)
- Fix failing tests of `internal/log` and modified to match the test template  (#368)
- Add enabled flag to each component in Helm chart (#372)
- Add configurations.md (#356)


## v0.0.35

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.35`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.35`
gateway | `docker pull vdaas/vald-gateway:v0.0.35`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.35`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.35`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.35`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.35`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.35`
index manager | `docker pull vdaas/vald-manager-index:v0.0.35`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.35`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.35)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.35/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.35/charts/vald-helm-operator/README.md)

### Changes
- add storage backup option to agent (#367)
- Add client-node dispatcher (#370)
- Bump github.com/tensorflow/tensorflow (#364)
- change fmt.Errorf to errors.Errorf (#361)
- add goleak (#359)


## v0.0.34

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.34`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.34`
gateway | `docker pull vdaas/vald-gateway:v0.0.34`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.34`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.34`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.34`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.34`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.34`
index manager | `docker pull vdaas/vald-manager-index:v0.0.34`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.34`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.34)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.34/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.34/charts/vald-helm-operator/README.md)

### Changes
- feature/internal/cassandra/add option (#358)
- update helm docs when version is published (#355)
- upgrade tools (#354)
- bugfix protoc-gen-validate resolve failure (#353)
- Fix conflicts between formatter and helm template (#350)
- Add more options and remove valdhelmoperatorrelease, valdrelease from vald-helm-operator chart (#334)


## v0.0.33

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.33`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.33`
gateway | `docker pull vdaas/vald-gateway:v0.0.33`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.33`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.33`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.33`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.33`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.33`
index manager | `docker pull vdaas/vald-manager-index:v0.0.33`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.33`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.33)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.33/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.33/charts/vald-helm-operator/README.md)

### Changes
- update k8s dependencies (#349)
- create missing test files by the our original test template (#348)
- create test template for using gotests (#327)
- Revise coverage CI settings (#347)
- fix tensorflow.go, option.go (#261)


## v0.0.32

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.32`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.32`
gateway | `docker pull vdaas/vald-gateway:v0.0.32`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.32`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.32`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.32`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.32`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.32`
index manager | `docker pull vdaas/vald-manager-index:v0.0.32`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.32`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.32)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.32/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.32/charts/vald-helm-operator/README.md)

### Changes
- bugfix ip discoverer disconnection too slow (#344)
- Compressor: backup vectors in queue using PostStop function (#345)
- Revise backup/meta Cassandra default values (#336)


## v0.0.31

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.31`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.31`
gateway | `docker pull vdaas/vald-gateway:v0.0.31`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.31`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.31`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.31`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.31`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.31`
index manager | `docker pull vdaas/vald-manager-index:v0.0.31`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.31`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.31)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.31/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.31/charts/vald-helm-operator/README.md)

### Changes
- Resolve busy-loop on worker (#339)


## v0.0.30

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.30`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.30`
gateway | `docker pull vdaas/vald-gateway:v0.0.30`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.30`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.30`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.30`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.30`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.30`
index manager | `docker pull vdaas/vald-manager-index:v0.0.30`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.30`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.30)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.30/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.30/charts/vald-helm-operator/README.md)

### Changes
- async compressor
- optimized gRPC pool connection
- update helm chart API version
- internal gRPC client for Vald
- Cassandra NewConvictionPolicy
- dicoverer now returns clone object
- new internal/singleflight package
- new internal/net package
- coding guideline


## v0.0.26

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.26`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.26`
gateway | `docker pull vdaas/vald-gateway:v0.0.26`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.26`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.26`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.26`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.26`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.26`
index manager | `docker pull vdaas/vald-manager-index:v0.0.26`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.26`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.26)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.26/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.26/charts/vald-helm-operator/README.md)

### Changes
- added helm operator
- added telepresence
- improved meta-Cassandra performance
- added doc users/get-started
- fixed some bugs
