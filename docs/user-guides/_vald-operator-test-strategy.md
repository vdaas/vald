# Test Strategy

## Background

The controller is built on a kubebuilder scaffold and uses two test styles: Ginkgo/Gomega
integration tests on top of `envtest`, and plain `testing` + testify unit tests. Stating
which layer owns which behavior avoids both duplication and gaps.

## Test layers

### Unit tests (`testing` + testify)

Cover pure Go logic — anything that does not depend on the Kubernetes API server or an
external client.

| Package | Under test |
|---|---|
| `internal/pkg/lifecycle` | `MakeCondition`, `GetIndex`, `GetNext` |
| `internal/pkg/domain/valdoperatorrelease` | `ConditionWaitForClusterCreate` (Check branches), `ResolveAgentNodePool` |
| `api/v1` | `GetNodePool`, `GetResourceList` |
| `internal/pkg/api/valdrelease` | `SetRelationalResources` and per-component resource methods |
| `internal/pkg/lifecycle/builder/vald` | `validate()`, the component builders, and `Build()` (golden file) |

### Integration tests (Ginkgo/Gomega + envtest)

`internal/controller/valdoperatorrelease_controller_test.go` and `resource_syncer_test.go` stand
up a real API server via `envtest` and drive the full reconcile loop. Behaviors that
require a Kubernetes client — phase transitions, status updates, CreateOrUpdate/prune —
are verified here.

## What is intentionally not unit-tested, and why

- **`zz_generated.deepcopy.go`** — kubebuilder-generated; not hand-written, so out of
  scope. It lowers the raw `go test -cover` number without being a real gap.
- **`desired.Resource.IsReady()` / `Build()`** — calls `client.Get` internally, so it
  cannot be unit-tested without a mock. The `ConditionWaitForCreateVrs` unit test stays at
  the level of "is `lc.Condition.Type` correct"; `Build`/`IsReady` behavior is left to the
  integration tests.
- **`ConditionCompleted`, `InitProgress`, `NewLifeCycleFlow`** — thin delegation to
  already-tested functions (e.g. `GetIndex`) or simple assignment; low marginal value.
  Covered indirectly by the integration phase-transition tests.

## Golden file test

`builder/vald/TestVrsBuilder_Build` compares the first item of the generated
`ValdRelease` list against `testdata/vrs.golden.yaml`. The expected file is overwritten
when the test runs, so intentional changes must be reviewed as a diff.

## Coverage philosophy

We prioritize unit coverage of hand-written logic, excluding generated code. The overall
`go test -cover` number looks low because generated code drags it down; per-function, the
main paths of the hand-written code are covered. Integration coverage is measured
separately via `make test-e2e` (requires an `envtest` environment).

## Framework choice

| Use | Framework | Reason |
|---|---|---|
| Unit tests | `testing` + testify | Standard, minimal deps, easy table-driven tests. |
| Integration tests | Ginkgo/Gomega | kubebuilder scaffold default; strong `envtest` affinity. |
