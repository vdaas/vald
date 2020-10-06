## {{ version }}

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:{{ version }}</code>
    </td>
  </tr>
</table>

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@{{ version }})
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/{{ version }}/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/{{ version }}/charts/vald-helm-operator/README.md)

### Changes
