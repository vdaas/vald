# CHANGELOG v1.3.x

## v1.3.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.3.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.3.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.3.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.3.1</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.3.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.3.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.3.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.3.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.3.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.3.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.3.1/charts/vald-helm-operator/README.md)

### Changes

- Add documentation comments to proto files (#1452)
- Add ngt index count panel to agent grafana dashboard (#1483)
- add grpc keepalive EnforcementPolicy support (#1487)
- fix timing of removeInvalidIndex (#1481)
- [bugfix] add validation to agent service option InitialDelayMaxDuration (#1482)

## v1.3.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.3.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.3.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.3.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.3.0</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.3.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.3.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.3.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.3.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.3.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.3.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.3.0/charts/vald-helm-operator/README.md)

### Changes

- add startupProbe support (#1473)
- add label / field selectors for discoverer (#1472)

