# CHANGELOG v1.4.x

## v1.4.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.4.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.4.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.4.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.4.1</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.4.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.4.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.4.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.4.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.4.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.4.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.4.1/charts/vald-helm-operator/README.md)

### Changes

- [bugfix] fix miss param for fp16 (#1563)
- [bugfix] add missing empty dir mount for s3 backup without pvs (#1562)

## v1.4.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.4.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.4.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.4.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.4.0</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.4.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.4.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.4.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.4.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.4.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.4.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.4.0/charts/vald-helm-operator/README.md)

### Changes

<!--- Describe your changes in detail -->

:sparkles: New features

- Add agent minnum field and support fp16 (#1558)
- Add pyroscope manifest running on persistent volume (#1551)
- Add settings for profiling with pyroscope (#1539)
- Add pyroscope manifest (#1520)
- Add linear search handler for gateway-lb / ingress-filter / agent-core-ngt (#1511)
- Add grpc custom codec (#1490)
- Update grpc codes (#1489)
- Add NGT linear search API (#1504)

:recycle: Refactor

- Delete unsupported library sptag (#1559)
- Improved search operation (#1546)
- Add description timeout of search config (#1541)
- Update dependencies version (#1538)
- Happy new year (#1525)
- Update libs version (#1524)
- Add the missing go.sum (#1517)
- Update license headers / Format codes (#1514)
- Add .gitattributes (#1512)
- Update get started with using kubernetes ingress (#1510)
- Fix command template (#1508)
- Update dependencies version (#1501)
- modify .pb.go & swagger (#1493)
- Add reshape vector proto, remove meta/backup proto (#1492)

:pencil2: Documents

- Add agent component overview document (#1544)
- Add build api proto document (#1540)
- Add remove/object api document (#1536)
- Add search api document (#1534)
- Add upsert api document (#1533)
- Add update api document (#1529)
- Add insert api document (#1516)

:bug: Bugfix

- Fix vulnerability problem of helm-operator (#1535)

:white_check_mark: Test

- Implement internal/net/grpc metrics & proto & types test (#1507)
- Implement and modify internal/io tests (#1509)
- Implement internal/net net.go&dialer.go test (#1505)
- Implement internal/net/grpc/credentials and health test (#1502)
- Implement internal/net/control test (#1500)
- Remove unsupported feature from public document (#1497)

:tada: Cellebration!

- Add liusy182 as a contributor for example (#1519)
- Add zchee as a contributor for a11y (#1513)
