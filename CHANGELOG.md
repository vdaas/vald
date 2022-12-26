# CHANGELOG

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

## v1.5.6

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.6</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.6</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.6</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.6</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.6</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.6</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.6</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.6)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.6/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.6/charts/vald-helm-operator/README.md)

### Changes

✨ Feature

- Add distance type [#1732](https://github.com/vdaas/vald/pull/1732)
- Add vald-helm-operator E2E [#1722](https://github.com/vdaas/vald/pull/1722)

⬆️ update dependencies

- Delete tensorflow [#1723](https://github.com/vdaas/vald/pull/1723)
- Update deps [#1724](https://github.com/vdaas/vald/pull/1724)

✏️ Documents

- Add capacity planning document [#1714](https://github.com/vdaas/vald/pull/1714)
- Update filter gateway document [#1721](https://github.com/vdaas/vald/pull/1721)
- Fix capacity planning doc [#1736](https://github.com/vdaas/vald/pull/1736)
- Fix file name of capacity planning document [#1737](https://github.com/vdaas/vald/pull/1737)

## v1.5.5

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.5</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.5</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.5</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.5</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.5</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.5</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.5</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.5)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.5/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.5/charts/vald-helm-operator/README.md)

### Changes

♻️ Refactor

- Make format [#1717](https://github.com/vdaas/vald/pull/1717)
- Remove unnecessary import path on pkg test [#1712](https://github.com/vdaas/vald/pull/1712)
- Fix to get only backoff metrics of discover RPC [#1706](https://github.com/vdaas/vald/pull/1706)

🐛 Bugfix

- Return uuid when exits rpc called [#1709](https://github.com/vdaas/vald/pull/1709)

✏️ Documents

- Add textlint for document [#1715](https://github.com/vdaas/vald/pull/1715)
- Modified design of troubleshooting image [#1705](https://github.com/vdaas/vald/pull/1705)
- Update tutorial images [#1704](https://github.com/vdaas/vald/pull/1704)
- Add troubleshooting flow chart document [#1688](https://github.com/vdaas/vald/pull/1688)
- Update data-flow docs for new images and using remove instead of delete [#1693](https://github.com/vdaas/vald/pull/1693)

:white_check_mark: Test

- Implement agent handler getObject test case [#1707](https://github.com/vdaas/vald/pull/1707)
- Implement stream insert test case [#1697](https://github.com/vdaas/vald/pull/1697)
- Implement upsert test cases [#1685](https://github.com/vdaas/vald/pull/1685)

⬆️ Update dependencies

- Update deps [#1719](https://github.com/vdaas/vald/pull/1719)
- Update deps [#1702](https://github.com/vdaas/vald/pull/1702)
- Automatically update k8s manifests [#1701](https://github.com/vdaas/vald/pull/1701)

## v1.5.4

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.4</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.4</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.4</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.4</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.4</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.4</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.4</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.4)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.4/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.4/charts/vald-helm-operator/README.md)

### Changes

✨ New features

- Add backoff metrics panel [#1694](https://github.com/vdaas/vald/pull/1694)
- Add backoff metrics for grpc [#1684](https://github.com/vdaas/vald/pull/1684)
- Setup jaeger operator [#1682](https://github.com/vdaas/vald/pull/1682)

⬆️ update dependencies

- Update deps [#1695](https://github.com/vdaas/vald/pull/1695)
- Update deps [#1699](https://github.com/vdaas/vald/pull/1699)

♻️ Refactor

- Split agent pkg handler implementation [#1690](https://github.com/vdaas/vald/pull/1690)
- Refactor pkg test helper functions [#1678](https://github.com/vdaas/vald/pull/1678)

🐛 Bugfix

- Fix error handling in readyForUpdate and return NotFound error when delete fails in multiUpdate [#1681](https://github.com/vdaas/vald/pull/1681)
- Fix race error of server package [#1689](https://github.com/vdaas/vald/pull/1689)

✏️ Documents

- Add API status code description [#1679](https://github.com/vdaas/vald/pull/1679)
- Modified data flow images [#1687](https://github.com/vdaas/vald/pull/1687)
- Correspond to update omission [#1686](https://github.com/vdaas/vald/pull/1686)
- Renew basic architecture image [#1680](https://github.com/vdaas/vald/pull/1680)

## v1.5.3

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.3</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.3</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.3</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.3</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.3</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.3</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.3</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.3)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.3/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.3/charts/vald-helm-operator/README.md)

### Changes

:sparkles: Feature

- add invalid id argument handling [#1667](https://github.com/vdaas/vald/pull/1667)

:bug: Bugfix

- fix update bug [#1660](https://github.com/vdaas/vald/pull/1660)
- fix typo argument in Makefile and Makefile.d/function.mk [#1673](https://github.com/vdaas/vald/pull/1673)

:green_heart: CI

- fix chaos test temporarily [#1665](https://github.com/vdaas/vald/pull/1665)

:memo: Document

- add search config details document [#1661](https://github.com/vdaas/vald/pull/1661)

:white_check_mark: Test

- implement pkg handler exists test cases [#1628](https://github.com/vdaas/vald/pull/1628)
- implement multi insert test case for pkg agent handler [#1612](https://github.com/vdaas/vald/pull/1612)
- create investigation test of max dim for NGT [#1633](https://github.com/vdaas/vald/pull/1633)
- implement pkg handler remove test cases [#1644](https://github.com/vdaas/vald/pull/1644)
- add e2e test for maxDimensionTest [#1650](https://github.com/vdaas/vald/pull/1650)
- implement update handler test cases [#1657](https://github.com/vdaas/vald/pull/1657)

:arrow_up: Update dependencies

- update manifests version [#1642](https://github.com/vdaas/vald/pull/1642)
- update go module [#1643](https://github.com/vdaas/vald/pull/1643)
- fix go tool installation [#1649](https://github.com/vdaas/vald/pull/1649)
- update kind version [#1668](https://github.com/vdaas/vald/pull/1668)

:art: Design

- update dataflow image [#1647](https://github.com/vdaas/vald/pull/1647)

:lock: Security

- fix CWE-285 [#1654](https://github.com/vdaas/vald/pull/1654)

## v1.5.2

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.2</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.2</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.2</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.2</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.2</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.2</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.2</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.2)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.2/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.2/charts/vald-helm-operator/README.md)

### Changes

:arrow_up: update dependencies

- update libs version (#1636)
- update libs version (#1632)

:sparkles: feature

- use ReadWriteOncePod instead of ReadWriteOnce and remove initializer (#1627)

:recycle: refactor

- refactor: improve memory high usage of vald-agent (#1617)

:lock: security

- security: fix vulnerability problem of helm operator (#1625)
- [Security] Fix vulnerability problem of helm operator (#1611)
- security fix Vulnerability due to usage of old golang.org/x:net in example depentency (#1641)

:green_heart: ci

- ci: Fix CodeQL warning (#1629)
- fix fails actions job & update version (#1620)
- [CI] Allow e2e deploy action jobs to run in parallel (#1616)

:white_check_mark: test

- fix superfluous response.WriteHeader call (#1631)
- implement search-by-id pkg test (#1624)

:memo: document

- document: update formats (#1634)
- docs: add dotdc as a contributor for doc (#1623)
- doc: fixed architecture link in get-started.md (#1619)
- add FAQ and Troubleshooting document (#1591)

## v1.5.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.1</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.1/charts/vald-helm-operator/README.md)

### Changes

:arrow_up: update dependencies

- Upgrade pyroscope version (#1605)

:bug: bugfix

- bugfix internal/file and add CreateTemp function and resolve go module failure (#1608)

:white_check_mark: Test

- fix search handler test (#1613)
- fix fails test of e2e chaos (#1603)
- Implement pkg ngt handler insert test (#1552)
- implement agent ngt handler search test (#1557)

:memo: document

- add vald users (#1601)
- update brand guidelines pdf (#1600)
- fix file name of search api (#1599)
- [ImgBot] Optimize images (#1598)
- cleanup document images (#1595)

## v1.5.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.5.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.5.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.5.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.5.0</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.5.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.5.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.5.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.5.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.5.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.5.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.5.0/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New features

- Add values yaml for the back up the agent index data (#1566)
- Implement uniform/gaussian distributed random float32/uint8 vector generator (#1573)
- Add min_num in each search service api (#1576)
- Add copy on write (#1578)
- Add example values for using Pyroscope (#1582)

:recycle: Refactor

- Deleted resource limits of agent ngt and added in memory mode example (#1571)
- Improve string conversion performace (#1577)
- Update dependencies version (#1593)

:pencil2: Documents

- Add gateway component overview document (#1549)
- Update API docs and fix format (#1568)
- Update dataflow images (#1572)
- Add discoverer component overview document (#1574)
- Add index manager component overview document (#1575)
- Update README (#1584)

:bug: Bugfix

- Add error handling when there is no data in google cloud storage. (#1556)
- Reviewdog markdown workflow (#1585)
- Fix invalid URL (#1589)

:white_check_mark: Test

- Implement test for net/grpc codec,logger and server (#1530)
- Update test for go1.17 update & -race test (#1431)

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

## v1.2.4

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.2.4</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.2.4</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.2.4</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.2.4</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.2.4</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.2.4</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.2.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.2.4</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.2.4)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.2.4/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.2.4/charts/vald-helm-operator/README.md)

### Changes

- update go patch version (#1464)
- [bugfix] sidecar e2e testing (#1465)
- update dependencies version including NGT (#1461)
- fix unlimited gorountine processing in kvsdb (#1458)
- Refactor hack pkg agent e2e benchmark (#1430)
- Remove non-exist components from doc and ci (#1450)
- :recycle: :pencil: add default_pool_size in example yml (#1457)
- reduce memory usage around ngt implementation & refactor agent/lb & auto-generate unit test (#1449)
- Remove rinx from several yamls (#1451)
- Add E2E scenario with SkipStrictExistCheck enabled (#1415)
- fix filter-gateway chart (#1454)
- Implement pkg/agent/core/ngt/service/option test (#1429)

## v1.2.3

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.2.3</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.2.3</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.2.3</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.2.3</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.2.3</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.2.3</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.2.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.2.3</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.2.3)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.2.3/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.2.3/charts/vald-helm-operator/README.md)

### Changes

- update go module dependencies & update go patch version to 1.17.1 (#1445)
- [bugfix] fix unknown status problems in search operation (#1439)
- Add tools.mk for installing CI tools (#1442)
- Fix icon image path in Chart.yaml (#1441)
- fix(Makefile): :bug: Remove components (#1440)
- Change golangci-lint version to the latest and Fix docker build error (#1435)
- Fix golangci-lint rule (#1380)
- Add vqueue test for Exists and GetVector method (#1425)

## v1.2.2

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.2.2</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.2.2</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.2.2</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.2.2</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.2.2</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.2.2</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.2.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.2.2</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.2.2)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.2.2/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.2.2/charts/vald-helm-operator/README.md)

### Changes

- Add stateful PBT for NGT service (#1384)
- add default logger when logger not initialized (#1424)
- feat: :sparkles: Add LB gateway dashboard (#1420)
- chore-deps: :arrow_up: Upgrade OSDK to v1.11.0 (#1422)
- change default epsilon value to 0.1 from 0.0.1 (#1421)
- add go vet for checking cpu compatibility and update deps and refactor small code (#1418)
- Delete insert and delete channel of vqueue (#1400)
- downgrade cloud.google.com/go to resolve runtime panic (#1413)
- [bugfix] resolve errgroup limitation channel close panic (#1412)
- Update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE (#1409)
- update go version to 1.17 and update go module dependencies (#1404)
- refactor grpc logging and use os.PathSeparator instead of / (#1405)
- Remove schema of security contexts (#1406)
- change security context uid to distroless nonroot uid (#1402)

## v1.2.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.2.1</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.2.1</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.2.1</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.2.1</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.2.1</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.2.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.2.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.2.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.2.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.2.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.2.1/charts/vald-helm-operator/README.md)

### Changes

- remove global errgroup from kvsdb goroutine limitation (#1398)
- add google grpc healthz & logger (#1394)
- fix e2e multi apis test code (#1397)
- :pencil: fix command in vald agent standalone on docker (#1393)
- add Makefile & Makefile.d to build kick rule (#1391)
- format codes (#1389)
- [proposal] add golines as vald default formatter (#1337)
- replace gogo protobuf to vt protobuf (#1378)
- refactor gateway request result location aggregation logic & use xxh3 for kvsdb hashing (#1376)
- Add E2E tests for multi-APIs (#1353)
- Add test for the function to get length of vqueue (#1382)
- remove unused path of format/yaml command (#1383)
- [Documentation] fix image filename (#1377)
- Vald moves to a vector search engine that enables more simple and high-speed retrieval. This is the first step of Simple-Vald. (#1365)
- Add e2e-profiling job (#1356)
- fix svg error on dataflow image (#1375)
- remove invalid gateway component reference from chart.NOTES (#1371)

## v1.2.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.2.0</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.2.0</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.2.0</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.2.0</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.2.0</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.2.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.2.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.2.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.2.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.2.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.2.0/charts/vald-helm-operator/README.md)

### Changes

# Release Vald v1.2.0

## Changes:

- fix build failure on arm cpu due to the xxh3 dependency version (#1368)
- Add core algorithm ngt test (#1364)
- add mote accurate directory detection logic for removing (#1366)
- remove meta backup compressor from charts (#1334)
- update go module version (#1363)
- Update architecture overview (#1304)
- Improve internal package tests (#1227)
- fix bug of bulkinsert when error occurs and refactor error message (#1361)
- change initial-index-directory removal logic (#1359)
- ci: :construction_worker: Add gnupg to ci-container (#1362)
- ci: :construction_worker: add condition to trigger importing gpg key (#1360)
- fix: :bug: Fix mount paths when using persistent volume claim template (#1358)
- ci: :construction_worker: Add Upsert operation tests (#1347)
- config: :gear: Use GPG key for signing commits (#1351)
- add time validation for vqueue (#1352)
- chore-deps: :arrow_up: Upgrade tools (#1355)
- comment out backup/meta/compressor build command in Makefile (#1346)
- remove: :heavy_minus_sign: Remove backup, meta components from CI (#1331)
- add line trace logging when log mode is glg and level is debug (#1348)
- update go version to 1.16.6 and update go module dependencies (#1345)
- [bugfix] agent createindex operation's time.Ticker purges too slow & buffer overflow due to the unnecessary error wrapping (#1343)
- [bugfix] change kvsdb and vqueue check order for Exists operation (#1341)
- update go module dependencies (#1336)
- bugfix nil pointer panic in agent's MultiUpsert operation (#1335)
- [bugfix] agent.GetObject API returns old indexed vector problem instead of vqueue's new data (#1333)
- add timestamp handler for agent timestamp controlled update (#1324)
- fix: :bug: Fix typo (#1330)
- remove owner and description info from resource info rich error for each grpc handler (#1327)
- Update documents: configurations (#1289

## v1.1.2

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.1.2</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.1.2</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.1.2</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.1.2</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.1.2</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.1.2</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.1.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.1.2</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.1.2)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.1.2/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.1.2/charts/vald-helm-operator/README.md)

### Changes

- bugfix remove api return wrong error of non exsiting replicas (#1318)
- update go modules and update go version to 1.16.5 from 1.16.4 (#1306)
- Implement test for internal/core/alogrithm/ngt/option and refactor (#1251)
- fix build error due to e2e benchmark invalid creating new client (#1221)
- bugfix correct error backoff handling for lb gateway (#1309)
- Apply ruleguard fixes (#1302)
- Add ruleguard rules (#1301)
- bugfix add feature vector duplication checking for LB-GW Upsert operation. (#1303)
- fix fails race condition test (#1299)
- Update dashboards, add operator dashboard and several panels (#1280)
- Add composite actions for E2E tests (#1257)
- add save index wait duration to index-manager for multiple agent index upload delay (#1292)
- Bugfix index count problem of ngt (#1288)
- fix example/helm/values.yaml (#1290)
- Fix time condition of agent & lb gateway (#1285)
- Use gofumpt for format workflow (#1281)
- [bugfix] kvsdb: do not increment the counter if key exists (#1282)
- Revert ":art: Use gofumpt as a default formatter (#1278)" (#1279)
- :art: Use gofumpt as a default formatter (#1278)
- :pencil: fix some wrongs (#1274)
- Refactor E2E tests: split operations into a new package (#1220)
- :bug: Fix the chart: agent-sidecar initContainer mode (#1271)

## v1.1.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.1.1</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.1.1</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.1.1</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.1.1</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.1.1</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.1.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.1.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.1.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.1.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.1.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.1.1/charts/vald-helm-operator/README.md)

### Changes

- [patch] release Vald v1.1.1 ([#1268](https://github.com/vdaas/vald/pull/1268))
- fix copy CI error ([#1269](https://github.com/vdaas/vald/pull/1269))
- bugfix nil internal/net/grpc.Client.Reconnect nil pointer fix ([#1270](https://github.com/vdaas/vald/pull/1270))
- Docs/readme/replace architecture overview ([#1267](https://github.com/vdaas/vald/pull/1267))
- format code ([#1266](https://github.com/vdaas/vald/pull/1266))
- :robot: Update license headers / Format Go codes and YAML files ([#1265](https://github.com/vdaas/vald/pull/1265))
- [bugfix] append timing of vqueue existing map ([#1264](https://github.com/vdaas/vald/pull/1264))
- add resource type infos to grpc error response ([#1262](https://github.com/vdaas/vald/pull/1262))
- :wrench: Fix clusterrole of scylla-cluster member ([#1263](https://github.com/vdaas/vald/pull/1263))
- refactor file I/O & replace io.Copy to vald original faster Copy function ([#1261](https://github.com/vdaas/vald/pull/1261))
- [bugfix] change exists cheking for agent vqueue ([#1256](https://github.com/vdaas/vald/pull/1256))
- :robot: Update license headers / Format Go codes and YAML files ([#1255](https://github.com/vdaas/vald/pull/1255))
- :fire: Remove invalid initialization option ([#1252](https://github.com/vdaas/vald/pull/1252))
- add single connection client for agent & vald ([#1254](https://github.com/vdaas/vald/pull/1254))
- bugfix remove unneccessary error return & add grpc status code handling for backoff ([#1253](https://github.com/vdaas/vald/pull/1253))
- bugfix agent vqueue & refactor tools/deps ([#1250](https://github.com/vdaas/vald/pull/1250))
- Add test for pkg/agent/core/ngt/service/vqueue/option ([#1233](https://github.com/vdaas/vald/pull/1233))
- Add config agent core ngt service kvs test ([#1223](https://github.com/vdaas/vald/pull/1223))
- Fix update-helm-chart workflow ([#1249](https://github.com/vdaas/vald/pull/1249))
- :robot: Automatically update k8s manifests ([#1248](https://github.com/vdaas/vald/pull/1248))
- Fix path to the yq binary ([#1247](https://github.com/vdaas/vald/pull/1247))

## v1.1.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.1.0</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.1.0</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.1.0</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.1.0</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.1.0</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.1.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.1.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.1.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.1.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.1.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.1.0/charts/vald-helm-operator/README.md)

### Changes

#### Feature

- add correct grpc status and error handling (#1224)
- Add Cloud Storage mode of Agent Sidecar. (#519)
- Add general info metrics package / Add NGT info metrics (#1195)
- improve internal/info implementation (#1177)
- Add saving flag to index.info payload (#1200)
- Fix mysql and cassandra schema definitions in initialize jobs (#1186)
- Remove values-xxx.yaml from Helm packages (#1107)
- Feature/internal net/use netaddr (#1154)
- Add Nop Logger for empty logger message (#1158)
- Remove arm64 build of filter-ingress-tf (#1153)
- add tensorflow filter (#951)
- Add OpenAPIv3 schema to CRDs (#1068)
- PR for review all logging statements (#1052)

#### Bugfix

- bugfix gateway search result merging process (#1238)
- remove nil check in Error() method in internal/errors (#1185)
- Bug fix of parallel test and goleak (#1231)
- :rotating_light: Fix linting errors (#1179)
- Fix mutex deadlock when calling error occurred (#1225)
- fix nil pointer panic when err is nil in Is() (#1216)
- Fix typo & missed schema definition (#1156)
- :lock: Fix validation about zip slip (#1150)
- :lock: Fix security issue: Add validation about path of extracted file (#1145)
- fix: nil pointer bug when config is nil and refactor variable name (#1139)
- :bug: Fix gateway service selectors in Helm chart (#1109)
- Bugfix correct error handling for agent apis (#1144)

#### Document

- separate tutorials into each document and add images (#1230)
- Update Tutotial/Get-Started (#1203)
- Update unit-test-guideline.md (#1213)
- create docs/user-guides/sdks (#1182)
- Fix typo in README. (#1163)

#### CI/CD

- Upgrade to GitHub-native Dependabot (#1211)
- Use kubectl create/replace for upgrading CRDs (#1199)
- :wrench: Fix scylla deploy task (#1159)

#### Test Code

- Add config agent core ngt test (#1219)
- implement pkg/agent/core/ngt/router test & refactor router implementation (#1214)
- Implement pkg/agent/core/ngt/handler/grpc/option test (#1215)
- refactor config error handling (#1190)
- Implement pkg/agent/core/ngt/router/option test (#1206)
- Implement safety bench code (#1171)
- create internal/config/redis test (#1147)
- Implement internal/config/transport test (#1172)
- Add internal/config/sidecar test (#1173)
- Add benchmark for internal/timeutil/time.go (#1086)
- Add internal/config/observability unit test (#1155)
- Add test for internal/config/mysql (#1151)
- create internal/config/cassandra test (#1117)
- create test for internal/config/server (#1175)
- Add internal/config/lb test (#1134)
- create internal/config/net test (#1140)
- :white_check_mark: create internal/config/meta test (#1133)
- Add internal/config/index test (#1129)
- :white_check_mark: create internal/config/backup test (#1132)
- delete internal/config/debug file (#1124)
- create internal/config/grpc test (#1130)
- Add internal/config/discoverer test (#1122)
- Add internal/config/filter.go test (#1125)
- create internal/config/gateway test (#1128)
- Implement internal/info/info.go benchmark (#1093)
- Add internal/config/client.go test (#1121)
- Implement internal/timeutil/location/loc.go benchmark test (#1091)
- :white_check_mark: remove unused variable from mysql test (#1120)
- feat: add blob test and comment (#1114)
- create bench code for internal/rand (#1089)
- Add test internal/config/compress.go (#1097)
- create internal/config/backoff test (#1104)
- change params for passing test of internal/backoff (#1108)
- remove unnecessary tests and update test (#1061)
- add internal backoff test (#1085)
- Add internal/core/algorithm/ngt/util.go test (#1066)
- Fix hack/tools/metrics failed builds (#1094)
- make timeouts for e2e and chaos tests longer (#1092)

#### Dependencies

- update go version to 1.16.4 (#1239)
- :arrow_up: Upgrade helm, valdcli, kubelinter and osdk (#1181)
- update codecov version for vulnerability (#1207)

#### Others

- update go modules (#1242)
- [ImgBot] Optimize images (#1241)
- Update license headers / Format codes (#1212)
- remove ngt version in filter-ingress-tensorflow (#1157)
- remove license from bot config json (#1164)
- improve internal/timeutil/location implementation (#1176)

## v1.0.4

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.0.4</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.0.4</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.0.4</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.0.4</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.0.4</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.0.4</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.0.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.0.4</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.0.4)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.0.4/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.0.4/charts/vald-helm-operator/README.md)

### Changes

- Fix error handling of SearchByID API in lb gateway ([#1084](https://github.com/vdaas/vald/pull/1084))
- Update Agent dashboard ([#1069](https://github.com/vdaas/vald/pull/1069))
- Revise gRPC error statuses and details in meta (Redis/Cassandra) ([#1013](https://github.com/vdaas/vald/pull/1013))
- add remove sample code for tutorial ([#1053](https://github.com/vdaas/vald/pull/1053))
- add grpc reflection ([#1064](https://github.com/vdaas/vald/pull/1064))
- remove vcache for vald agent due to vcache delete timing control failure and time ordered concurrent vector queue called vqueue ([#1028](https://github.com/vdaas/vald/pull/1028))
- refactor discoverer client ([#1056](https://github.com/vdaas/vald/pull/1056))
- bugfix nil pointer no target discovered ([#1055](https://github.com/vdaas/vald/pull/1055))
- Upgrade Operator SDK version to v1.4.2 ([#1038](https://github.com/vdaas/vald/pull/1038))

## v1.0.3

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.0.3</code><br/>
      <code>docker pull vdaas/vald-backup-gateway:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-backup-gateway:v1.0.3</code><br/>
      <code>docker pull vdaas/vald-lb-gateway:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.0.3</code><br/>
      <code>docker pull vdaas/vald-meta-gateway:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-gateway:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Backup managers</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.0.3</code><br/>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Metas</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.0.3</code><br/>
      <code>docker pull vdaas/vald-meta-cassandra:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.0.3</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.0.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.0.3</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.0.3)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.0.3/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.0.3/charts/vald-helm-operator/README.md)

### Changes

- fix MySQL panic test ([#996](https://github.com/vdaas/vald/pull/996))
- Update tutorial to support v1 ([#1009](https://github.com/vdaas/vald/pull/1009))
- remove not supported control flag on darwin ([#1025](https://github.com/vdaas/vald/pull/1025))
- Add strategy section to Docker build workflows ([#1024](https://github.com/vdaas/vald/pull/1024), [#1019](https://github.com/vdaas/vald/pull/1019))
- move internal/net/tcp package to internal/net package and support unix domain socket ([#1010](https://github.com/vdaas/vald/pull/1010))
- Implement internal/info/info test ([#862](https://github.com/vdaas/vald/pull/862))
- add test for internal/errors/runner ([#1007](https://github.com/vdaas/vald/pull/1007))
- add logo guideline ([#973](https://github.com/vdaas/vald/pull/973))
- Fix invalid changelogs / update changelog workflows ([#1002](https://github.com/vdaas/vald/pull/1002))
- Add test case for internal/errors/observability.go ([#993](https://github.com/vdaas/vald/pull/993))

## v1.0.2

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.0.2</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.0.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.0.2</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.0.2)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.0.2/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.0.2/charts/vald-helm-operator/README.md)

### Changes

- v1.0.2 Release ([#998](https://github.com/vdaas/vald/pull/998))
- :robot: Automatically update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE ([#997](https://github.com/vdaas/vald/pull/997))
- Bug fix on StreamGetObject API and DNS cache expiration / refactor net connection ([#986](https://github.com/vdaas/vald/pull/986))
- Remove base docker image ([#995](https://github.com/vdaas/vald/pull/995))
- Use namespaced names for priorityclasses of new gateways ([#992](https://github.com/vdaas/vald/pull/992))
- Add E2E Chaos tests running on GitHub Actions ([#899](https://github.com/vdaas/vald/pull/899))
- Add zap logger to chart schema ([#985](https://github.com/vdaas/vald/pull/985))
- add test for internal/errors/runtime ([#984](https://github.com/vdaas/vald/pull/984))
- add test for internal/errors/tensorflow ([#982](https://github.com/vdaas/vald/pull/982))
- add test for internal/errors/unit ([#979](https://github.com/vdaas/vald/pull/979))
- update chatops permission ([#990](https://github.com/vdaas/vald/pull/990))
- :robot: Automatically update k8s manifests ([#981](https://github.com/vdaas/vald/pull/981))

## v1.0.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.0.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.0.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.0.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.0.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.0.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.0.1/charts/vald-helm-operator/README.md)

### Changes

- bugfix lb-gateway's Insert rpc nil pointer panic ([#980](https://github.com/vdaas/vald/pull/980))
- Implement internal/errors/worker test ([#952](https://github.com/vdaas/vald/pull/952))
- create test for internal/errors/errors.go ([#929](https://github.com/vdaas/vald/pull/929))
- Add test case for internal/errors/net.go ([#969](https://github.com/vdaas/vald/pull/969))
- :robot: Automatically update k8s manifests ([#975](https://github.com/vdaas/vald/pull/975))

## v1.0.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.0.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.0.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.0.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.0.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.0.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.0.0/charts/vald-helm-operator/README.md)

### Changes

- v1.0.0 Release ([#974](https://github.com/vdaas/vald/pull/974))
- Bug fix for gateway ([#963](https://github.com/vdaas/vald/pull/963))
- Add test case for internal/errors/ngt.go ([#936](https://github.com/vdaas/vald/pull/936))
- create new test and refactor for errors/option ([#950](https://github.com/vdaas/vald/pull/950))
- Use checksum of configmap for sidecar-enabled Agents. ([#970](https://github.com/vdaas/vald/pull/970))
- :wrench: Use cass-operator to deploy cassandra for dev cluster ([#968](https://github.com/vdaas/vald/pull/968))
- add test for internal/errors/vald ([#958](https://github.com/vdaas/vald/pull/958))
- Remove actions/cache to improve workflow speed / Refactoring docker-build workflows ([#957](https://github.com/vdaas/vald/pull/957))
- add filter gateway ([#948](https://github.com/vdaas/vald/pull/948))
- :robot: Update license headers / Format Go codes and YAML files ([#954](https://github.com/vdaas/vald/pull/954))
- Add Zap logger and access log interceptor ([#944](https://github.com/vdaas/vald/pull/944))
- Add test case for internal/errors/mysql.go ([#918](https://github.com/vdaas/vald/pull/918))
- Fix invalid data loading code in example ([#949](https://github.com/vdaas/vald/pull/949))
- Add option to use networking.k8s.io/v1beta1 ingresses ([#945](https://github.com/vdaas/vald/pull/945))
- bugfix gateway-lb nil pointer panic due to nil filter configuration ([#943](https://github.com/vdaas/vald/pull/943))
- update rbac.authorization.k8s.io v1beta1 to v1 ([#942](https://github.com/vdaas/vald/pull/942))
- Update Grafana dashboards ([#940](https://github.com/vdaas/vald/pull/940))
- Add a gRPC interceptor for embedding payloads into trace spans ([#900](https://github.com/vdaas/vald/pull/900))
- add goleak.IgnoreCurrent option for Parallel testing ([#941](https://github.com/vdaas/vald/pull/941))
- update go modules ([#939](https://github.com/vdaas/vald/pull/939))
- :robot: Automatically update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE ([#937](https://github.com/vdaas/vald/pull/937))
- :zap: Revise build command for multiplatforms ([#890](https://github.com/vdaas/vald/pull/890))
- update dependencies ([#935](https://github.com/vdaas/vald/pull/935))
- E2E deploy test: rewrite in go tests ([#814](https://github.com/vdaas/vald/pull/814))
- create unit test guideline ([#869](https://github.com/vdaas/vald/pull/869))
- feature/apis/change grpc error object ([#934](https://github.com/vdaas/vald/pull/934))
- Add test case for internal/errors/http.go ([#908](https://github.com/vdaas/vald/pull/908))
- :robot: Update license headers / Format Go codes and YAML files ([#933](https://github.com/vdaas/vald/pull/933))
- revise kubelinter config / Add securityContext section to Helm chart ([#833](https://github.com/vdaas/vald/pull/833))
- :robot: Update license headers / Format Go codes and YAML files ([#932](https://github.com/vdaas/vald/pull/932))
- :robot: Automatically update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE ([#931](https://github.com/vdaas/vald/pull/931))
- os.free nil pointer failure in ngt cgo due to create index hang up ([#930](https://github.com/vdaas/vald/pull/930))
- change grpc bidi-stream error handling and change grpc API interface ([#928](https://github.com/vdaas/vald/pull/928))
- fix unclosed string literal in Dockerfile's ARG MAINTAINER ([#923](https://github.com/vdaas/vald/pull/923))
- Revise building workflow of ci and dev containers ([#922](https://github.com/vdaas/vald/pull/922))
- bugfix add nil check for grpc connection pool objects in grpc/client.go ([#921](https://github.com/vdaas/vald/pull/921))
- remove unneccessary pr-tag definition from chart ([#920](https://github.com/vdaas/vald/pull/920))
- :pencil: Fix typo in gateway-vald configmap template ([#919](https://github.com/vdaas/vald/pull/919))
- change docker base image PRIMARY_TAG name from nightly to latest ([#917](https://github.com/vdaas/vald/pull/917))
- :green_heart: Use vdaas-ci token for making commits ([#895](https://github.com/vdaas/vald/pull/895))
- Vald V1 New Design APIs ([#826](https://github.com/vdaas/vald/pull/826))
- :robot: Automatically update k8s manifests ([#914](https://github.com/vdaas/vald/pull/914))

## v0.0.66

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.66</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.66</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.66</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.66)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.66/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.66/charts/vald-helm-operator/README.md)

### Changes

- bugfix: do not create metadata file when create/append flag is not set for the agent/agent-sidecar ([#904](https://github.com/vdaas/vald/pull/904))
- :robot: Update license headers / Format Go codes and YAML files ([#913](https://github.com/vdaas/vald/pull/913))
- :green_heart: Add formatter for main branch ([#911](https://github.com/vdaas/vald/pull/911))
- :page_facing_up: Update license headers for .github yamls ([#907](https://github.com/vdaas/vald/pull/907))
- Add test case for internal/errors/io.go ([#910](https://github.com/vdaas/vald/pull/910))
- Add test case for internal/errors/grpc.go ([#903](https://github.com/vdaas/vald/pull/903))
- :robot: Automatically update k8s manifests ([#906](https://github.com/vdaas/vald/pull/906))

## v0.0.65

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.65</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.65</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.65</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.65)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.65/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.65/charts/vald-helm-operator/README.md)

### Changes

- Happy new year ([#905](https://github.com/vdaas/vald/pull/905))
- :white_check_mark: Add test case for internal/errors/file.go ([#893](https://github.com/vdaas/vald/pull/893))
- :robot: Automatically update k8s manifests ([#902](https://github.com/vdaas/vald/pull/902))

## v0.0.64

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.64</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.64</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.64</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.64)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.64/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.64/charts/vald-helm-operator/README.md)

### Changes

- :wrench: Rename PriorityClass names to contain namespace ([#901](https://github.com/vdaas/vald/pull/901))
- :bug: Fix bug on updating status of VR & VHOR resources ([#892](https://github.com/vdaas/vald/pull/892))
- :white_check_mark: Implement internal/errors/blob test ([#888](https://github.com/vdaas/vald/pull/888))
- :white_check_mark: Add test for cassandra error ([#865](https://github.com/vdaas/vald/pull/865))
- :white_check_mark: Add test case for internal/errors/discoverer.go ([#874](https://github.com/vdaas/vald/pull/874))
- :bug: :white_check_mark: remove invalid import package from internal/errors/mysql_test.go ([#894](https://github.com/vdaas/vald/pull/894))
- :white_check_mark: Add test for internal/errors/compressor.go ([#870](https://github.com/vdaas/vald/pull/870))
- :robot: Automatically update k8s manifests ([#887](https://github.com/vdaas/vald/pull/887))

## v0.0.63

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.63</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.63</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.63</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.63)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.63/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.63/charts/vald-helm-operator/README.md)

### Changes

- Remove go mod tidy from base image / Add valid token to semver job ([#886](https://github.com/vdaas/vald/pull/886))
- bugfix agent duplicated data update execution ([#885](https://github.com/vdaas/vald/pull/885))
- :package: Build schemas when building helm-operator image ([#879](https://github.com/vdaas/vald/pull/879))
- :white_check_mark: Add s3 test and Refactor ([#837](https://github.com/vdaas/vald/pull/837))
- :white_check_mark: Add internal/net/http/client test ([#858](https://github.com/vdaas/vald/pull/858))
- :white_check_mark: add test case for json package ([#857](https://github.com/vdaas/vald/pull/857))
- Add FOSSA scan workflow & .fossa.yml ([#846](https://github.com/vdaas/vald/pull/846))
- :white_check_mark: create internal/net/http/client option test ([#831](https://github.com/vdaas/vald/pull/831))
- :green_heart: fix checkout-v2 fetch depths ([#832](https://github.com/vdaas/vald/pull/832))
- :wrench: update Helm 3.4.1, helm-docs ([#829](https://github.com/vdaas/vald/pull/829))
- :white_check_mark: test/internal/nosql/cassandra test ([#809](https://github.com/vdaas/vald/pull/809))
- :pencil: fix coding style for mock ([#806](https://github.com/vdaas/vald/pull/806))
- CI: Add reviewdog-k8s ([#824](https://github.com/vdaas/vald/pull/824))
- Fix e2e-bench-agent CI to fail correctly ([#800](https://github.com/vdaas/vald/pull/800))
- Add internal/db/cassandra/conviction test ([#799](https://github.com/vdaas/vald/pull/799))
- add test ([#798](https://github.com/vdaas/vald/pull/798))
- :pencil: Fix coding guideline about constructor due to mock implementation ([#792](https://github.com/vdaas/vald/pull/792))
- :robot: Automatically update k8s manifests ([#781](https://github.com/vdaas/vald/pull/781))

## v0.0.62

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.62</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.62</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.62</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.62)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.62/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.62/charts/vald-helm-operator/README.md)

### Changes

- add 3 new distance type support for agent-ngt ([#780](https://github.com/vdaas/vald/pull/780))
- upgrade KinD, Helm, valdcli, telepresence, tensorlfow, operator-sdk, helm-docs ([#776](https://github.com/vdaas/vald/pull/776))
- :robot: Automatically update k8s manifests ([#774](https://github.com/vdaas/vald/pull/774))

## v0.0.61

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.61</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.61</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.61</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.61)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.61/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.61/charts/vald-helm-operator/README.md)

### Changes

- fix search result sorting codes ([#772](https://github.com/vdaas/vald/pull/772))
- :robot: Automatically update k8s manifests ([#771](https://github.com/vdaas/vald/pull/771))

## v0.0.60

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Discoverer k8s</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Gateway</td>
    <td>
      <code>docker pull vdaas/vald-gateway:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-gateway:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager MySQL</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-mysql:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-mysql:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Backup manager Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-manager-backup-cassandra:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-backup-cassandra:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Compressor</td>
    <td>
      <code>docker pull vdaas/vald-manager-compressor:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-compressor:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Meta Redis</td>
    <td>
      <code>docker pull vdaas/vald-meta-redis:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-redis:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Meta Cassandra</td>
    <td>
      <code>docker pull vdaas/vald-meta-cassandra:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-meta-cassandra:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v0.0.60</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v0.0.60</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v0.0.60</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.60)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.60/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.60/charts/vald-helm-operator/README.md)

### Changes

- Fix fails test for s3 reader ([#770](https://github.com/vdaas/vald/pull/770))
- CI: Make docker builds fast again ([#756](https://github.com/vdaas/vald/pull/756))
- :robot: Automatically update k8s manifests ([#769](https://github.com/vdaas/vald/pull/769))

## v0.0.59

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.59`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.59`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.59`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.59`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.59`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.59` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.59`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.59`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.59`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.59`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.59`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.59)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.59/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.59/charts/vald-helm-operator/README.md)

### Changes

- bugfix gateway index out of bounds ([#768](https://github.com/vdaas/vald/pull/768))
- :robot: Automatically update k8s manifests ([#766](https://github.com/vdaas/vald/pull/766))

## v0.0.58

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.58`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.58`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.58`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.58`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.58`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.58` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.58`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.58`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.58`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.58`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.58`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.58)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.58/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.58/charts/vald-helm-operator/README.md)

### Changes

- change gateway vald's mutex lock ([#765](https://github.com/vdaas/vald/pull/765))
- patch add more effective Close function for internal/core ([#764](https://github.com/vdaas/vald/pull/764))
- :white_check_mark: Fix mysql test failure ([#750](https://github.com/vdaas/vald/pull/750))
- :green_heart: remove deprecated set-env commands ([#752](https://github.com/vdaas/vald/pull/752))
- bugfix discoverer nil map reference ([#745](https://github.com/vdaas/vald/pull/745))
- add test of internal/db/rdb/mysql ([#659](https://github.com/vdaas/vald/pull/659))
- :white_check_mark: :recycle: Add test for internal/db/storage/blob/s3/reader ([#718](https://github.com/vdaas/vald/pull/718))
- :white_check_mark: Cassandra option test (part 2) ([#724](https://github.com/vdaas/vald/pull/724))
- CI: Build multi-platform Docker images ([#727](https://github.com/vdaas/vald/pull/727))
- :white_check_mark: Add test for s3 session option ([#736](https://github.com/vdaas/vald/pull/736))
- :white_check_mark: Create s3/session test ([#702](https://github.com/vdaas/vald/pull/702))
- :fire: remove dependencies to gql proto ([#731](https://github.com/vdaas/vald/pull/731))
- :robot: Automatically update k8s manifests ([#730](https://github.com/vdaas/vald/pull/730))

## v0.0.57

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.57`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.57`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.57`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.57`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.57`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.57` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.57`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.57`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.57`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.57`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.57`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.57)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.57/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.57/charts/vald-helm-operator/README.md)

### Changes

- fix duplicated search result ([#729](https://github.com/vdaas/vald/pull/729))
- :recycle: enable to inject only agent-sidecar on initContainer mode without enabling sidecar mode ([#726](https://github.com/vdaas/vald/pull/726))
- :sparkles: implement billion scale data ([#612](https://github.com/vdaas/vald/pull/612))
- Add devcontainer ([#620](https://github.com/vdaas/vald/pull/620))
- :white_check_makr: :recycle: Add test for s3/writer and Refactor. ([#672](https://github.com/vdaas/vald/pull/672))
- CI-container: upgrade dependencies of & remove workdir contents ([#711](https://github.com/vdaas/vald/pull/711))
- :robot: Automatically update k8s manifests ([#708](https://github.com/vdaas/vald/pull/708))

## v0.0.56

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.56`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.56`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.56`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.56`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.56`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.56` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.56`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.56`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.56`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.56`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.56`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.56)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.56/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.56/charts/vald-helm-operator/README.md)

### Changes

- add C.free & delete ivc before core.BulkInsert C' function executing for reducing memory usage ([#701](https://github.com/vdaas/vald/pull/701))
- Add cassandra option test ([#644](https://github.com/vdaas/vald/pull/644))
- :memo: build single artifact from pbdocs task ([#699](https://github.com/vdaas/vald/pull/699))
- improve CI builds: use DOCKER_BUILDKIT ([#706](https://github.com/vdaas/vald/pull/706))
- :robot: Automatically update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE ([#695](https://github.com/vdaas/vald/pull/695))
- :pencil: add AVX2 to Requirements section ([#686](https://github.com/vdaas/vald/pull/686))
- :pencil: update contributing guide ([#678](https://github.com/vdaas/vald/pull/678))
- Use runtime.GC for reducing indexing memory & replace saveMu with atomic busy loop for race control ([#682](https://github.com/vdaas/vald/pull/682))
- :white_check_mark: :recycle: Implement zstd test ([#676](https://github.com/vdaas/vald/pull/676))
- :sparkles: use internal client ([#618](https://github.com/vdaas/vald/pull/618))
- :robot: Automatically update k8s manifests ([#684](https://github.com/vdaas/vald/pull/684))

## v0.0.55

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.55`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.55`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.55`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.55`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.55`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.55` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.55`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.55`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.55`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.55`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.55`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.55)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.55/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.55/charts/vald-helm-operator/README.md)

### Changes

- pass CFLAGS, CXXFLAGS to NGT build command ([#683](https://github.com/vdaas/vald/pull/683))
- :robot: Automatically update k8s manifests ([#681](https://github.com/vdaas/vald/pull/681))

## v0.0.54

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.54`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.54`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.54`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.54`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.54`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.54` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.54`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.54`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.54`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.54`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.54`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.54)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.54/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.54/charts/vald-helm-operator/README.md)

### Changes

- bugfix error assertion ([#680](https://github.com/vdaas/vald/pull/680))
- :robot: Automatically update k8s manifests ([#679](https://github.com/vdaas/vald/pull/679))

## v0.0.53

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.53`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.53`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.53`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.53`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.53`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.53` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.53`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.53`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.53`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.53`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.53`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.53)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.53/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.53/charts/vald-helm-operator/README.md)

### Changes

- remove cockroachdb/errors ([#677](https://github.com/vdaas/vald/pull/677))
- :white_check_mark: Add test case for storage/blob/s3/writer/option ([#656](https://github.com/vdaas/vald/pull/656))
- :white_check_mark: fix: failing tset ([#671](https://github.com/vdaas/vald/pull/671))
- :bug: fix & upgrade manifests to operator-sdk v1.0.0 compatible ([#667](https://github.com/vdaas/vald/pull/667))
- :robot: Automatically update k8s manifests ([#666](https://github.com/vdaas/vald/pull/666))

## v0.0.52

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.52`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.52`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.52`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.52`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.52`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.52` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.52`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.52`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.52`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.52`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.52`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.52)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.52/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.52/charts/vald-helm-operator/README.md)

### Changes

- add build stage for operator-sdk docker v1.0.0 permission changes ([#665](https://github.com/vdaas/vald/pull/665))
- :robot: Automatically update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE ([#663](https://github.com/vdaas/vald/pull/663))
- :robot: Automatically update k8s manifests ([#664](https://github.com/vdaas/vald/pull/664))

## v0.0.51

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.51`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.51`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.51`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.51`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.51`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.51` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.51`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.51`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.51`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.51`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.51`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.51)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.51/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.51/charts/vald-helm-operator/README.md)

### Changes

- update deps ([#660](https://github.com/vdaas/vald/pull/660))
- add metrics for indexer and sidecar ([#642](https://github.com/vdaas/vald/pull/642))
- :pencil2: fix indents in helm chart of vald-helm-operator ([#658](https://github.com/vdaas/vald/pull/658))
- :white_check_mark: Add test for internal/compress/gob.go ([#646](https://github.com/vdaas/vald/pull/646))
- Upgrade go mod default: k8s.io/xxx v0.18.8 ([#645](https://github.com/vdaas/vald/pull/645))
- [Coding guideline] Add implementation and grouping section ([#641](https://github.com/vdaas/vald/pull/641))
- :white_check_mark: add internal/compress/lz4 test ([#643](https://github.com/vdaas/vald/pull/643))
- Add test case for `s3/option.go` ([#640](https://github.com/vdaas/vald/pull/640))
- Refactoring and Add test code for `compress` ([#622](https://github.com/vdaas/vald/pull/622))
- [ImgBot] Optimize images ([#639](https://github.com/vdaas/vald/pull/639))
- Add operation guide ([#541](https://github.com/vdaas/vald/pull/541))
- :white_check_mark: add internal/s3/reader/option test ([#630](https://github.com/vdaas/vald/pull/630))
- :bug: Fix indexer's creation_pool_size field ([#637](https://github.com/vdaas/vald/pull/637))
- :wrench: revise languagetool rules: disable EN_QUOTES ([#635](https://github.com/vdaas/vald/pull/635))
- :wrench: revise languagetool rules: disable TYPOS, DASH_RULE ([#634](https://github.com/vdaas/vald/pull/634))
- :wrench: revise languagetool rules ([#633](https://github.com/vdaas/vald/pull/633))
- [ImgBot] Optimize images ([#632](https://github.com/vdaas/vald/pull/632))
- Add upsert flow in architecture doc ([#627](https://github.com/vdaas/vald/pull/627))
- Add DB metrics & traces: Redis, MySQL ([#623](https://github.com/vdaas/vald/pull/623))
- :white_check_mark: add internal/db/rdb/mysql/model test ([#628](https://github.com/vdaas/vald/pull/628))
- :wrench: upload sarif only for HIGH or CRITICAL ([#629](https://github.com/vdaas/vald/pull/629))
- :white_check_mark: add internal/db/rdb/mysql/option test ([#626](https://github.com/vdaas/vald/pull/626))
- Add internal/net test ([#615](https://github.com/vdaas/vald/pull/615))
- :white_check_mark: Add test for gzip_option ([#625](https://github.com/vdaas/vald/pull/625))
- :pencil: change showing image method ([#624](https://github.com/vdaas/vald/pull/624))
- :white_check_mark: Add internal/compress/zstd_option test ([#621](https://github.com/vdaas/vald/pull/621))
- use distroless for base image ([#605](https://github.com/vdaas/vald/pull/605))
- :pencil: Coding guideline: Add error checking section ([#614](https://github.com/vdaas/vald/pull/614))
- :white_check_mark: add internal/compress/lz4_option test ([#619](https://github.com/vdaas/vald/pull/619))
- :white_check_mark: fix test fail ([#616](https://github.com/vdaas/vald/pull/616))
- :white_check_mark: Add test of internal/worker/queue_option ([#613](https://github.com/vdaas/vald/pull/613))
- [ImgBot] Optimize images ([#617](https://github.com/vdaas/vald/pull/617))
- :pencil: Add update dataflow in architecture document ([#601](https://github.com/vdaas/vald/pull/601))
- Add internal/worker/worker test ([#602](https://github.com/vdaas/vald/pull/602))
- :recycle: refactor load test ([#552](https://github.com/vdaas/vald/pull/552))
- :white_check_mark: Add test case for internal/kvs/redis/option.go ([#611](https://github.com/vdaas/vald/pull/611))
- :white_check_mark: create internal/worker/queue test ([#606](https://github.com/vdaas/vald/pull/606))
- :white_check_mark: :recycle: Add internal roundtrip test code ([#589](https://github.com/vdaas/vald/pull/589))
- :pencil: Documentation/performance/loadtest ([#610](https://github.com/vdaas/vald/pull/610))
- :robot: Automatically update k8s manifests ([#609](https://github.com/vdaas/vald/pull/609))

## v0.0.50

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.50`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.50`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.50`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.50`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.50`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.50` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.50`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.50`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.50`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.50`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.50`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.50)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.50/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.50/charts/vald-helm-operator/README.md)

### Changes

- Add warn logging messages to agent-sidecar & ignore io.EOF error when reading metadata.json ([#608](https://github.com/vdaas/vald/pull/608))
- Add DB metrics: Cassandra ([#587](https://github.com/vdaas/vald/pull/587))
- :recycle: Improve Singleflight performance ([#580](https://github.com/vdaas/vald/pull/580))
- [ImgBot] Optimize images ([#607](https://github.com/vdaas/vald/pull/607))
- :pencil: add delete dataflow in architecture document ([#591](https://github.com/vdaas/vald/pull/591))
- :recycle: Add gaussian random vector generation ([#595](https://github.com/vdaas/vald/pull/595))
- :recycle: :white_check_mark: Refactoring `internal/db/kvs/redis` package ([#590](https://github.com/vdaas/vald/pull/590))
- :green_heart: Add reviewdog - markdown: LanguageTool ([#604](https://github.com/vdaas/vald/pull/604))
- :green_heart: Add reviewdog - hadolint ([#603](https://github.com/vdaas/vald/pull/603))
- :robot: Automatically update k8s manifests ([#599](https://github.com/vdaas/vald/pull/599))

## v0.0.49

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.49`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.49`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.49`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.49`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.49`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.49` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.49`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.49`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.49`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.49`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.49`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.49)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.49/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.49/charts/vald-helm-operator/README.md)

### Changes

- :bug: fix agent sidecar behavior ([#598](https://github.com/vdaas/vald/pull/598))
- :robot: Automatically update k8s manifests ([#597](https://github.com/vdaas/vald/pull/597))

## v0.0.48

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.48`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.48`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.48`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.48`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.48`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.48` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.48`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.48`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.48`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.48`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.48`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.48)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.48/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.48/charts/vald-helm-operator/README.md)

### Changes

- :bug: fix behavior when index path is empty ([#596](https://github.com/vdaas/vald/pull/596))
- :white_check_mark: add internal/net/http/transport/option test ([#594](https://github.com/vdaas/vald/pull/594))
- tensorflow savedmodel warmup ([#539](https://github.com/vdaas/vald/pull/539))
- :robot: Automatically update k8s manifests ([#592](https://github.com/vdaas/vald/pull/592))

## v0.0.47

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.47`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.47`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.47`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.47`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.47`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.47` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.47`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.47`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.47`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.47`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.47`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.47)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.47/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.47/charts/vald-helm-operator/README.md)

### Changes

- [agent-NGT, sidecar] Improve S3 backup/recover behavior ([#556](https://github.com/vdaas/vald/pull/556))
- :white_check_mark: add internal/cache/option test ([#586](https://github.com/vdaas/vald/pull/586))
- :robot: Automatically update k8s manifests ([#588](https://github.com/vdaas/vald/pull/588))

## v0.0.46

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.46`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.46`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.46`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.46`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.46`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.46` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.46`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.46`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.46`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.46`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.46`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.46)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.46/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.46/charts/vald-helm-operator/README.md)

### Changes

- Test/internal/tcp ([#501](https://github.com/vdaas/vald/pull/501))
- :white_check_mark: add internal/cache/gache test ([#583](https://github.com/vdaas/vald/pull/583))
- :white_check_mark: add cache test ([#576](https://github.com/vdaas/vald/pull/576))
- :white_check_mark: add internal/cache/gache/option test ([#575](https://github.com/vdaas/vald/pull/575))
- :bug: :white_check_mark: fix: fails test ([#578](https://github.com/vdaas/vald/pull/578))
- :white_check_mark: Add test for `internal/file/watch` ([#526](https://github.com/vdaas/vald/pull/526))
- :art: update k8s manifests only on publish tags ([#574](https://github.com/vdaas/vald/pull/574))
- :robot: Automatically update k8s manifests ([#572](https://github.com/vdaas/vald/pull/572))
- :robot: Automatically update k8s manifests ([#571](https://github.com/vdaas/vald/pull/571))
- :robot: Automatically update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE ([#570](https://github.com/vdaas/vald/pull/570))

## v0.0.45

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.45`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.45`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.45`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.45`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.45`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.45` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.45`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.45`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.45`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.45`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.45`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.45)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.45/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.45/charts/vald-helm-operator/README.md)

### Changes

- bugfix gateway & internal/net/grpc ([#569](https://github.com/vdaas/vald/pull/569))
- fix update-k8s workflow & update sample manifests ([#567](https://github.com/vdaas/vald/pull/567))
- :white_check_mark: Add test for `internal/config/mysql.go` ([#563](https://github.com/vdaas/vald/pull/563))
- :bug: :white_check_mark: fix failed test ([#561](https://github.com/vdaas/vald/pull/561))
- :white_check_mark: internal/tls test ([#485](https://github.com/vdaas/vald/pull/485))
- pass tparse by tee command ([#562](https://github.com/vdaas/vald/pull/562))
- fix global cache ([#560](https://github.com/vdaas/vald/pull/560))
- :white_check_mark: add internal/config/ngt test ([#554](https://github.com/vdaas/vald/pull/554))
- :white_check_mark: internal/cache/cacher test ([#553](https://github.com/vdaas/vald/pull/553))
- :white_check_mark: Add test case for `internal/file` ([#550](https://github.com/vdaas/vald/pull/550))
- :white_check_mark: add internal/singleflight test ([#542](https://github.com/vdaas/vald/pull/542))
- not to force rebuild gotests ([#548](https://github.com/vdaas/vald/pull/548))
- :pencil: Add use case document ([#482](https://github.com/vdaas/vald/pull/482))
- :white_check_mark: add internal/log/mock/retry test ([#549](https://github.com/vdaas/vald/pull/549))
- feat: options test ([#518](https://github.com/vdaas/vald/pull/518))
- :white_check_mark: add log/mock/logger test ([#538](https://github.com/vdaas/vald/pull/538))
- :bug: Fix condition check of chatops ([#544](https://github.com/vdaas/vald/pull/544))
- exclude hack codes ([#543](https://github.com/vdaas/vald/pull/543))

## v0.0.44

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.44`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.44`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.44`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.44`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.44`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.44` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.44`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.44`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.44`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.44`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.44`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.44)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.44/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.44/charts/vald-helm-operator/README.md)

### Changes

- use Len and InsertVCacheLen method for IndexInfo / add mutex for (Create|Save)Index ([#536](https://github.com/vdaas/vald/pull/536))
- documentation: tutorial/agent-on-docker ([#516](https://github.com/vdaas/vald/pull/516))
- Revise log messages along with the coding guideline ([#504](https://github.com/vdaas/vald/pull/504))
- :art: Add images of usecase ([#537](https://github.com/vdaas/vald/pull/537))
- :white_check_mark: add internal/config/tls test ([#534](https://github.com/vdaas/vald/pull/534))
- :bug: Add cancel hook for file watcher ([#535](https://github.com/vdaas/vald/pull/535))
- :green_heart: Add test workflow ([#531](https://github.com/vdaas/vald/pull/531))
- added internal/config/log test ([#530](https://github.com/vdaas/vald/pull/530))
- add codeql config ([#532](https://github.com/vdaas/vald/pull/532))
- Added test case for `internal/info` pacakge. ([#514](https://github.com/vdaas/vald/pull/514))
- Add `internal/runner` test ([#505](https://github.com/vdaas/vald/pull/505))
- Added test case for `internal/unit` pacakge ([#515](https://github.com/vdaas/vald/pull/515))

## v0.0.43

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.43`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.43`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.43`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.43`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.43`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.43` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.43`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.43`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.43`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.43`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.43`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.43)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.43/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.43/charts/vald-helm-operator/README.md)

### Changes

- Revise S3 reader/writer: compatible with IBM Cloud Object Storage ([#509](https://github.com/vdaas/vald/pull/509))
- :bug: Close [#502](https://github.com/vdaas/vald/pull/502) / Fix roundtrip error handling (#508)
- Feature/drawio ([#500](https://github.com/vdaas/vald/pull/500))
- Added test case for `internal/errorgroup` ([#494](https://github.com/vdaas/vald/pull/494))
- Update Helm Chart info ([#496](https://github.com/vdaas/vald/pull/496))
- Revise triggers of workflow run & Fix reading changelogs from PR comments ([#495](https://github.com/vdaas/vald/pull/495))

## v0.0.42

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.42`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.42`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.42`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.42`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.42`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.42` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.42`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.42`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.42`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.42`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.42`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.42)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.42/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.42/charts/vald-helm-operator/README.md)

### Changes

- ✨ Add Stackdriver Monitoring, Tracing and Profiler support ([#479](https://github.com/vdaas/vald/pull/479))
- :green_heart: Add CodeQL workflow instead of LGTM.com ([#486](https://github.com/vdaas/vald/pull/486))
- Add `internal/params` pacakge test ([#474](https://github.com/vdaas/vald/pull/474))
- :sparkles: aws region can be specified with empty string ([#477](https://github.com/vdaas/vald/pull/477))
- Fix failed test case of internal/safety package ([#464](https://github.com/vdaas/vald/pull/464))
- send a request to GoProxy after a new version is published ([#475](https://github.com/vdaas/vald/pull/475))
- internal/db/storage/blob/s3: remove ctx from struct ([#473](https://github.com/vdaas/vald/pull/473))

## v0.0.41

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.41`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.41`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.41`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.41`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.41`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.41` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.41`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.41`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.41`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.41`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.41`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.41)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.41/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.41/charts/vald-helm-operator/README.md)

### Changes

- Refactor agent-sidecar: fix S3 reader & add backoff logic ([#467](https://github.com/vdaas/vald/pull/467))
- :bug: :pencil: fix link ([#471](https://github.com/vdaas/vald/pull/471))
- Fix /changelog command format ([#470](https://github.com/vdaas/vald/pull/470))
- fix: failing test ([#469](https://github.com/vdaas/vald/pull/469))
- fix: failing test ([#468](https://github.com/vdaas/vald/pull/468))
- ✨ Add options for AWS client ([#460](https://github.com/vdaas/vald/pull/460))
- Fix /format command ([#466](https://github.com/vdaas/vald/pull/466))
- Fix /format command ([#465](https://github.com/vdaas/vald/pull/465))
- Fix `internal/log/retry` pacakge ([#458](https://github.com/vdaas/vald/pull/458))
- [ImgBot] Optimize images ([#461](https://github.com/vdaas/vald/pull/461))
- :art: trim white margin at data flow images ([#459](https://github.com/vdaas/vald/pull/459))

## v0.0.40

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.40`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.40`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.40`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.40`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.40`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.40` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.40`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.40`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.40`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.40`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.40`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.40)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.40/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.40/charts/vald-helm-operator/README.md)

### Changes

- Documentation: add concept section to the architecture document ([#438](https://github.com/vdaas/vald/pull/438))
- feat: level pacakge test ([#455](https://github.com/vdaas/vald/pull/455))
- [ImgBot] Optimize images ([#457](https://github.com/vdaas/vald/pull/457))
- fix document and added png images ([#456](https://github.com/vdaas/vald/pull/456))
- Fix test template bug ([#452](https://github.com/vdaas/vald/pull/452))
- Add WithOperation func to loadtest usecase ([#454](https://github.com/vdaas/vald/pull/454))
- :bug: Fix k8s manifests for loadtest jobs ([#453](https://github.com/vdaas/vald/pull/453))
- bugfix change final stage of loadtest ([#451](https://github.com/vdaas/vald/pull/451))
- :bug: Fix bug on /changelog command ([#450](https://github.com/vdaas/vald/pull/450))
- add loadtest job and container ([#449](https://github.com/vdaas/vald/pull/449))
- 🐛 Fix bug on changelog command ([#448](https://github.com/vdaas/vald/pull/448))

## v0.0.39

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.39`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.39`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.39`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.39`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.39`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.39` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.39`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.39`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.39`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.39`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.39`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.39)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.39/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.39/charts/vald-helm-operator/README.md)

### Changes

- [patch] fix doc file path ([#444](https://github.com/vdaas/vald/pull/444))
- Add changelog command (ChatOps) ([#447](https://github.com/vdaas/vald/pull/447))
- Fix inconsistent wording ([#442](https://github.com/vdaas/vald/pull/442))
- [ImgBot] Optimize images ([#443](https://github.com/vdaas/vald/pull/443))
- [Document] Apply design template to flow diagram ([#441](https://github.com/vdaas/vald/pull/441))
- Document to deploy standalone agent ([#407](https://github.com/vdaas/vald/pull/407))
- implement load tester prototype ([#363](https://github.com/vdaas/vald/pull/363))
- 🐛 Add gRPC interceptor to recover panic in handlers ([#440](https://github.com/vdaas/vald/pull/440))
- tensorflow test ([#378](https://github.com/vdaas/vald/pull/378))
- :bento: update architecture overview svg to add agent sidecar ([#437](https://github.com/vdaas/vald/pull/437))
- Example program: Add indexing interval description & fix logging message ([#405](https://github.com/vdaas/vald/pull/405))
- :pencil2: Fix typo ([#436](https://github.com/vdaas/vald/pull/436))

## v0.0.38

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.38`                |
| agent sidecar            | `docker pull vdaas/vald-agent-sidecar:v0.0.38`            |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.38`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.38`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.38`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.38` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.38`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.38`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.38`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.38`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.38`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.38)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.38/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.38/charts/vald-helm-operator/README.md)

### Changes

- send PR when K8s manifests are updated ([#435](https://github.com/vdaas/vald/pull/435))
- Implementation of agent-sidecar storage backup logic ([#409](https://github.com/vdaas/vald/pull/409))
- Fix structure of grpc java package ([#431](https://github.com/vdaas/vald/pull/431))
- Remove AUTHORS/CONTRIBUTORS file ([#428](https://github.com/vdaas/vald/pull/428))
- Contribute document ([#390](https://github.com/vdaas/vald/pull/390))
- Separate the component section from architecture doc ([#430](https://github.com/vdaas/vald/pull/430))
- fix: delete fch channel because fch causes channel blocking ([#429](https://github.com/vdaas/vald/pull/429))
- Update Operator SDK version ([#412](https://github.com/vdaas/vald/pull/412))
- Documentation: About vald ([#374](https://github.com/vdaas/vald/pull/374))
- Vald architecture document ([#366](https://github.com/vdaas/vald/pull/366))
- Add JSON schema for Vald Helm Chart ([#365](https://github.com/vdaas/vald/pull/365))
- Revise ChatOps not to add go.mod & go.sum when /format runs ([#406](https://github.com/vdaas/vald/pull/406))
- add agent sidecar flame for implementation ([#404](https://github.com/vdaas/vald/pull/404))
- Upgrade Operator SDK version / Remove useless GO111MODULE=off ([#402](https://github.com/vdaas/vald/pull/402))
- Upgrade tools version ([#399](https://github.com/vdaas/vald/pull/399))
- Add app.kubernetes.io/xxx labels to all resources ([#397](https://github.com/vdaas/vald/pull/397))
- Vald contacts document ([#373](https://github.com/vdaas/vald/pull/373))
- add trace spans and metrics for agent-ngt and index-manager ([#389](https://github.com/vdaas/vald/pull/389))
- Add gen-test command for chatops ([#379](https://github.com/vdaas/vald/pull/379))
- Add internal/db/storage/blob ([#388](https://github.com/vdaas/vald/pull/388))

## v0.0.37

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.37`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.37`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.37`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.37`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.37` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.37`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.37`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.37`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.37`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.37`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.37)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.37/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.37/charts/vald-helm-operator/README.md)

### Changes

- add agent auto save indexing feature ([#385](https://github.com/vdaas/vald/pull/385))
- :bug: fix ngt `distance_type` ([#384](https://github.com/vdaas/vald/pull/384))
- Add topology spread constraints ([#383](https://github.com/vdaas/vald/pull/383))

## v0.0.36

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.36`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.36`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.36`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.36`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.36` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.36`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.36`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.36`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.36`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.36`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.36)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.36/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.36/charts/vald-helm-operator/README.md)

### Changes

- update dependencies version ([#381](https://github.com/vdaas/vald/pull/381))
- Fix missing value on compressor health servers ([#377](https://github.com/vdaas/vald/pull/377))
- Fix compressor readiness shutdown_duration / Fix cassandra … ([#376](https://github.com/vdaas/vald/pull/376))
- Bump gopkg.in/yaml.v2 from 2.2.8 to 2.3.0 ([#375](https://github.com/vdaas/vald/pull/375))
- Fix`internal/log/format` to match the test template ([#369](https://github.com/vdaas/vald/pull/369))
- Fix `internal/log/logger` to match the test template ([#371](https://github.com/vdaas/vald/pull/371))
- Fix failing tests of `internal/log` and modified to match the test template ([#368](https://github.com/vdaas/vald/pull/368))
- Add enabled flag to each component in Helm chart ([#372](https://github.com/vdaas/vald/pull/372))
- Add configurations.md ([#356](https://github.com/vdaas/vald/pull/356))

## v0.0.35

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.35`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.35`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.35`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.35`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.35` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.35`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.35`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.35`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.35`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.35`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.35)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.35/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.35/charts/vald-helm-operator/README.md)

### Changes

- add storage backup option to agent ([#367](https://github.com/vdaas/vald/pull/367))
- Add client-node dispatcher ([#370](https://github.com/vdaas/vald/pull/370))
- Bump github.com/tensorflow/tensorflow ([#364](https://github.com/vdaas/vald/pull/364))
- change fmt.Errorf to errors.Errorf ([#361](https://github.com/vdaas/vald/pull/361))
- add goleak ([#359](https://github.com/vdaas/vald/pull/359))

## v0.0.34

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.34`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.34`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.34`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.34`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.34` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.34`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.34`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.34`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.34`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.34`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.34)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.34/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.34/charts/vald-helm-operator/README.md)

### Changes

- feature/internal/cassandra/add option ([#358](https://github.com/vdaas/vald/pull/358))
- update helm docs when version is published ([#355](https://github.com/vdaas/vald/pull/355))
- upgrade tools ([#354](https://github.com/vdaas/vald/pull/354))
- bugfix protoc-gen-validate resolve failure ([#353](https://github.com/vdaas/vald/pull/353))
- Fix conflicts between formatter and helm template ([#350](https://github.com/vdaas/vald/pull/350))
- Add more options and remove valdhelmoperatorrelease, valdrelease from vald-helm-operator chart ([#334](https://github.com/vdaas/vald/pull/334))

## v0.0.33

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.33`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.33`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.33`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.33`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.33` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.33`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.33`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.33`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.33`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.33`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.33)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.33/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.33/charts/vald-helm-operator/README.md)

### Changes

- update k8s dependencies ([#349](https://github.com/vdaas/vald/pull/349))
- create missing test files by the our original test template ([#348](https://github.com/vdaas/vald/pull/348))
- create test template for using gotests ([#327](https://github.com/vdaas/vald/pull/327))
- Revise coverage CI settings ([#347](https://github.com/vdaas/vald/pull/347))
- fix tensorflow.go, option.go ([#261](https://github.com/vdaas/vald/pull/261))

## v0.0.32

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.32`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.32`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.32`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.32`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.32` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.32`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.32`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.32`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.32`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.32`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.32)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.32/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.32/charts/vald-helm-operator/README.md)

### Changes

- bugfix ip discoverer disconnection too slow ([#344](https://github.com/vdaas/vald/pull/344))
- Compressor: backup vectors in queue using PostStop function ([#345](https://github.com/vdaas/vald/pull/345))
- Revise backup/meta Cassandra default values ([#336](https://github.com/vdaas/vald/pull/336))

## v0.0.31

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.31`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.31`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.31`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.31`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.31` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.31`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.31`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.31`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.31`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.31`            |

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.31)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.31/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.31/charts/vald-helm-operator/README.md)

### Changes

- Resolve busy-loop on worker ([#339](https://github.com/vdaas/vald/pull/339))

## v0.0.30

### Docker images

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.30`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.30`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.30`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.30`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.30` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.30`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.30`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.30`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.30`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.30`            |

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

| component                | docker pull                                               |
| ------------------------ | --------------------------------------------------------- |
| agent NGT                | `docker pull vdaas/vald-agent-ngt:v0.0.26`                |
| discoverer K8s           | `docker pull vdaas/vald-discoverer-k8s:v0.0.26`           |
| gateway                  | `docker pull vdaas/vald-gateway:v0.0.26`                  |
| backup manager MySQL     | `docker pull vdaas/vald-manager-backup-mysql:v0.0.26`     |
| backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.26` |
| compressor               | `docker pull vdaas/vald-manager-compressor:v0.0.26`       |
| meta Redis               | `docker pull vdaas/vald-meta-redis:v0.0.26`               |
| meta Cassandra           | `docker pull vdaas/vald-meta-cassandra:v0.0.26`           |
| index manager            | `docker pull vdaas/vald-manager-index:v0.0.26`            |
| Helm operator            | `docker pull vdaas/vald-helm-operator:v0.0.26`            |

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
