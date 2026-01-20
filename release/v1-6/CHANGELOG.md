# CHANGELOG v1.6.x

## v1.6.3

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.6.3</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.6.3</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.6.3</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.6.3</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.6.3</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.6.3</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.6.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.6.3</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.6.3)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.6.3/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.6.3/charts/vald-helm-operator/README.md)

### Changes

🐛 Bugfix

- Bugfix Circuit Breaker half-open error handling [#1811](https://github.com/vdaas/vald/pull/1811)

📝 Document fix

- Fix dead link [#1807](https://github.com/vdaas/vald/pull/1807)

:arrow_up: Dependencies

- Update go modules and add small test for strings [#1812](https://github.com/vdaas/vald/pull/1812)

## v1.6.2

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.6.2</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.6.2</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.6.2</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.6.2</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.6.2</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.6.2</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.6.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.6.2</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.6.2)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.6.2/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.6.2/charts/vald-helm-operator/README.md)

### Changes

🐛 Bugfix

- Bugfix success handling in the half-open and add flow control [#1805](https://github.com/vdaas/vald/pull/1805)
- Fix string concat buffer overflow [#1806](https://github.com/vdaas/vald/pull/1806)

:white_check_mark: Test

- Implement pkg/agent/handler createAndSaveIndex test case [#1794](https://github.com/vdaas/vald/pull/1794)

:pencil: Document

- Add cluster role document [#1796](https://github.com/vdaas/vald/pull/1796)
- Fix document format [#1804](https://github.com/vdaas/vald/pull/1804)

## v1.6.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.6.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.6.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.6.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.6.1</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.6.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.6.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.6.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.6.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.6.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.6.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.6.1/charts/vald-helm-operator/README.md)

### Changes

🐛 Bugfix

- fix metrics bug [#1800](https://github.com/vdaas/vald/pull/1800)

:white_check_mark: Test

- Add test for attributesFromError method [#1801](https://github.com/vdaas/vald/pull/1801)

## v1.6.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.6.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.6.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.6.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.6.0</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.6.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.6.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.6.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.6.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.6.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.6.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.6.0/charts/vald-helm-operator/README.md)

### Changes

### Description:

:sparkles: New feature

- Introduce opentelemetry [#1778](https://github.com/vdaas/vald/pull/1778)
- Change opencensus tracing to opentelemetry tracing [#1767](https://github.com/vdaas/vald/pull/1767)
- Add circuit breaker implementation [#1738](https://github.com/vdaas/vald/pull/1738)

⬆️ Dependency update

- Deleted deprecated dependency for OTEL [#1795](https://github.com/vdaas/vald/pull/1795)
- Delete stackdriver dependencies [#1761](https://github.com/vdaas/vald/pull/1761)

🐛 Bug fix

- Fix fossa workflow bug [#1787](https://github.com/vdaas/vald/pull/1787)
- Fix failed github workflow [#1745](https://github.com/vdaas/vald/pull/1745)
- Add safe directory configuration [#1748](https://github.com/vdaas/vald/pull/1748)

:recycle: Refactor

- Refactor internal package (net, file, logger) [#1768](https://github.com/vdaas/vald/pull/1768)
- :recycle: Set default image tag as latest [#1766](https://github.com/vdaas/vald/pull/1766)
- Upgrade ubuntu version [#1743](https://github.com/vdaas/vald/pull/1743)
- Deleted vald\_ prefix of dashboard [#1785](https://github.com/vdaas/vald/pull/1785)
- Use gotestfmt instead of tparse [#1788](https://github.com/vdaas/vald/pull/1788)

✅ Test

- Implement agent handler saveIndex test case [#1731](https://github.com/vdaas/vald/pull/1731)
- Ignore gorules test [#1790](https://github.com/vdaas/vald/pull/1790)
- Fix chaos test [#1757](https://github.com/vdaas/vald/pull/1757)
- Implement agent handler createIndex test case [#1710](https://github.com/vdaas/vald/pull/1710)
- Implement agent handler indexInfo test case [#1708](https://github.com/vdaas/vald/pull/1708)

📝 Document

- Update testing guideline & template [#1791](https://github.com/vdaas/vald/pull/1791)
- Add Client API Config document [#1783](https://github.com/vdaas/vald/pull/1783)
- Add backup configuration document [#1754](https://github.com/vdaas/vald/pull/1754)
- Add upgrade document [#1777](https://github.com/vdaas/vald/pull/1777)
- Update document by feedback [#1773](https://github.com/vdaas/vald/pull/1733)
- Add filter config document [#1755](https://github.com/vdaas/vald/pull/1755)
- Add deployment document [#1758](https://github.com/vdaas/vald/pull/1758)
- Update the images of agent page [#1753](https://github.com/vdaas/vald/pull/1753)
- Update config document [#1751](https://github.com/vdaas/vald/pull/1751)

