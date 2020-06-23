## {{ version }}

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:{{ version }}`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:{{ version }}`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:{{ version }}`
gateway | `docker pull vdaas/vald-gateway:{{ version }}`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:{{ version }}`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:{{ version }}`
compressor | `docker pull vdaas/vald-manager-compressor:{{ version }}`
meta Redis | `docker pull vdaas/vald-meta-redis:{{ version }}`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:{{ version }}`
index manager | `docker pull vdaas/vald-manager-index:{{ version }}`
Helm operator | `docker pull vdaas/vald-helm-operator:{{ version }}`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@{{ version }})
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/{{ version }}/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/{{ version }}/charts/vald-helm-operator/README.md)

### Changes
