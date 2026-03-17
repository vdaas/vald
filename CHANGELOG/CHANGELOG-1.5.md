# CHANGELOG v1.5.x

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

