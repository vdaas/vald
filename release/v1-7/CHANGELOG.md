# CHANGELOG v1.7.x

## v1.7.17

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.17</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.17</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.17</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.17</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.17</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.17</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.17</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.17</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.17)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.17/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.17/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature
[VALD-125] feat: Implement index exporter [#2746](https://github.com/vdaas/vald/pull/2746)

[VALD-147] add server_config [#2958](https://github.com/vdaas/vald/pull/2958)

[VALD-147] implement AccessLogMiddleware [#2991](https://github.com/vdaas/vald/pull/2991)

[VALD-148] Generate grafana boards by grafana-foundation-sdk Go [#2937](https://github.com/vdaas/vald/pull/2937)

[VALD-322] e2e/v2: Support selector for wait action [#2956](https://github.com/vdaas/vald/pull/2956)

[VALD-325] E2E V2: Index Correction Job [#3000](https://github.com/vdaas/vald/pull/3000)

[VALD-336] Implement expect in E2E v2 to assert API results [#2971](https://github.com/vdaas/vald/pull/2971)

[VALD-351] Add Rust version VQueue prototype implementation [#2998](https://github.com/vdaas/vald/pull/2998)

Add git-lfs to dev container image [#2978](https://github.com/vdaas/vald/pull/2978)

Fix dashboard variables in overview board (#3030) (#3042)

[VALD-359] Support Mac C flags (#3040) (#3044)

Support prometheus without container label (#3043) (#3045)

:zap: Improve performance

improve rangeConns's performance for Release v1.7.17 (#3039) (#3041)

:bug: Bugfix
[BUGFIX] use client config for tls dialer [#2938](https://github.com/vdaas/vald/pull/2938)

[Bugfix] prevent nil pointer panic for internal/config/grpc.go [#3017](https://github.com/vdaas/vald/pull/3017)

[VALD-148] hotfix: Fix name & despcription of metrics [#2981](https://github.com/vdaas/vald/pull/2981)

Avoid concurrent assignment of stream client in RoundRobin [#2954](https://github.com/vdaas/vald/pull/2954)

fix: make format has conflicts in some targets [#2977](https://github.com/vdaas/vald/pull/2977)

:pencil2: Document
docs: Fix a typo in the value that specified the service name [#3005](https://github.com/vdaas/vald/pull/3005)

docs: add Kynea0b as a contributor for doc [#3006](https://github.com/vdaas/vald/pull/3006)

docs: add Matts966 as a contributor for code, infra, and 2 more [#3007](https://github.com/vdaas/vald/pull/3007)

:white_check_mark: Testing
E2E V2 CI [#2950](https://github.com/vdaas/vald/pull/2950)

:green_heart: CI
feat: automatically resolve go.mod & go.sum conflicts in backport workflow [#2980](https://github.com/vdaas/vald/pull/2980)

:chart_with_upwards_trend: Metrics/Tracing
[VALD-148] Generate grafana boards by grafana-foundation-sdk Go [#2937](https://github.com/vdaas/vald/pull/2937)

## v1.7.16

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.16</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.16</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.16</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.16</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.16</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.16</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.16</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.16</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.16)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.16/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.16/charts/vald-helm-operator/README.md)

### Changes

:recycle: Refactor

- Fix format of proto files [#2778](https://github.com/vdaas/vald/pull/2778) ([#2783](https://github.com/vdaas/vald/pull/2783))
- Refactor merge docker and github actions workflow gen logic [#2769](https://github.com/vdaas/vald/pull/2769) ([#2774](https://github.com/vdaas/vald/pull/2774))

:pencil2: Document

- Change symlink API documents [#2741](https://github.com/vdaas/vald/pull/2741) ([#2776](https://github.com/vdaas/vald/pull/2776))

:green_heart: CI

- Refactor github actions [#2773](https://github.com/vdaas/vald/pull/2773) ([#2779](https://github.com/vdaas/vald/pull/2779))
  Change make command [#2765](https://github.com/vdaas/vald/pull/2765) ([#2770](https://github.com/vdaas/vald/pull/2770))

:arrow_up: Update dependencies

- Update libs dependency [#2775](https://github.com/vdaas/vald/pull/2775) ([#2785](https://github.com/vdaas/vald/pull/2785))

## v1.7.15

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.15</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.15</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.15</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.15</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.15</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.15</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.15</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.15</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.15)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.15/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.15/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Add affinity to jobTemplate [#2758](https://github.com/vdaas/vald/pull/2758) ([#2760](https://github.com/vdaas/vald/pull/2760))
- feat: Implement delete expired index job [#2702](https://github.com/vdaas/vald/pull/2702) ([#2722](https://github.com/vdaas/vald/pull/2722))
- Add QUIC support [#1771](https://github.com/vdaas/vald/pull/1771)
- add example-client docker image [#2705](https://github.com/vdaas/vald/pull/2705) ([#2709](https://github.com/vdaas/vald/pull/2709))

:recycle: Refactor

- refactor dockerfiles and update gitattributes [#2743](https://github.com/vdaas/vald/pull/2743) ([#2745](https://github.com/vdaas/vald/pull/2745))

:bug: Bugfix

- :bug: Fix update deps workflow: buf is not found [#2737](https://github.com/vdaas/vald/pull/2737) ([#2739](https://github.com/vdaas/vald/pull/2739))
- [BUGFIX] resolve agent GetGraphStatistics API double-free error problem [#2733](https://github.com/vdaas/vald/pull/2733)
- fix rust-analyzer [#2731](https://github.com/vdaas/vald/pull/2731) ([#2732](https://github.com/vdaas/vald/pull/2732))
- Fix installation command for arm64 [#2729](https://github.com/vdaas/vald/pull/2729) ([#2730](https://github.com/vdaas/vald/pull/2730))
- fix not found error [#2726](https://github.com/vdaas/vald/pull/2726) ([#2727](https://github.com/vdaas/vald/pull/2727))
- Fix bind DOCKER_OPTS option [#2718](https://github.com/vdaas/vald/pull/2718) ([#2719](https://github.com/vdaas/vald/pull/2719))

:pencil2: Document

- Update README.md [#2724](https://github.com/vdaas/vald/pull/2724) ([#2725](https://github.com/vdaas/vald/pull/2725))
- :pencil: Remove clj link [#2710](https://github.com/vdaas/vald/pull/2710) ([#2714](https://github.com/vdaas/vald/pull/2714))

:green_heart: CI

- :green_heart: Multi-PF build for example-client [#2713](https://github.com/vdaas/vald/pull/2713)
- Add auto deps version update workflow [#2707](https://github.com/vdaas/vald/pull/2707) ([#2717](https://github.com/vdaas/vald/pull/2717))

:arrow_up: Update dependencies

- :green_heart: use ci-container for update deps cron job [#2744](https://github.com/vdaas/vald/pull/2744) ([#2748](https://github.com/vdaas/vald/pull/2748))
- update ubuntu version for devcontainer [#2736](https://github.com/vdaas/vald/pull/2736) ([#2750](https://github.com/vdaas/vald/pull/2750))
- :arrow_up: update versions/BUF_VERSION [#2703](https://github.com/vdaas/vald/pull/2703) ([#2704](https://github.com/vdaas/vald/pull/2704))

:handshake: Contributor

- docs: add highpon as a contributor for code [#2721](https://github.com/vdaas/vald/pull/2721)

## v1.7.14

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.14</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.14</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.14</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.14</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.14</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.14</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.14</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.14</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.14)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.14/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.14/charts/vald-helm-operator/README.md)

### Changes

### :sparkles: New feature

- Add String sorted topologicalSort [#2696](https://github.com/vdaas/vald/pull/2696) [#2698](https://github.com/vdaas/vald/pull/2698)
- Add CPU_INFO_FLAGS for Apple Silicon [#2694](https://github.com/vdaas/vald/pull/2694) [#2697](https://github.com/vdaas/vald/pull/2697)
- Add New gRPC Options and Add Reconnect Logic for connection Pool [#2685](https://github.com/vdaas/vald/pull/2685) [#2693](https://github.com/vdaas/vald/pull/2693)
- Add option to disable dns resolve [#2634](https://github.com/vdaas/vald/pull/2634) [#2641](https://github.com/vdaas/vald/pull/2641)
- Backport PR #2584 to release/v1.7 for Implement ngt property get API [#2588](https://github.com/vdaas/vald/pull/2588)
- add HTTP2 support for http.Client and Vald HTTP Server [#2572](https://github.com/vdaas/vald/pull/2572) [#2575](https://github.com/vdaas/vald/pull/2575)

### :zap: Improve performance

- Refactor grpc/status.withDetails function for performance [#2664](https://github.com/vdaas/vald/pull/2664) [#2668](https://github.com/vdaas/vald/pull/2668)

### :recycle: Refactor

- Refactor use Absolute path for Makefile [#2673](https://github.com/vdaas/vald/pull/2673)
- Refactor internal/net/grpc/client.go [#2675](https://github.com/vdaas/vald/pull/2675)
- modify ParseError to FromError for agent handler [#2667](https://github.com/vdaas/vald/pull/2667) [#2679](https://github.com/vdaas/vald/pull/2679)
- Backport PR #2674 to release/v1.7 for Refactor internal/net/grpc/client.go [#2675](https://github.com/vdaas/vald/pull/2675)
- Backport PR #2670 to release/v1.7 for Refactor use Absolute path for Makefile [#2673](https://github.com/vdaas/vald/pull/2673)
- Refactor grpc/status.withDetails function for performance [#2664](https://github.com/vdaas/vald/pull/2664) [#2668](https://github.com/vdaas/vald/pull/2668)
- Refactor for release v1.7.14 [#2639](https://github.com/vdaas/vald/pull/2639) [#2648](https://github.com/vdaas/vald/pull/2648)
- refactor(gateway): delete unused file [#2644](https://github.com/vdaas/vald/pull/2644) [#2646](https://github.com/vdaas/vald/pull/2646)
- Refactor test checkFunc condition [#2599](https://github.com/vdaas/vald/pull/2599) [#2602](https://github.com/vdaas/vald/pull/2602)
- Backport PR #2586 to release/v1.7 for modify rust package structure [#2590](https://github.com/vdaas/vald/pull/2590)
- Backport PR #2577 to release/v1.7 for refactor docker and change buildkit-syft-scanner reference to ghcr.io [#2578](https://github.com/vdaas/vald/pull/2578)

### :bug: Bugfix

- Fix gRPC error handling for gateway/filter handler [#2669](https://github.com/vdaas/vald/pull/2669) [#2689](https://github.com/vdaas/vald/pull/2689)
- fix: increase limit [#2683](https://github.com/vdaas/vald/pull/2683) [#2686](https://github.com/vdaas/vald/pull/2686)
- Fix gRPC error handling for mirror-gateway handler [#2665](https://github.com/vdaas/vald/pull/2665) [#2681](https://github.com/vdaas/vald/pull/2681)
- Fix gRPC error msg handling for lb-gateway handler [#2663](https://github.com/vdaas/vald/pull/2663) [#2682](https://github.com/vdaas/vald/pull/2682)
- Bugfix ingress route settings [#2636](https://github.com/vdaas/vald/pull/2636) [#2642](https://github.com/vdaas/vald/pull/2642)
- Fix broken links in the document files [#2611](https://github.com/vdaas/vald/pull/2611) [#2614](https://github.com/vdaas/vald/pull/2614)
- Fix: make command name [#2610](https://github.com/vdaas/vald/pull/2610) [#2612](https://github.com/vdaas/vald/pull/2612)
- Bugfix NGT flush logic [#2598](https://github.com/vdaas/vald/pull/2598) [#2606](https://github.com/vdaas/vald/pull/2606)

### :pencil2: Document

- Fix broken links in the document files [#2611](https://github.com/vdaas/vald/pull/2611) [#2614](https://github.com/vdaas/vald/pull/2614)

### :white_check_mark: Testing

- Refactor test checkFunc condition [#2599](https://github.com/vdaas/vald/pull/2599) [#2602](https://github.com/vdaas/vald/pull/2602)

### :green_heart: CI

- Buf CLI migrate to v2 [#2691](https://github.com/vdaas/vald/pull/2691) [#2695](https://github.com/vdaas/vald/pull/2695)
- [create-pull-request] automated change [#2677](https://github.com/vdaas/vald/pull/2677) [#2678](https://github.com/vdaas/vald/pull/2678)
- automatically generate workflows [#2595](https://github.com/vdaas/vald/pull/2595) [#2603](https://github.com/vdaas/vald/pull/2603)

### :chart_with_upwards_trend: Metrics/Tracing

- Introduce an observability crate using opentelemetry-rust [#2535](https://github.com/vdaas/vald/pull/2535) [#2609](https://github.com/vdaas/vald/pull/2609)

<!-- This is an auto-generated comment: release notes by coderabbit.ai -->

## Summary by CodeRabbit

- **New Features**
  - Added several new contributors to the project, enhancing community involvement.
  - Introduced a new configuration file for spell checking, improving documentation quality.
  - Expanded the project with new configuration files, documentation, and source code for enhanced functionality.

- **Bug Fixes**
  - Updated version information in issue templates for accuracy.

- **Documentation**
  - Improved clarity in the pull request template and updated version information.

- **Chores**
  - Modified GitHub Actions for better handling of Docker image tags.

<!-- end of auto-generated comment: release notes by coderabbit.ai -->

## v1.7.13

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.13</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.13</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.13</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.13</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.13</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.13</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.13</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.13</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.13)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.13/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.13/charts/vald-helm-operator/README.md)

### Changes

### [Bugfix]

- **General Fixes**
  1. Fix index correction process [#2565](https://github.com/vdaas/vald/pull/2565) ([#2566](https://github.com/vdaas/vald/pull/2566))
  2. libquadmath is not available on ARM [#2559](https://github.com/vdaas/vald/pull/2559)
  3. fix: add checkout option [#2545](https://github.com/vdaas/vald/pull/2545) ([#2546](https://github.com/vdaas/vald/pull/2546))
  4. fix: make format [#2534](https://github.com/vdaas/vald/pull/2534) ([#2540](https://github.com/vdaas/vald/pull/2540))
  5. fix conflict bug [#2537](https://github.com/vdaas/vald/pull/2537)
  6. Bugfix that caused an error when argument has 3 or more nil arguments [#2517](https://github.com/vdaas/vald/pull/2517) ([#2520](https://github.com/vdaas/vald/pull/2520))
  7. Bugfix recreate benchmark job when operator reboot [#2463](https://github.com/vdaas/vald/pull/2463) ([#2464](https://github.com/vdaas/vald/pull/2464))
  8. Fix agent-faiss build failed [#2418](https://github.com/vdaas/vald/pull/2418) ([#2419](https://github.com/vdaas/vald/pull/2419))
  9. Fix the logic to determine docker image [#2410](https://github.com/vdaas/vald/pull/2410) ([#2420](https://github.com/vdaas/vald/pull/2420))

- **Backport and Release-Related**
  1. Fix workflow trigger for backport pr creation [#2471](https://github.com/vdaas/vald/pull/2471) ([#2472](https://github.com/vdaas/vald/pull/2472))
  2. Fix output settings to determine-docker-image-tag action and release branch build tag name [#2423](https://github.com/vdaas/vald/pull/2423) ([#2425](https://github.com/vdaas/vald/pull/2425))

- **E2E and Index**
  1. Fix e2e for read replica and add e2e for index operator [#2455](https://github.com/vdaas/vald/pull/2455) ([#2459](https://github.com/vdaas/vald/pull/2459))
  2. Fix index job logic to pass DNS A record [#2438](https://github.com/vdaas/vald/pull/2438) ([#2448](https://github.com/vdaas/vald/pull/2448))

- **Documentation and Other**
  1. fix: typo of execution rule [#2426](https://github.com/vdaas/vald/pull/2426) ([#2427](https://github.com/vdaas/vald/pull/2427))
  2. :pencil: Fix typo of file name [#2413](https://github.com/vdaas/vald/pull/2413) ([#2415](https://github.com/vdaas/vald/pull/2415))

### [Enhancement]

- **General Improvements**
  1. Update dependencies, C++ standard, and improve Dockerfiles for better build systems and localization [#2549](https://github.com/vdaas/vald/pull/2549) ([#2557](https://github.com/vdaas/vald/pull/2557))
  2. Implement ngt Statistics API [#2539](https://github.com/vdaas/vald/pull/2539) ([#2547](https://github.com/vdaas/vald/pull/2547))
  3. refactor index manager service add index service API to expose index informations [#2525](https://github.com/vdaas/vald/pull/2525) ([#2532](https://github.com/vdaas/vald/pull/2532))

- **API and Logic Changes**
  1. Change default image tag from latest to nightly [#2516](https://github.com/vdaas/vald/pull/2516) ([#2518](https://github.com/vdaas/vald/pull/2518))
  2. update large top-K ratio handling logic [#2509](https://github.com/vdaas/vald/pull/2509) ([#2511](https://github.com/vdaas/vald/pull/2511))
  3. add inner product distance type for ngt [#2454](https://github.com/vdaas/vald/pull/2454) ([#2458](https://github.com/vdaas/vald/pull/2458))
  4. Implement index operator logic for read replica rotation [#2444](https://github.com/vdaas/vald/pull/2444) ([#2456](https://github.com/vdaas/vald/pull/2456))

- **Performance and Optimization**
  1. update deps & add validation for Flush API when agent is Read Only [#2433](https://github.com/vdaas/vald/pull/2433) ([#2436](https://github.com/vdaas/vald/pull/2436))
  2. Add `index-operator` template implementation [#2375](https://github.com/vdaas/vald/pull/2375) ([#2424](https://github.com/vdaas/vald/pull/2424))

- **Testing and Metrics**
  1. Implement client metrics interceptor for continuous benchmark job [#2477](https://github.com/vdaas/vald/pull/2477) ([#2480](https://github.com/vdaas/vald/pull/2480))
  2. Add tests for index information export [#2412](https://github.com/vdaas/vald/pull/2412) ([#2414](https://github.com/vdaas/vald/pull/2414))

### [CI]

1. [create-pull-request] automated change [#2552](https://github.com/vdaas/vald/pull/2552) ([#2556](https://github.com/vdaas/vald/pull/2556))
2. Add workflow to check git conflict for backport PR [#2548](https://github.com/vdaas/vald/pull/2548) ([#2550](https://github.com/vdaas/vald/pull/2550))
3. [CI] Add workflow to synchronize ubuntu base image [#2526](https://github.com/vdaas/vald/pull/2526) ([#2527](https://github.com/vdaas/vald/pull/2527))
4. Automatically add backport main label for release-pr [#2473](https://github.com/vdaas/vald/pull/2473) ([#2475](https://github.com/vdaas/vald/pull/2475))
5. change external docker image reference to ghcr.io registry [#2567](https://github.com/vdaas/vald/pull/2567) ([#2568](https://github.com/vdaas/vald/pull/2568))

### [Backport]

1. Backport PR #2542, #2538 to release/v1.7 [#2543](https://github.com/vdaas/vald/pull/2543)
2. Backport docs updates to release/v1.7 [#2521](https://github.com/vdaas/vald/pull/2521)
3. Backport Flush API [#2434](https://github.com/vdaas/vald/pull/2434)

### [Documentation]

1. capitalize faq [#2512](https://github.com/vdaas/vald/pull/2512) ([#2522](https://github.com/vdaas/vald/pull/2522))
2. add faiss in values.yaml & valdrelease.yaml [#2514](https://github.com/vdaas/vald/pull/2514) ([#2519](https://github.com/vdaas/vald/pull/2519))
3. add read replica and rotator docs [#2497](https://github.com/vdaas/vald/pull/2497) ([#2499](https://github.com/vdaas/vald/pull/2499))
4. Update continuous benchmark docs [#2485](https://github.com/vdaas/vald/pull/2485) ([#2486](https://github.com/vdaas/vald/pull/2486))
5. docs: add hrichiksite as a contributor for doc [#2441](https://github.com/vdaas/vald/pull/2441) ([#2442](https://github.com/vdaas/vald/pull/2442))

### [Other]

1. Add base of benchmark operator dashboard [#2430](https://github.com/vdaas/vald/pull/2430) ([#2453](https://github.com/vdaas/vald/pull/2453))
2. Add client metrics panels for continuous benchmark job [#2481](https://github.com/vdaas/vald/pull/2481) ([#2483](https://github.com/vdaas/vald/pull/2483))
3. Add unit tests for index operator [#2460](https://github.com/vdaas/vald/pull/2460) ([#2461](https://github.com/vdaas/vald/pull/2461))
4. add reviewer guideline [#2507](https://github.com/vdaas/vald/pull/2507) ([#2508](https://github.com/vdaas/vald/pull/2508))
5. Sync release/v1.7 to main [#2495](https://github.com/vdaas/vald/pull/2495)
6. Add snapshot timestamp annotations to read replica agent [#2428](https://github.com/vdaas/vald/pull/2428) ([#2443](https://github.com/vdaas/vald/pull/2443))
7. Update build rule for nightly image [#2421](https://github.com/vdaas/vald/pull/2421) ([#2422](https://github.com/vdaas/vald/pull/2422))

## v1.7.12

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.12</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.12</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.12</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.12</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.12</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.12</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.12</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.12</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.12)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.12/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.12/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Add rotate-all option to rotator [#2305](https://github.com/vdaas/vald/pull/2305) [#2393](https://github.com/vdaas/vald/pull/2393)
- Make agent export index metrics to Pod k8s resource [#2319](https://github.com/vdaas/vald/pull/2319) [#2372](https://github.com/vdaas/vald/pull/2372)

:recycle: Refactor

- Delete unnecessary code for mirror [#2366](https://github.com/vdaas/vald/pull/2366) [#2391](https://github.com/vdaas/vald/pull/2391)

:bug: Bugfix

- Resolve kvs already closed before last saving [#2390](https://github.com/vdaas/vald/pull/2390) [#2394](https://github.com/vdaas/vald/pull/2394)

:pencil2: Document

- Create continous benchmark doc [#2352](https://github.com/vdaas/vald/pull/2352) [#2395](https://github.com/vdaas/vald/pull/2395)

:white_check_mark: Testing

- Fix: build error of internal kvs test [#2396](https://github.com/vdaas/vald/pull/2396) [#2398](https://github.com/vdaas/vald/pull/2398)

:green_heart: CI

- Fix: disable protobuf dispatch for client [#2401](https://github.com/vdaas/vald/pull/2401) [#2403](https://github.com/vdaas/vald/pull/2403)
- Add Con-Bench helm chart to the Vald charts [#2388](https://github.com/vdaas/vald/pull/2388) [#2389](https://github.com/vdaas/vald/pull/2389)
- Update workflow to release readreplica chart [#2383](https://github.com/vdaas/vald/pull/2383) [#2387](https://github.com/vdaas/vald/pull/2387)
- Backport ci deps others [#2386](https://github.com/vdaas/vald/pull/2386)
- Update docker build target platform selection rules [#2370](https://github.com/vdaas/vald/pull/2370) [#2374](https://github.com/vdaas/vald/pull/2374)
- Add commit hash build image [#2359](https://github.com/vdaas/vald/pull/2359) [#2371](https://github.com/vdaas/vald/pull/2371)
- Refactor code using golangci-lint [#2362](https://github.com/vdaas/vald/pull/2362) [#2365](https://github.com/vdaas/vald/pull/2365)
- Change docker scan timeout longer [#2363](https://github.com/vdaas/vald/pull/2363) [#2364](https://github.com/vdaas/vald/pull/2364)

:arrow_up: Update dependencies

- Update deps [#2404](https://github.com/vdaas/vald/pull/2404) [#2405](https://github.com/vdaas/vald/pull/2405)

:lock: Security

- Create SECURITY.md [#2367](https://github.com/vdaas/vald/pull/2367) [#2368](https://github.com/vdaas/vald/pull/2368)

:art: Design

- Change JP logo to EN logo [#2369](https://github.com/vdaas/vald/pull/2369) [#2392](https://github.com/vdaas/vald/pull/2392)

## v1.7.11

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.11</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.11</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.11</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.11</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.11</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.11</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.11</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.11</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.11)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.11/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.11/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Add HPA for read replica [#2307](https://github.com/vdaas/vald/pull/2307)
- Add internal KVS pogreb package [#2302](https://github.com/vdaas/vald/pull/2302)
- Two version deploy support [#2171](https://github.com/vdaas/vald/pull/2171)
- Add mirror gateway definitions and Implementations [#2262](https://github.com/vdaas/vald/pull/2262)
- Initialize dev env for Rust agent [#2293](https://github.com/vdaas/vald/pull/2293)
- Add new grafana dashboard for agent memory metrics [#2279](https://github.com/vdaas/vald/pull/2279)
- Implement continuous benchmark tool [#2216](https://github.com/vdaas/vald/pull/2216)

:recycle: Refactor

- Add newline between params to avoid false formatting [#2347](https://github.com/vdaas/vald/pull/2347)
- Fix golangci-lint config and apply tagalign [#2326](https://github.com/vdaas/vald/pull/2326)
- Refactor postAttachCommand [#2312](https://github.com/vdaas/vald/pull/2312)
- Refactor ignore rule [#2339](https://github.com/vdaas/vald/pull/2339)
- Fix NGT default params [#2332](https://github.com/vdaas/vald/pull/2332)
- Format yaml using google/yamlfmt & update go version and dependencies [#2322](https://github.com/vdaas/vald/pull/2322)
- Refactor update opentelemetry-go & faiss [#2303](https://github.com/vdaas/vald/pull/2303)
- Change discoverer client to broadcast to read replicas [#2276](https://github.com/vdaas/vald/pull/2276)
- Add stern and telepresence [#2320](https://github.com/vdaas/vald/pull/2320)
- Add issue metrics [#2308](https://github.com/vdaas/vald/pull/2308)
- Add dispatch workflow for update contents of vdaas/web repo [#2294](https://github.com/vdaas/vald/pull/2294)
- Fix: add release build for bench and mirror [#2300](https://github.com/vdaas/vald/pull/2300)
- Fix deeepsource errors [#2299](https://github.com/vdaas/vald/pull/2299)
- Add go cache for improvement docker build performance [#2297](https://github.com/vdaas/vald/pull/2297)
- Add detailed log for readreplica rotator [#2281](https://github.com/vdaas/vald/pull/2281)
- Add isSymlink function and test to gen license to avoid for symlink to become normal file. [#2290](https://github.com/vdaas/vald/pull/2290)
- Add owner reference to the resources made by rotator to delete them when read replica resources are deleted [#2287](https://github.com/vdaas/vald/pull/2287)
- Make vald-readreplica values.yaml to symbolic link [#2286](https://github.com/vdaas/vald/pull/2286)
- Separate readreplica chart [#2283](https://github.com/vdaas/vald/pull/2283)
- Happy New Year 2024 [#2284](https://github.com/vdaas/vald/pull/2284)

:bug: Bugfix

- Fix: disable arm64 [#2354](https://github.com/vdaas/vald/pull/2354)
- gcc environment for ARM [#2334](https://github.com/vdaas/vald/pull/2334)
- Revert dev Dockerfile to use official devcontainer image [#2335](https://github.com/vdaas/vald/pull/2335)
- Revert docker-image.yaml change [#2336](https://github.com/vdaas/vald/pull/2336)
- Fix release pr workflow [#2333](https://github.com/vdaas/vald/pull/2333)
- Fix e2e regressions [#2327](https://github.com/vdaas/vald/pull/2327)
- Bugfix grpc ip direct connection status check [#2316](https://github.com/vdaas/vald/pull/2316)
- Fix k3d connectivity error [#2317](https://github.com/vdaas/vald/pull/2317)
- Change lincense/gen/main.go to skip shebang [#2313](https://github.com/vdaas/vald/pull/2313)
- Stop using ENV ARCH and add --platform in Dockerfile [#2304](https://github.com/vdaas/vald/pull/2304)
- gRPC pool connection health check for DNS Addr may fail during VIP member disconnection [#2277](https://github.com/vdaas/vald/pull/2277)
- Fix isSymlink function to correctly check for symbolic links [#2292](https://github.com/vdaas/vald/pull/2292)
- Disable disconnection during non-IP-direct connection [#2291](https://github.com/vdaas/vald/pull/2291)
- Fix git add chart directory for release [#2356](https://github.com/vdaas/vald/pull/2356)

:pencil2: Document

- Add search optimization document [#2306](https://github.com/vdaas/vald/pull/2306)
- Update capacity planning doc [#2295](https://github.com/vdaas/vald/pull/2295)

:green_heart: CI

- Change to dynamically switch CI container image tag [#2310](https://github.com/vdaas/vald/pull/2310)
- Add E2E tests for read replica feature [#2298](https://github.com/vdaas/vald/pull/2298)
- CI, Docker EXTRA_ARGS not working problem [#2278](https://github.com/vdaas/vald/pull/2278)

## v1.7.10

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.10</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.10</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.10</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.10</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.10</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.10</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.10</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.10</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.10)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.10/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.10/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Implement malloc metrics [#2161](https://github.com/vdaas/vald/pull/2161)
- Add svc discoverer for readreplica svc [#2263](https://github.com/vdaas/vald/pull/2263)
- Add agent readreplica resources [#2258](https://github.com/vdaas/vald/pull/2258)
- Add cronjob for readreplica rotator [#2242](https://github.com/vdaas/vald/pull/2242)

:recycle: Refactor

- Apply make proto/all [#2266](https://github.com/vdaas/vald/pull/2266)
- Migratation to buf [#2236](https://github.com/vdaas/vald/pull/2236)
- Update schema [#2265](https://github.com/vdaas/vald/pull/2265)

:bug: Bugfix

- Resolve duplicated cluster wide resources name problem [#2274](https://github.com/vdaas/vald/pull/2274)

:pencil2: Document

- Add caution sentence for deploy multi-Vald clusters [#2271](https://github.com/vdaas/vald/pull/2271)

:green_heart: CI

- Disable BUILDKIT_INLINE_CACHE on GitHub Actions [#2270](https://github.com/vdaas/vald/pull/2270)
- Fix docker build for scanning [#2269](https://github.com/vdaas/vald/pull/2269)
- change login user and token for ghcr.io & small refactor [#2268](https://github.com/vdaas/vald/pull/2268)
- Add e2e job for index management job [#2239](https://github.com/vdaas/vald/pull/2239)
- Add docker buildx cache [#2261](https://github.com/vdaas/vald/pull/2261)

## v1.7.9

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.9</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.9</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.9</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.9</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.9</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.9</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.9</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.9</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.9)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.9/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.9/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Add read replica rotator [#2241](https://github.com/vdaas/vald/pull/2241)
- Implement manifests for Index Management Job [#2235](https://github.com/vdaas/vald/pull/2235)
- Add job error to report index correction error status [#2231](https://github.com/vdaas/vald/pull/2231)
- Add implementation for save index job [#2227](https://github.com/vdaas/vald/pull/2227)
- Add implementation for create index job [#2223](https://github.com/vdaas/vald/pull/2223)
- Add index correction metrics [#2215](https://github.com/vdaas/vald/pull/2215)
- Add index correction document [#2217](https://github.com/vdaas/vald/pull/2217)
- Add make command to update template [#2212](https://github.com/vdaas/vald/pull/2212)
- Add job to check format difference [#2214](https://github.com/vdaas/vald/pull/2214)
- Add verification for index correction e2e and add clusterrole cronjobs for operator to deploy index correction [#2205](https://github.com/vdaas/vald/pull/2205)
- Add StreamListObject to LB [#2203](https://github.com/vdaas/vald/pull/2203)
- Add index correction helm templates and E2E [#2200](https://github.com/vdaas/vald/pull/2200)
- Add index correction internal logic [#2194](https://github.com/vdaas/vald/pull/2194)
- Add bbolt as internal/db/kvs [#2177](https://github.com/vdaas/vald/pull/2177)

:zap: Improve performance

- Improve index correction performance [#2234](https://github.com/vdaas/vald/pull/2234)

:recycle: Refactor

- Refactor Index Management Job [#2232](https://github.com/vdaas/vald/pull/2232)
- Fix invalid network policy schema [#2230](https://github.com/vdaas/vald/pull/2230)
- Add minikube to create volume snapshot development environment locally [#2228](https://github.com/vdaas/vald/pull/2228)
- Enable ingress resource in the get started document [#2211](https://github.com/vdaas/vald/pull/2211)
- Add step to get k3s latest version [#2206](https://github.com/vdaas/vald/pull/2206)
- Update telepresence and helm-docs installer and update deps [#2195](https://github.com/vdaas/vald/pull/2195)
- Replace x/slices with standard slices pkg [#2193](https://github.com/vdaas/vald/pull/2193)
- add benchmark and check program for core ngt [#2179](https://github.com/vdaas/vald/pull/2179)

:bug: Bugfix

- Revert vtpool for ResourceExhausted problem [#2255](https://github.com/vdaas/vald/pull/2255)
- Fix deleted contour ingress controller apply [#2229](https://github.com/vdaas/vald/pull/2229)

:pencil2: Document

- Add document for RemoveByTimestamp RPC [#2238](https://github.com/vdaas/vald/pull/2238)

:green_heart: CI

- Disable exhaustruct [#2240](https://github.com/vdaas/vald/pull/2240)
- Fix fails when there are format differences [#2226](https://github.com/vdaas/vald/pull/2226)

:arrow_up: Update dependencies

- update deps [#2208](https://github.com/vdaas/vald/pull/2208)
- update dependencies [#2260](https://github.com/vdaas/vald/pull/2260)

## v1.7.8

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.8</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.8</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.8</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.8</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.8</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.8</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.8</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.8</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.8)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.8/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.8/charts/vald-helm-operator/README.md)

### Changes

♻️ Refactor

- change default creation poolsize [#2190](https://github.com/vdaas/vald/pull/2190)
- List kvs and vqueue data [#2188](https://github.com/vdaas/vald/pull/2188)
- refactor semver ci [#2189](https://github.com/vdaas/vald/pull/2189)

## v1.7.7

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.7</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.7</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.7</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.7</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.7</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.7</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.7</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.7</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.7)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.7/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.7/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Add CopyBuffer to copy.go [#2167](https://github.com/vdaas/vald/pull/2167)
- Add Vald standard internal/sync package [#2153](https://github.com/vdaas/vald/pull/2153)
- Add RemoveByTimestamp RPC [#2158](https://github.com/vdaas/vald/pull/2158)
- Implement StreamListObject and its tests [#2145](https://github.com/vdaas/vald/pull/2145)
- Add apiversion capabilities check to helm template [#2137](https://github.com/vdaas/vald/pull/2137)
- Add timestamp field to Object.Vector [#2136](https://github.com/vdaas/vald/pull/2136)
- Add gache's generic Map as internal/sync.Map and replace standard sync.Map with it [#2115](https://github.com/vdaas/vald/pull/2115)
- Make internal/cache generic [#2104](https://github.com/vdaas/vald/pull/2104)
- Install additional tools for docker image for devcontainer [#2101](https://github.com/vdaas/vald/pull/2101)
- Install buf and apply buf format [#2094](https://github.com/vdaas/vald/pull/2094)
- Add backup origin when CoW enabled and failed to load primary [#2091](https://github.com/vdaas/vald/pull/2091)
- Add decode kvsdb tool [#2059](https://github.com/vdaas/vald/pull/2059)
- Add user custom network policy [#2078](https://github.com/vdaas/vald/pull/2078)

:recycle: Refactor

- Refactor agent ngt core. [#2172](https://github.com/vdaas/vald/pull/2172)
- Refactor proto [#2173](https://github.com/vdaas/vald/pull/2173)
- Refactor search status [#2168](https://github.com/vdaas/vald/pull/2168)
- Refactor internal/core/algorithm/ngt mutex lock timing [#2144](https://github.com/vdaas/vald/pull/2144)
- Refactor github actions [#2141](https://github.com/vdaas/vald/pull/2141)
- Update license text [#2169](https://github.com/vdaas/vald/pull/2169)
- Refactor agent error not to wrap with details for performance issue [#2154](https://github.com/vdaas/vald/pull/2154)
- Use internal comparator instead of go-cmp [#2132](https://github.com/vdaas/vald/pull/2132)
- Refactor context [#2121](https://github.com/vdaas/vald/pull/2121)
- Propagate context to Search operation. [#2117](https://github.com/vdaas/vald/pull/2117)
- Refactor fix url http to https [#2090](https://github.com/vdaas/vald/pull/2090)
- Update "make gotests/gen" command [#2085](https://github.com/vdaas/vald/pull/2085)

:bug: Bugfix

- Fix duplicate make command [#2165](https://github.com/vdaas/vald/pull/2165)
- Add timestamp check for GetObject e2e [#2142](https://github.com/vdaas/vald/pull/2142)
- Modified apiversion capabilities check [#2149](https://github.com/vdaas/vald/pull/2149)
- Fix ngt index path of test case [#2130](https://github.com/vdaas/vald/pull/2130)
- Fix hack/benchmark search interface change [#2129](https://github.com/vdaas/vald/pull/2129)
- Fix internal/gache definition variable type [#2123](https://github.com/vdaas/vald/pull/2123)
- Use GOBIN instead of GOPATH/bin [#2102](https://github.com/vdaas/vald/pull/2102)
- Fix jaeger operator wait logic [#2114](https://github.com/vdaas/vald/pull/2114)
- Fix make k8s/metrics/jaeger/deploy failure [#2077](https://github.com/vdaas/vald/pull/2077)
- Bugfix Makefile KUBECONFIG recursive reference [#2089](https://github.com/vdaas/vald/pull/2089)
- Fix deploy command [#2088](https://github.com/vdaas/vald/pull/2088)
- Fix non-trusted module problem of v1.7.6 and disable not found debug message [#2076](https://github.com/vdaas/vald/pull/2076)
- Bugfix lb gateway pacicked caused by pairing heap search aggregator makes nil pointer when empty search result [#2181](https://github.com/vdaas/vald/pull/2181)

:pencil2: Document

- Update testing guideline for updated testing policy [#2131](https://github.com/vdaas/vald/pull/2131)
- Add troubleshooting for each rpc [#2163](https://github.com/vdaas/vald/pull/2163)
- Fix format network policy document [#2108](https://github.com/vdaas/vald/pull/2108)
- Add broken index backup document [#2096](https://github.com/vdaas/vald/pull/2096)
- Add network policy document [#2095](https://github.com/vdaas/vald/pull/2095)
- Fix 404 URL link [#2098](https://github.com/vdaas/vald/pull/2098)
- Update observability document [#2086](https://github.com/vdaas/vald/pull/2086)
- Fix typo of contribution guide [#2087](https://github.com/vdaas/vald/pull/2087)
- Update docs: search API and client API config [#2081](https://github.com/vdaas/vald/pull/2081)

:white_check_mark: Testing

- Re-Generate test codes [#2107](https://github.com/vdaas/vald/pull/2107)
- Update golangci-lint configuration: use white-list configuration pattern [#2106](https://github.com/vdaas/vald/pull/2106)

:green_heart: CI

- Fix coverage CI error [#2150](https://github.com/vdaas/vald/pull/2150)
- Remove some linters to make ci faster [#2116](https://github.com/vdaas/vald/pull/2116)

:chart_with_upwards_trend: Metrics/Tracing

- Divide latency of CreateIndex and SaveIndex metrics [#2099](https://github.com/vdaas/vald/pull/2099)
- Add broken index count metrics [#2083](https://github.com/vdaas/vald/pull/2083)

:arrow_up: Update dependencies

- Update go modules [#2092](https://github.com/vdaas/vald/pull/2092)

:art: Design

- Modified svg images [#2178](https://github.com/vdaas/vald/pull/2178)

## v1.7.6

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.6</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.6</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.6</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.6</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.6</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.6</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.6</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.6</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.6)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.6/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.6/charts/vald-helm-operator/README.md)

### Changes

⚠️ ⚠️ ⚠️ Warning ⚠️ ⚠️ ⚠️

- `v1.7.6` does not support `vald-agent-sidecar` for some reason.
- You can use the `vald-agent-sidecar` by setting the `vald-agent-ngt` image tag as `v1.7.5` or earlier.
- We will support `vald-agent-sidecar` in the future version again.

:sparkles: New feature

- Add search algorithm benchmark and update search aggregation algo [#2044](https://github.com/vdaas/vald/pull/2044)
- Add broken index backup [#2034](https://github.com/vdaas/vald/pull/2034)
- Add network policy [#2022](https://github.com/vdaas/vald/pull/2022)

:recycle: Refactor

- Add save index operation log [#2048](https://github.com/vdaas/vald/pull/2048)
- Added flg that can disable to ingress defaultBackend [#1976](https://github.com/vdaas/vald/pull/1976)
- Refactor and Add test for service/ngt.go [#2040](https://github.com/vdaas/vald/pull/2040)
- Add e2e envs to devcontainer [#2032](https://github.com/vdaas/vald/pull/2032)
- Update RoundTrip retry condition [#2033](https://github.com/vdaas/vald/pull/2033)

:bug: Bugfix

- Fix fp16 problems [#2049](https://github.com/vdaas/vald/pull/2049)
- Add KUBECTL_VERSION value to workflow [#2052](https://github.com/vdaas/vald/pull/2052)
- Remove sudo from kubectl and small refactor around os/arch [#2037](https://github.com/vdaas/vald/pull/2037)
- Disable vtproto pooling due to the performance degradation [#2063](https://github.com/vdaas/vald/pull/2063)
- Fix to create index_path when it does not exists [#2060](https://github.com/vdaas/vald/pull/2060)

:pencil2: Document

- Add documentation for devcontiner [#2042](https://github.com/vdaas/vald/pull/2042)
- Create README for each docker image [#2014](https://github.com/vdaas/vald/pull/2014)

:green_heart: CI

- Disable deepsource TestCoverage due to the Deepsource Coverage collect server timeout is too short for Vald testing [#2038](https://github.com/vdaas/vald/pull/2038)
- Update Docker Build workflow with forked sources [#2036](https://github.com/vdaas/vald/pull/2036)
- Fix e2e-max-dim test [#2028](https://github.com/vdaas/vald/pull/2028)
- Fix E2E actions on PR [#2025](https://github.com/vdaas/vald/pull/2025)
- Change E2E actions to use local charts on PR [#2024](https://github.com/vdaas/vald/pull/2024)
- Update format chatops [#2021](https://github.com/vdaas/vald/pull/2021)
- Format code with prettier and gofumpt [#2015](https://github.com/vdaas/vald/pull/2015)

:chart_with_upwards_trend: Metrics/Tracing

- Add command to deploy monitoring stack [#2030](https://github.com/vdaas/vald/pull/2030)
- Fixed duplicate counting in CPU graphs [#2019](https://github.com/vdaas/vald/pull/2019)

:arrow_up: Update dependencies

- Update go modules [#2053](https://github.com/vdaas/vald/pull/2053)
- Update NGT version [#2026](https://github.com/vdaas/vald/pull/2026)

:handshake: Contributor

- Add takuyaymd as a contributor for maintenance [#2020](https://github.com/vdaas/vald/pull/2020)

## v1.7.5

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.5</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.5</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.5</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.5</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.5</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.5</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.5</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.5</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.5)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.5/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.5/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New features

- Add index with timestamp [#1984](https://github.com/vdaas/vald/pull/1984)

:recycle: Refactor

- Improve errors.Join performance [#2010](https://github.com/vdaas/vald/pull/2010)
- Refactor error returning timing of doSearch function [#1996](https://github.com/vdaas/vald/pull/1996)
- Refactor makefile for non GOPATH strict environment #2 [#1998](https://github.com/vdaas/vald/pull/1998)
- Refactor makefile for non GOPATH strict environment [#1997](https://github.com/vdaas/vald/pull/1997)

:bug: Bugfix

- Correction of a bug that returned NotFound as success 0 when balancedUpdate is disabled and all ReplicaAgents are AlreadyExists (already have the exact same Index). [#2011](https://github.com/vdaas/vald/pull/2011)
- Refactor replace errors wrap with join [#2001](https://github.com/vdaas/vald/pull/2001)
- Remove nvimlog [#1994](https://github.com/vdaas/vald/pull/1994)

:green_heart: CI

- Fix chatops format workflow [#2007](https://github.com/vdaas/vald/pull/2007)
- Fix incorrect error output of gen-test chatopts command [#2004](https://github.com/vdaas/vald/pull/2004)
- Fix Makefile bug and update deps for checking bugfix [#2002](https://github.com/vdaas/vald/pull/2002)
- Output error to chatops comment [#1999](https://github.com/vdaas/vald/pull/1999)
- Fix ChatOpts /gen-test command error [#1993](https://github.com/vdaas/vald/pull/1993)

:pencil2: Document

- Update unit test guideline for unimplemented test [#1983](https://github.com/vdaas/vald/pull/1983)

:white_check_mark: Testing

- Implement generic function tests [#2008](https://github.com/vdaas/vald/pull/2008)
- Generate empty test using /gen-test ChatOpts command [#2005](https://github.com/vdaas/vald/pull/2005)
- Update internal/info test and add new case for coverage [#2003](https://github.com/vdaas/vald/pull/2003)

:arrow_up: Update dependencies

- Update go module and libs [#2012](https://github.com/vdaas/vald/pull/2012)

:handshake: Contributor

- Add ykadowak as a contributor for code, and test [#2009](https://github.com/vdaas/vald/pull/2009)

## v1.7.4

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.4</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.4</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.4</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.4</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.4</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.4</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.4</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.4</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.4)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.4/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.4/charts/vald-helm-operator/README.md)

### Changes

:bug: Bug fix

- Fix range concurrency branch rule [#1986](https://github.com/vdaas/vald/pull/1986)
- Update makefile for "not implemented" placeholder [#1967](https://github.com/vdaas/vald/pull/1977)
- Non-gRPC style error parse result returns Unknown status, it should be re-parse to find inside status [#1981](https://github.com/vdaas/vald/pull/1981)
- Enable gorules [#1980](https://github.com/vdaas/vald/pull/1980)
- Format code with prettier and gofumpt [#1971](https://github.com/vdaas/vald/pull/1971)

:memo: Document

- Fix documentation typo disable_balanced_update [#1978](https://github.com/vdaas/vald/pull/1978)

:handshake: Contributor

- docs: add junsei-ando as a contributor for doc [#1979](https://github.com/vdaas/vald/pull/1979)

## v1.7.3

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.3</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.3</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.3</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.3</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.3</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.3</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.3</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.3</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.3)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.3/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.3/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New feature

- Add disable balanced update features & use generic type for BidirectionalStream [#1964](https://github.com/vdaas/vald/pull/1964)

:bug: Bug fix

- Fix grafana dashboard query for backoff retry count [#1961](https://github.com/vdaas/vald/pull/1961)

:recycle: Refactor

- Refactor conv.go [#1968](https://github.com/vdaas/vald/pull/1968)

:memo: Document

- Add new API parameter and update observability docs [#1966](https://github.com/vdaas/vald/pull/1966)

:arrow_up: Dependency update

- Update deps [#1969](https://github.com/vdaas/vald/pull/1969)

## v1.7.2

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.2</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.2</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.2</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.2</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.2</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.2</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.2</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.2</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.2)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.2/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.2/charts/vald-helm-operator/README.md)

### Changes

:bug: Bug fix

- Vald gRPC Client and Pool logic makes huge backoff [#1953](https://github.com/vdaas/vald/pull/1953)
- Missing backoff metrics [#1958](https://github.com/vdaas/vald/pull/1958)

:recycle: Refactor

- Update test template to exclude deepsource warning [#1954](https://github.com/vdaas/vald/pull/1954)

:white_check_mark: Test

- Remove non-implemented test [#1952](https://github.com/vdaas/vald/pull/1952)

## v1.7.1

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.1</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.1</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.1</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.1</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.1</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.1</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.1</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.1</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.1)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.1/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.1/charts/vald-helm-operator/README.md)

### Changes

#### 🐛 Bugfix

- [bugfix] add target length validation for each gRPC client exection method [#1939](https://github.com/vdaas/vald/pull/1939)

#### ♻️ Refactor

- update gRPC status code for API docs [#1943](https://github.com/vdaas/vald/pull/1943)
- Refactor: Add t.Helper() on test helper function [#1935](https://github.com/vdaas/vald/pull/1935)
- Fix syntax error on dump context workflow [#1936](https://github.com/vdaas/vald/pull/1936)
- format codes [#1934](https://github.com/vdaas/vald/pull/1934)

## v1.7.0

### Docker images

<table>
  <tr>
    <th>component</th>
    <th>Docker pull</th>
  </tr>
  <tr>
    <td>Agent NGT</td>
    <td>
      <code>docker pull vdaas/vald-agent-ngt:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-ngt:v1.7.0</code>
    </td>
  </tr>
  <tr>
    <td>Agent sidecar</td>
    <td>
      <code>docker pull vdaas/vald-agent-sidecar:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-agent-sidecar:v1.7.0</code>
    </td>
  </tr>
  <tr>
    <td>Discoverers</td>
    <td>
      <code>docker pull vdaas/vald-discoverer-k8s:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-discoverer-k8s:v1.7.0</code>
    </td>
  </tr>
  <tr>
    <td>Gateways</td>
    <td>
      <code>docker pull vdaas/vald-lb-gateway:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-lb-gateway:v1.7.0</code><br/>
      <code>docker pull vdaas/vald-filter-gateway:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-filter-gateway:v1.7.0</code>
    </td>
  </tr>
  <tr>
    <td>Index Manager</td>
    <td>
      <code>docker pull vdaas/vald-manager-index:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-manager-index:v1.7.0</code>
    </td>
  </tr>
  <tr>
    <td>Helm Operator</td>
    <td>
      <code>docker pull vdaas/vald-helm-operator:v1.7.0</code><br/>
      <code>docker pull ghcr.io/vdaas/vald/vald-helm-operator:v1.7.0</code>
    </td>
  </tr>
</table>

### Documents

- [GoDoc](https://pkg.go.dev/github.com/vdaas/vald@v1.7.0)
- [Helm Chart Reference](https://github.com/vdaas/vald/blob/v1.7.0/charts/vald/README.md)
- [Helm Operator Chart Reference](https://github.com/vdaas/vald/blob/v1.7.0/charts/vald-helm-operator/README.md)

### Changes

:sparkles: New Feature

- Introduce OTLP for metrics and trace [#1824](https://github.com/vdaas/vald/pull/1824)
- Add manifest to deploy opentelemetry-operator [#1819](https://github.com/vdaas/vald/pull/1819)

:arrow_up: Dependencies update

- Update go modules and go version [#1922](https://github.com/vdaas/vald/pull/1922)
- Update go modules [#1904](https://github.com/vdaas/vald/pull/1904)

:bug: Bug fix

- Divide agent errors for QBG and Faiss implementation [#1924](https://github.com/vdaas/vald/pull/1924)
- Bugfix status handling for idle connection [#1921](https://github.com/vdaas/vald/pull/1921)
- Fix invalid character error [#1914](https://github.com/vdaas/vald/pull/1914)
- Fixed duplicate counts of working_memory_set_bytes [#1911](https://github.com/vdaas/vald/pull/1911)
- Bug fix using undefined a variable at maxDimensionTest [#1856](https://github.com/vdaas/vald/pull/1856)
- Bug fix prometheus export return value [#1817](https://github.com/vdaas/vald/pull/1817)

:recycle: Refactor

- Happy New Year 2023 [#1918](https://github.com/vdaas/vald/pull/1918)
- Add auto-update libs & deps make command [#1917](https://github.com/vdaas/vald/pull/1917)
- Add canceled status for CreateIndex API [#1892](https://github.com/vdaas/vald/pull/1892)
- Update concurrent cancellation group name [#1912](https://github.com/vdaas/vald/pull/1912)
- Remove blank when all parameters are not used and Add ErrJobFuncNotFound [#1879](https://github.com/vdaas/vald/pull/1879)
- Rename doXXX() [#1878](https://github.com/vdaas/vald/pull/1878)
- Remove deprecated functional option (internal/net/grpc) [#1877](https://github.com/vdaas/vald/pull/1877)
- Fix deepsource: RVV-B0001 Confusing naming of struct fields or methods [#1875](https://github.com/vdaas/vald/pull/1875)
- Fix deepsource: VET-V0008 lock erroneously passed by value internal/net test [#1874](https://github.com/vdaas/vald/pull/1874)
- Fix deepsource: RVV-B0001 confusing naming of struct fields or methods [#1844](https://github.com/vdaas/vald/pull/1844)
- Fix deepsource: SCC-U1000 Unused code [#1873](https://github.com/vdaas/vald/pull/1873)
- Fix deepsource: RVV-B0006 Method modifies receiver [#1872](https://github.com/vdaas/vald/pull/1872)
- Fix deepsource: SCC-SA4006 Value assigned to a variable is never read before being overwritten [#1871](https://github.com/vdaas/vald/pull/1871)
- Fix deepsource: VET-V0008 Lock erroneously passed by value (pkg/agent) [#1868](https://github.com/vdaas/vald/pull/1868)
- Fix deepsource: VET-V0008 lock erroneously passed by value internal/info test [#1869](https://github.com/vdaas/vald/pull/1869)
- Fix deepsource: DOK-W1001 found consecutive run command [#1870](https://github.com/vdaas/vald/pull/1870)
- Fix deepsource: VET-V0008 lock erroneously passed by value (internal/net) [#1867](https://github.com/vdaas/vald/pull/1867)
- Fix deepsource: CRT-D0001 append possibly assigns to a wrong variable [#1866](https://github.com/vdaas/vald/pull/1866)
- Fix deepsource: VET-V0008 Lock erroneously passed by value (pkg/manager) [#1861](https://github.com/vdaas/vald/pull/1861)
- Fix deepsource: VET-V0008 lock erroneously passed by value pkg/discoverer [#1857](https://github.com/vdaas/vald/pull/1857)
- Fix deepsource: RVV-B0006 Method modifies receiver [#1865](https://github.com/vdaas/vald/pull/1865)
- Fix deepsource: VET-V0008 Lock erroneously passed by value (internal/test, singleflight, observability) [#1863](https://github.com/vdaas/vald/pull/1863)
- Fix deepsource: VET-V0008 lock erroneously passed by value info [#1864](https://github.com/vdaas/vald/pull/1864)
- Fix deepsource: VET-V0008 Lock erroneously passed by value internal/client [#1862](https://github.com/vdaas/vald/pull/1862)
- Fix deepsource: VET-V0008 lock erroneously passed by value internal/info,iocopy,errgroup [#1860](https://github.com/vdaas/vald/pull/1860)
- Fix deepsource: VET-V0008 Lock erroneously passed by value internal/db, backoff, circuitbreaker [#1859](https://github.com/vdaas/vald/pull/1859)
- Fix deepsource: RVV-B0009 Redefinition of builtin [#1858](https://github.com/vdaas/vald/pull/1858)
- Fix deepsource: CRT-A0014 switch single case can be rewritten as if or if-else [#1855](https://github.com/vdaas/vald/pull/1855)
- Fix deepsource: RVV-A0003 Exit inside non-main function ./hack [#1854](https://github.com/vdaas/vald/pull/1854)
- Fix deepsource: SCC-S1003 replace call to strings.Index with strings.Contains [#1853](https://github.com/vdaas/vald/pull/1853)
- Fix deepsource: RVV-B0013 Unused method receiver [#1852](https://github.com/vdaas/vald/pull/1852)
- Fix deepsource: CRT-D0007 Duplicate cases found in switch statement [#1851](https://github.com/vdaas/vald/pull/1851)
- Fix deepsource: GO-W1009 using a deprecated function, variable, constant or field [#1846](https://github.com/vdaas/vald/pull/1846)
- Fix deepsource: RVV-B0011 exported function returning value of unexported type [#1848](https://github.com/vdaas/vald/pull/1848)
- Fix deepsource: GSC-G103 Function call made to an unsafe package [#1850](https://github.com/vdaas/vald/pull/1850)
- Fix deepsource: RVV-B0012 Unused parameter in the function
- Fix deepsource: DOK-SC2002, DOK-W1001 Useless cat and Multiple consecutive RUN [#1847](https://github.com/vdaas/vald/pull/1847)
- Fix deepsource: dockerfile warning [#1835](https://github.com/vdaas/vald/pull/1835)
- Fix deepsource RVV-B0013 [#1832](https://github.com/vdaas/vald/pull/1832)
- Fix deepsource VET-V0007 unkeyed composite literals [#1837](https://github.com/vdaas/vald/pull/1837)
- Fix deepsource: Audit required: Insecure gRPC server [#1833](https://github.com/vdaas/vald/pull/1833)
- Fix deepsource: Potential slowloris attack [#1834](https://github.com/vdaas/vald/pull/1834)
- Fix deepsource: Unsafe defer of os.Close [#1836](https://github.com/vdaas/vald/pull/1836)
- Fix deepsource: RVV-A0003 exit inside non-main function [#1838](https://github.com/vdaas/vald/pull/1838)
- Fix deepsource: GSC-G404 Audit the random number generation source (rand) [#1839](https://github.com/vdaas/vald/pull/1839)
- Fix makefile [#1828](https://github.com/vdaas/vald/pull/1828)
- Refactor circuitbreaker [#1816](https://github.com/vdaas/vald/pull/1816)

:white_check_mark: Test

- Refactor Insert Upsert Testing [#1919](https://github.com/vdaas/vald/pull/1919)

:green_heart: CI

- Ci/GitHub action docker/update docker login action [#1903](https://github.com/vdaas/vald/pull/1903)
- Add actions workflow validation [#1902](https://github.com/vdaas/vald/pull/1902)
- Change docker build permission [#1901](https://github.com/vdaas/vald/pull/1901)
- Update docker login action [#1900](https://github.com/vdaas/vald/pull/1900)
- Format code with prettier and gofumpt [#1886](https://github.com/vdaas/vald/pull/1886)
- Update deepsource configuration [#1881](https://github.com/vdaas/vald/pull/1881)
- Update gotestfmt org [#1880](https://github.com/vdaas/vald/pull/1880)
- Add escape for e2e workflow [#1845](https://github.com/vdaas/vald/pull/1845)
- Resolve GitHub Actions warning [#1818](https://github.com/vdaas/vald/pull/1818)

:memo: Document

- Create observability configuration document [#1882](https://github.com/vdaas/vald/pull/1882)
- Add takuyaymd as a contributor for bug, and code [#1913](https://github.com/vdaas/vald/pull/1913)
- Update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE [#1885](https://github.com/vdaas/vald/pull/1885)
- Fix typo comment [#1831](https://github.com/vdaas/vald/pull/1831)
- Add filter gateway api doc [#1821](https://github.com/vdaas/vald/pull/1821)
- Fix dead link [#1823](https://github.com/vdaas/vald/pull/1823)
- Update pull request template [#1820](https://github.com/vdaas/vald/pull/1820)


