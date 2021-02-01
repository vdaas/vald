# CHANGELOG

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
- :green_heart: Add formatter for master branch ([#911](https://github.com/vdaas/vald/pull/911))
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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.59`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.59`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.59`
gateway | `docker pull vdaas/vald-gateway:v0.0.59`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.59`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.59`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.59`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.59`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.59`
index manager | `docker pull vdaas/vald-manager-index:v0.0.59`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.59`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.59)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.59/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.59/charts/vald-helm-operator/README.md)

### Changes
- bugfix gateway index out of bounds ([#768](https://github.com/vdaas/vald/pull/768))
- :robot: Automatically update k8s manifests ([#766](https://github.com/vdaas/vald/pull/766))


## v0.0.58

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.58`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.58`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.58`
gateway | `docker pull vdaas/vald-gateway:v0.0.58`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.58`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.58`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.58`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.58`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.58`
index manager | `docker pull vdaas/vald-manager-index:v0.0.58`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.58`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.57`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.57`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.57`
gateway | `docker pull vdaas/vald-gateway:v0.0.57`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.57`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.57`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.57`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.57`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.57`
index manager | `docker pull vdaas/vald-manager-index:v0.0.57`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.57`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.56`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.56`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.56`
gateway | `docker pull vdaas/vald-gateway:v0.0.56`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.56`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.56`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.56`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.56`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.56`
index manager | `docker pull vdaas/vald-manager-index:v0.0.56`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.56`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.55`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.55`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.55`
gateway | `docker pull vdaas/vald-gateway:v0.0.55`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.55`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.55`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.55`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.55`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.55`
index manager | `docker pull vdaas/vald-manager-index:v0.0.55`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.55`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.55)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.55/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.55/charts/vald-helm-operator/README.md)

### Changes
- pass CFLAGS, CXXFLAGS to NGT build command ([#683](https://github.com/vdaas/vald/pull/683))
- :robot: Automatically update k8s manifests ([#681](https://github.com/vdaas/vald/pull/681))


## v0.0.54

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.54`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.54`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.54`
gateway | `docker pull vdaas/vald-gateway:v0.0.54`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.54`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.54`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.54`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.54`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.54`
index manager | `docker pull vdaas/vald-manager-index:v0.0.54`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.54`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.54)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.54/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.54/charts/vald-helm-operator/README.md)

### Changes
- bugfix error assertion ([#680](https://github.com/vdaas/vald/pull/680))
- :robot: Automatically update k8s manifests ([#679](https://github.com/vdaas/vald/pull/679))


## v0.0.53

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.53`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.53`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.53`
gateway | `docker pull vdaas/vald-gateway:v0.0.53`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.53`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.53`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.53`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.53`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.53`
index manager | `docker pull vdaas/vald-manager-index:v0.0.53`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.53`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.52`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.52`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.52`
gateway | `docker pull vdaas/vald-gateway:v0.0.52`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.52`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.52`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.52`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.52`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.52`
index manager | `docker pull vdaas/vald-manager-index:v0.0.52`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.52`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.51`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.51`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.51`
gateway | `docker pull vdaas/vald-gateway:v0.0.51`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.51`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.51`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.51`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.51`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.51`
index manager | `docker pull vdaas/vald-manager-index:v0.0.51`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.51`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.50`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.50`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.50`
gateway | `docker pull vdaas/vald-gateway:v0.0.50`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.50`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.50`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.50`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.50`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.50`
index manager | `docker pull vdaas/vald-manager-index:v0.0.50`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.50`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.50)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.50/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.50/charts/vald-helm-operator/README.md)

### Changes
- Add warn logging messages to agent-sidecar &  ignore io.EOF error when reading metadata.json ([#608](https://github.com/vdaas/vald/pull/608))
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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.49`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.49`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.49`
gateway | `docker pull vdaas/vald-gateway:v0.0.49`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.49`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.49`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.49`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.49`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.49`
index manager | `docker pull vdaas/vald-manager-index:v0.0.49`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.49`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.49)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.49/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.49/charts/vald-helm-operator/README.md)

### Changes
- :bug: fix agent sidecar behavior ([#598](https://github.com/vdaas/vald/pull/598))
- :robot: Automatically update k8s manifests ([#597](https://github.com/vdaas/vald/pull/597))


## v0.0.48

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.48`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.48`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.48`
gateway | `docker pull vdaas/vald-gateway:v0.0.48`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.48`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.48`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.48`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.48`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.48`
index manager | `docker pull vdaas/vald-manager-index:v0.0.48`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.48`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.47`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.47`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.47`
gateway | `docker pull vdaas/vald-gateway:v0.0.47`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.47`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.47`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.47`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.47`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.47`
index manager | `docker pull vdaas/vald-manager-index:v0.0.47`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.47`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.46`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.46`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.46`
gateway | `docker pull vdaas/vald-gateway:v0.0.46`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.46`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.46`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.46`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.46`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.46`
index manager | `docker pull vdaas/vald-manager-index:v0.0.46`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.46`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.45`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.45`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.45`
gateway | `docker pull vdaas/vald-gateway:v0.0.45`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.45`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.45`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.45`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.45`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.45`
index manager | `docker pull vdaas/vald-manager-index:v0.0.45`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.45`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.44`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.44`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.44`
gateway | `docker pull vdaas/vald-gateway:v0.0.44`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.44`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.44`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.44`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.44`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.44`
index manager | `docker pull vdaas/vald-manager-index:v0.0.44`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.44`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.43`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.43`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.43`
gateway | `docker pull vdaas/vald-gateway:v0.0.43`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.43`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.43`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.43`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.43`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.43`
index manager | `docker pull vdaas/vald-manager-index:v0.0.43`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.43`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.43)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.43/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.43/charts/vald-helm-operator/README.md)

### Changes
- Revise S3 reader/writer: compatible with IBM Cloud Object Storage ([#509](https://github.com/vdaas/vald/pull/509))
- :bug: Close [#502](https://github.com/vdaas/vald/pull/502) / Fix roundtrip error handling (#508)
- Feature/drawio ([#500](https://github.com/vdaas/vald/pull/500))
- Added test case for `internal/errorgroup`  ([#494](https://github.com/vdaas/vald/pull/494))
- Update Helm Chart info ([#496](https://github.com/vdaas/vald/pull/496))
- Revise triggers of workflow run & Fix reading changelogs from PR comments ([#495](https://github.com/vdaas/vald/pull/495))


## v0.0.42

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.42`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.42`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.42`
gateway | `docker pull vdaas/vald-gateway:v0.0.42`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.42`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.42`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.42`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.42`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.42`
index manager | `docker pull vdaas/vald-manager-index:v0.0.42`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.42`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.41`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.41`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.41`
gateway | `docker pull vdaas/vald-gateway:v0.0.41`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.41`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.41`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.41`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.41`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.41`
index manager | `docker pull vdaas/vald-manager-index:v0.0.41`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.41`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.40`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.40`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.40`
gateway | `docker pull vdaas/vald-gateway:v0.0.40`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.40`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.40`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.40`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.40`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.40`
index manager | `docker pull vdaas/vald-manager-index:v0.0.40`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.40`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.39`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.39`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.39`
gateway | `docker pull vdaas/vald-gateway:v0.0.39`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.39`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.39`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.39`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.39`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.39`
index manager | `docker pull vdaas/vald-manager-index:v0.0.39`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.39`

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
- 🐛  Add gRPC interceptor to recover panic in handlers ([#440](https://github.com/vdaas/vald/pull/440))
- tensorflow test ([#378](https://github.com/vdaas/vald/pull/378))
- :bento: update architecture overview svg to add agent sidecar ([#437](https://github.com/vdaas/vald/pull/437))
- Example program: Add indexing interval description & fix logging message ([#405](https://github.com/vdaas/vald/pull/405))
- :pencil2: Fix typo ([#436](https://github.com/vdaas/vald/pull/436))


## v0.0.38

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.38`
agent sidecar | `docker pull vdaas/vald-agent-sidecar:v0.0.38`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.38`
gateway | `docker pull vdaas/vald-gateway:v0.0.38`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.38`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.38`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.38`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.38`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.38`
index manager | `docker pull vdaas/vald-manager-index:v0.0.38`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.38`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.37`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.37`
gateway | `docker pull vdaas/vald-gateway:v0.0.37`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.37`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.37`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.37`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.37`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.37`
index manager | `docker pull vdaas/vald-manager-index:v0.0.37`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.37`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.36`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.36`
gateway | `docker pull vdaas/vald-gateway:v0.0.36`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.36`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.36`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.36`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.36`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.36`
index manager | `docker pull vdaas/vald-manager-index:v0.0.36`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.36`

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
- Fix failing tests of `internal/log` and modified to match the test template  ([#368](https://github.com/vdaas/vald/pull/368))
- Add enabled flag to each component in Helm chart ([#372](https://github.com/vdaas/vald/pull/372))
- Add configurations.md ([#356](https://github.com/vdaas/vald/pull/356))


## v0.0.35

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.35`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.35`
gateway | `docker pull vdaas/vald-gateway:v0.0.35`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.35`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.35`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.35`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.35`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.35`
index manager | `docker pull vdaas/vald-manager-index:v0.0.35`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.35`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.34`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.34`
gateway | `docker pull vdaas/vald-gateway:v0.0.34`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.34`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.34`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.34`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.34`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.34`
index manager | `docker pull vdaas/vald-manager-index:v0.0.34`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.34`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.33`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.33`
gateway | `docker pull vdaas/vald-gateway:v0.0.33`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.33`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.33`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.33`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.33`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.33`
index manager | `docker pull vdaas/vald-manager-index:v0.0.33`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.33`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.32`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.32`
gateway | `docker pull vdaas/vald-gateway:v0.0.32`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.32`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.32`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.32`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.32`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.32`
index manager | `docker pull vdaas/vald-manager-index:v0.0.32`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.32`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.31`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.31`
gateway | `docker pull vdaas/vald-gateway:v0.0.31`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.31`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.31`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.31`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.31`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.31`
index manager | `docker pull vdaas/vald-manager-index:v0.0.31`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.31`

### Documents
- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v0.0.31)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v0.0.31/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v0.0.31/charts/vald-helm-operator/README.md)

### Changes
- Resolve busy-loop on worker ([#339](https://github.com/vdaas/vald/pull/339))


## v0.0.30

### Docker images

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.30`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.30`
gateway | `docker pull vdaas/vald-gateway:v0.0.30`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.30`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.30`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.30`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.30`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.30`
index manager | `docker pull vdaas/vald-manager-index:v0.0.30`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.30`

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

component | docker pull
--------- | -----------
agent NGT | `docker pull vdaas/vald-agent-ngt:v0.0.26`
discoverer K8s | `docker pull vdaas/vald-discoverer-k8s:v0.0.26`
gateway | `docker pull vdaas/vald-gateway:v0.0.26`
backup manager MySQL | `docker pull vdaas/vald-manager-backup-mysql:v0.0.26`
backup manager Cassandra | `docker pull vdaas/vald-manager-backup-cassandra:v0.0.26`
compressor | `docker pull vdaas/vald-manager-compressor:v0.0.26`
meta Redis | `docker pull vdaas/vald-meta-redis:v0.0.26`
meta Cassandra | `docker pull vdaas/vald-meta-cassandra:v0.0.26`
index manager | `docker pull vdaas/vald-manager-index:v0.0.26`
Helm operator | `docker pull vdaas/vald-helm-operator:v0.0.26`

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
