# CHANGELOG v1.2.x

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


