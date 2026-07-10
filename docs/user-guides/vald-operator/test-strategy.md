# Test Strategy

## Background

The controller uses a single test style: plain `testing` + testify unit tests.
All behavior is covered at the unit level — there are no integration tests or envtest
dependencies.

## Test layers

### Unit tests (`testing` + testify)

Cover pure Go logic — anything that does not depend on a live Kubernetes API server.

| Package                                  | Under test                                                                           |
| ---------------------------------------- | ------------------------------------------------------------------------------------ |
| `pkg/operator/vald/service` (builder)    | `validate()`, component builders, `Build()`, `mergeOverlay` (golden file)            |
| `pkg/operator/vald/service` (capability) | `AlwaysAvailable`, `ResolveNodePoolCapability`                                       |
| `pkg/operator/vald/service` (rules)      | `resolveAgentNodePool` fallback logic                                                |
| `pkg/operator/vald/service` (reconciler) | Reconcile loop, phase transitions, status updates (client faked via interface mocks) |
| `pkg/operator/vald/service` (operator)   | `New`, `Operator.Start` cancellation                                                 |
| `pkg/operator/vald/config`               | Config loading and defaults                                                          |

## What is intentionally not unit-tested, and why

- **`zz_generated.deepcopy.go`** — generated code; not hand-written, so out of scope.
  It lowers the raw `go test -cover` number without being a real gap.
- **`reconcilePhase` Kubernetes side-effects** — `CreateOrUpdate` / `Prune` against a real
  API server are not exercised in unit tests. The reconciler tests fake the client, so
  the unit suite validates the phase logic, not the wire protocol. E2E tests cover the
  full reconcile against a real cluster.

## Golden file test

`TestVrsBuilder_Build` compares the first item of the generated `ValdRelease` list
against `testdata/vrs.golden.yaml`. The golden file is updated when the test is run with
`-update`, so intentional changes must be reviewed as part of the diff.

## Coverage philosophy

We prioritize unit coverage of hand-written logic, excluding generated code. The overall
`go test -cover` number looks low because generated code drags it down; per-function, the
main paths of the hand-written code are covered.

## Framework choice

| Use        | Framework           | Reason                                           |
| ---------- | ------------------- | ------------------------------------------------ |
| Unit tests | `testing` + testify | Standard, minimal deps, easy table-driven tests. |
