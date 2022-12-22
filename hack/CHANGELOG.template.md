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
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:{{ version }}</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:{{ version }}</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:{{ version }}</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:{{ version }}</code>
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
