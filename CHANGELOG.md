# CHANGELOG

<!-- NEW ENTRY -->

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
