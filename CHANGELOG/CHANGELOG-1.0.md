# CHANGELOG v1.0.x

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
- bugfix add nil check for gRPC connection pool objects in grpc/client.go ([#921](https://github.com/vdaas/vald/pull/921))
- remove unneccessary pr-tag definition from chart ([#920](https://github.com/vdaas/vald/pull/920))
- :pencil: Fix typo in gateway-vald configmap template ([#919](https://github.com/vdaas/vald/pull/919))
- change docker base image PRIMARY_TAG name from nightly to latest ([#917](https://github.com/vdaas/vald/pull/917))
- :green_heart: Use vdaas-ci token for making commits ([#895](https://github.com/vdaas/vald/pull/895))
- Vald V1 New Design APIs ([#826](https://github.com/vdaas/vald/pull/826))
- :robot: Automatically update k8s manifests ([#914](https://github.com/vdaas/vald/pull/914))
