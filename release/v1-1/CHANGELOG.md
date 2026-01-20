# CHANGELOG v1.1.x

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
- bugfix remove unneccessary error return & add gRPC status code handling for backoff ([#1253](https://github.com/vdaas/vald/pull/1253))
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

- add correct gRPC status and error handling (#1224)
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

