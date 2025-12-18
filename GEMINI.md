# GEMINI.md — Vald Development Guide

This document serves as a guide for AI agents (and human developers) working on the Vald repository. It outlines the development workflow, tooling, coding standards, and common tasks.

> **Role:** Senior Software Engineer / Collaborative Agent
> **Goal:** High-quality, tested, and verifiable code changes.

---

## 1. Core Workflow

1.  **Explore**: Understand the codebase using `list_files`, `read_file`, and `grep` (via `run_in_bash_session`).
2.  **Plan**: Create a step-by-step plan using `set_plan`.
3.  **Execute**: Implement changes incrementally.
4.  **Verify**: **Always** verify changes.
    *   Run linters: `make lint`
    *   Format code: `make format`
    *   Run tests: `make test` (unit) or `make e2e/v2` (end-to-end)
5.  **Submit**: Create a pull request with a clear description.

---

## 2. Build System & Make Targets

Vald uses a comprehensive `Makefile` system. **Always check `make help`** for the most up-to-date targets.

### Testing
*   **Unit Tests (Go)**:
    *   `make test`: Runs tests for `cmd`, `internal`, `pkg`.
    *   `make test/all`: Runs all Go tests.
    *   `make test/pkg`: Runs tests for `pkg/`.
    *   `make test/internal`: Runs tests for `internal/`.
*   **Rust Tests**:
    *   `make test/rust`: Runs Rust agent and QBG tests.
*   **E2E Tests**:
    *   `make e2e/v2`: Runs the **V2** end-to-end tests (Preferred).
    *   `make e2e`: Runs the classic E2E tests (`TestE2EStandardCRUD`).
    *   **Note**: E2E tests often require a local cluster (see below).

### Formatting & Linting
*   **Format All**: `make format` (Runs Go, Proto, JSON, YAML, MD formatters).
*   **Format Go**: `make format/go` (Runs `golines`, `gofumpt`, `goimports`, `strictgoimports`, `crlfmt`).
*   **Lint**: `make lint` (Runs `go vet`, `golangci-lint`, `textlint`, `cspell`, `reviewdog`).
*   **Update**: `make update` (Updates deps, licenses, and runs formatters).

### Local Development (Kubernetes)
*   **Start Cluster**: `make k3d/start` (Starts a local k3d cluster).
*   **Deploy Vald**: `make k8s/vald/deploy` (Deploys Vald to the current context).
*   **Delete Vald**: `make k8s/vald/delete`.
*   **Stop Cluster**: `make k3d/stop`.

---

## 3. Project Structure

*   `cmd/`: Main applications (Agent, Gateway, Discoverer, etc.).
*   `pkg/`: Public library code.
*   `internal/`: Private library code (Core logic often lives here).
*   `apis/`: Protobuf definitions (`.proto`) and generated code.
*   `rust/`: Rust components (AgentNGT, QBG).
*   `tests/v2/e2e/`: End-to-End Test Suite V2.
*   `hack/`: Scripts and build tools.
*   `charts/`: Helm charts.
*   `Makefile` & `Makefile.d/`: Build configuration.

---

## 4. Coding Standards

### Go
*   **Formatting**: Strictly enforced. **Always run `make format/go`** before verifying/submitting.
    *   Tools used: `golines` (line wrapping), `gofumpt` (strict fmt), `goimports` (imports management), `strictgoimports`, `crlfmt`.
*   **Linting**: `golangci-lint` must pass.
*   **Testing**:
    *   Use table-driven tests.
    *   Place unit tests in `*_test.go` files next to the source.
    *   Mock interfaces where appropriate using the internal mock framework or standard techniques.

### Rust
*   Standard `cargo fmt` and `cargo test` workflows apply, wrapped by Make targets (`make test/rust`).

### Protocol Buffers
*   Modify `.proto` files in `apis/proto/`.
*   Run `make proto/all` to regenerate Go/Rust code.

---

## 5. E2E Testing Workflow (V2)

The V2 E2E suite (`tests/v2/e2e`) is the modern way to verify system behavior.

1.  **Prerequisites**:
    *   Running Kubernetes cluster (`make k3d/start`).
    *   Datasets (e.g., Fashion-MNIST) in `hack/benchmark/assets/`.
2.  **Execution**:
    ```bash
    make k8s/vald/deploy
    make e2e/v2
    ```
3.  **Configuration**:
    *   Tests are configured via YAML files in `tests/v2/e2e/assets/`.
    *   You can run specific scenarios using `make e2e/v2` which defaults to `TestE2EStrategy`.

---

## 6. Contribution Guidelines

*   **Branch Naming**: `[type]/[area]/[description]`
    *   Types: `feature`, `bug`, `refactoring`, `test`, `ci`, `doc`.
    *   Example: `feature/gateway/add-new-filter`
*   **Pull Requests**:
    *   Keep them small and focused.
    *   Include tests and benchmark results if performance is affected.
    *   Update documentation if behavior changes.

---

## 7. Troubleshooting

*   **"Command not found"**: Ensure you are running commands via `run_in_bash_session`. The environment should have `go`, `make`, `kubectl`, `helm` pre-installed or accessible.
*   **Lint Failures**: Read the output of `make lint` carefully. `golangci-lint` often gives specific instructions.
*   **Test Failures**: Use `go test -v` on the specific package to see detailed logs.
*   **Dirty Git State after Build**: Run `make format` to ensure generated code matches the expected format.

---

**Note**: This document is a living guide. Update it if you find new patterns or tools that improve the development workflow.
