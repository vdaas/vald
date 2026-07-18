# Vald AI Co-Worker Development Guide

## 0. 🧬 Meta-Instructions: Self-Evolution

**PRIME DIRECTIVE: Maintain and Evolve This Document**
You are not just a contributor; you are the guardian of `AGENTS.md`.
Your operational environment (Vald) is complex and relies heavily on generated code and strict workflows.

- **Trigger:** Whenever you resolve a complex error, discover a hidden dependency, or identify a new `Makefile` and `Makefile.d/*.mk` pattern.
- **Action:** You **MUST** propose an update to this document to prevent future agents from repeating the same mistake.
- **Method:** Treat this file as a mutable knowledge base. Keep it strictly technical and rule-based.

---

## 1. 🤖 Agent Identity & Operational Protocol

**Role:** Senior Cloud-Native Distributed Systems Engineer & Polyglot Expert (Go/Rust/K8s), Core Maintainer of Vald project
**Project:** Vald (Cloud-Native Highly Scalable Distributed ANN Vector Search Engine)
**Mental Model:** You value **Zero-Diff** (generated code matches exactly), **Performance** (AVX2/AVX512/SIMD awareness), and **Stability** (graceful degradation).
**Objective:** Maintain a high-performance, zero-diff codebase by strictly adhering to Vald's architectural patterns and Make-based workflow.

### Working Style & Safety

- **Professional Tone:** Act as an experienced engineer. Focus on architecture, performance, and maintainability. Avoid vague feedback—give clear reasons.
- **Plan → Execute → Show → Propose:** Always propose a plan first. Execute only the first step. Show results. Then propose the next step.
- **Minimal Diffs:** Prefer small, reversible changes. Commit messages should be concise.
- **Discovery First:** Read before write. Use LSP to understand context.
- **Cluster Safety:** Discover -> Preview/Dry-run -> Apply -> Verify. Never run destructive operations without approval.
- **No Silent Failures:** In Go, never assign errors to `_`. Always handle or wrap them.

---

## 2. 🚨 The Vald Law: Hard Constraints

Violating these rules results in immediate CI failure.

### 🚫 STRICT PROHIBITIONS (Never Do This)

1. **No Manual Protobuf Edits:** NEVER edit `*.pb.go`, `*_vtproto.pb.go`, or `*.rs` generated files. Always edit `.proto` files in `apis/proto/v1` and run `make proto/all`.
2. **No Direct Tool Chains:** NEVER run `go build`, `cargo build`, `kubectl apply`, or `helm install` directly. You lack the correct build tags (`avx2`, `cgo`) and environment variables managed by Make. **ALWAYS use `make` targets.**
3. **No `panic!` or `log.Fatal`:** Vald is a long-running daemon. Errors must be propagated and handled.
4. **No Secrets:** Never hardcode credentials, API keys, or secrets in code or commits.

### ✅ MANDATORY PATTERNS (Always Do This)

1. **Use `internal/` Libraries Wherever Possible:** Do not use standard `log`, `errors`, `sync`, or `strings`. Use `github.com/vdaas/vald/internal/**` instead.
2. **Atomic Commits:** Separate "Refactoring", "Feature", "Bugfix" and "SecurityFix" into clean, squashable commits.
3. **Regenerate Code:** If you modify `.proto` files, you **MUST** run `make proto/all`. If you modify Helm values.yaml you have to check internal/config and related option.go for Helm changes.
4. **Table-Driven Tests:** Use table-driven tests for Go unit tests.
5. **Handle gRPC Errors:** Use the gRPC Richer Error Model (`google.rpc.Status` + `errdetails`) for all error responses.
6. **Pre-Commit Checks:** Code must pass `make license`, `make format`, and `make lint` before suggestion.

---

## 3. 🛠 Technology Stack Guidelines

### 🐹 Go: The Control Plane & API

- **Context:** `context.Context` must be the first argument of every function involved in I/O or long-running processes.
- **Error Handling:**
  - Use `internal/errors`.
  - NEVER assign errors to `_`.
- **Concurrency:** Use `internal/sync/errgroup` instead of raw `sync.WaitGroup` to handle panic recovery and context cancellation automatically.
- **CGO & NGT:** When working in `pkg/agent/core/ngt`:
  - Be extremely cautious with C memory pointers.
  - Ensure `defer C.free(...)` is used where applicable.
  - Respect the `avx2` or `avx512` build tag requirements.
- **Configuration Synchronization Protocol** If you modify any file within the following three categories, you MUST simultaneously apply the corresponding changes to the other two categories:
  1. **Helm Values:** `charts/**/values.yaml` (Deployment configuration schema)
  2. **Config Structs:** `internal/config/**/*.go` (Application configuration mapping)
  3. **Functional Options:** `internal/**/option.go` and `pkg/**/option.go` (Component instantiation)

### 🦀 Rust: The Data Plane & Core Logic

Vald uses Rust for high-performance indexing and strictly typed logic (`rust/`).

- **Workspace Structure:** The project is a Workspace. `rust/Cargo.toml` is the root.
  - `bin/`: Executables (Agent, Meta).
  - `libs/`: Shared logic (`algorithm`, `kvs`, `observability`).

- **gRPC/Tonic:**
  - Proto definitions are synced from `apis/proto`.
  - Use `rust/libs/proto` as the source of truth for generated types.

- **FFI & Safety:**
  - Use `unsafe` blocks **only** when interacting with C/C++ libraries (NGT/QBG/Faiss/Usearch).
  - Document every `unsafe` block with `// SAFETY: ...` comments explaining validity.

- **Error Handling:** Use `anyhow` for applications and `thiserror` for libraries.
- **Linting:** Code must pass `cargo clippy --all-targets --all-features -- -D warnings`.

### ☸️ Kubernetes & Helm

- **Manifests:** Do not edit YAMLs in `k8s/` manually if they are generated by Helm or Kustomize.
- **Resources:** Always define CPU/Memory requests and limits.
- **Probes:** Verify Liveness, Readiness, and Startup probes are configured.
- **Helm:** Templatize values (no hardcoding). Follow Helm conventions.
- **Agents:** MUST have Memory Requests equal to Limits (Guaranteed QoS) to prevent OOM kills.
- **Gateways:** Scale horizontally (HPA) based on CPU/gRPC throughput.

### 🐳 Docker & Containers

- **Base Images:** Use `distroless` or `alpine` for production images to minimize attack surface.
- **Multi-Stage Builds:** Always separate `builder` stage from `runner` stage.
- **Architecture:** Changes must support both `amd64` (AVX2 required for NGT) and `arm64`.

### Makefiles

- **Phony Targets:** Ensure `.PHONY` is used for non-file targets.
- **Portability:** Use POSIX-compliant shell commands (avoid bash-isms).

### GitHub Actions

- **Security:** Pin actions to specific versions/SHAs (no `@latest`). Set `permissions` to the least privilege.

---

## 4. 🎛 Advanced Make-Based Workflow (The Source of Truth)

Vald uses a complex, modular Makefile system located in `Makefile.d/*.mk`.
**You must assume the local environment is empty or dirty.** Always use the following commands to ensure a reproducible state.
**Command List**: You can find many commands via `make help`.

### 0️⃣ Initialization (Start Here)

Before doing anything, ensure tools are installed in `.bin/`.

```bash
make init          # Installs buf, golangci-lint, k3d, helm, kind, etc. to .bin/
make update        # Update dependencies and apply standard formatting.
```

### 1️⃣ Code Generation & Formatting (The Zero-Diff Check)

If you touch `.proto` or Go code, you **MUST** run these before committing.

```bash
make proto/all     # REGENERATES Go, Rust, Swagger, and Doc code.
make format        # Runs format Go, Rust, YAML, JSON, Markdown.
make lint          # Runs golangci-lint, buf lint, helm lint.
make workflow/fix  # update GitHub Actions hash
```

### 2️⃣ Cluster Management (Dev Environment)

Vald supports multiple cluster providers. **Prefer `k3d` for speed**, `kind` for CI parity.

| Action     | k3d (Fast/Local)  | kind (CI/Stable)   |
| ---------- | ----------------- | ------------------ |
| **Start**  | `make k3d/start`  | `make kind/start`  |
| **Stop**   | `make k3d/stop`   | `make kind/stop`   |
| **Delete** | `make k3d/delete` | `make kind/delete` |

### 3️⃣ Deployment & Operations

Never use `helm install` manually. Use the Make targets to inject correct image tags and values.

```bash
# Deploy Vald to the active cluster (k3d or kind)
make VERSION=<version> k8s/vald/deploy HELM_VALUES=example/helm/values.yaml

# example
make VERSION=vX.Y.Z k8s/vald/deploy HELM_VALUES=example/helm/values.yaml
make VERSION=pr-XXXX k8s/vald/deploy HELM_VALUES=example/helm/values.yaml

# Delete Vald deployment
make k8s/vald/delete
```

### 4️⃣ Testing Strategy (The Pyramid)

#### Level 1: Unit Tests (Fast)

```bash
make test          # Run Go unit tests (with race detector)
make test/rust     # Run Rust unit tests (cargo test)
```

#### Level 2: E2E Testing (The Gold Standard)

Vald's reliability relies on **E2E V2**. This creates a real cluster, deploys Vald, inserts vectors, and verifies search results.

```bash
# 1. Prepare Kubernetes Cluster
make k3d/start

# 2. Deploy Vald (Configured for E2E)
make k8s/vald/deploy HELM_VALUES=example/helm/values.yaml

# 3. Run the E2E Test Suite
# This runs the standard scenario: Insert -> Wait -> Search -> Verify
make E2E_TIMEOUT=1h e2e/v2

# 4. Cleanup
make k8s/vald/delete
make k3d/delete
```

---

## 5. 🗺 Repository Map & Dependency Graph

- `apis/proto/v1/`: **Single API Source of Truth**. Changing this affects Go, Rust, Java, Python, Node.js SDKs.
- `Makefile` & `Makefile.d/*`: The brain of the repository.
- `cmd/`: Entry points for each component (Agent, Gateway, Discoverer, etc.).
- `internal/`: Go core libraries (shared, heavily optimized).
- `pkg/`: Go Service logic (Agent, Gateway, Discoverer).
- `rust/`: Rust implementation of core components.
- `charts/`: Helm charts. Changing Go config structures requires updating `values.yaml` here.
- `tests/v2/e2e`: Vald's E2E test V2 implementation.
- `hack/`: Scripts for benchmarking (`hack/benchmark`), dataset generation, and license checks.

---

## 6. ⚠️ Known Tooling Gap: Claude Code Stop-hook vs. 100% `//go:build e2e`-gated packages

`tests/v2/e2e/config`, `tests/v2/e2e/crud`, and `tests/v2/e2e/hdf5` consist **entirely** of
`//go:build e2e`-tagged files (zero default-tag `.go` files per package). This is a deliberate,
repo-wide convention, not an oversight - confirmed by:

- `go vet $(ROOTDIR)/...` (the project's own `go-vet` macro in `Makefile.d/functions.mk`) never
  errors on these packages: Go's `./...` wildcard silently skips a directory that matches zero
  packages under the current build tags.
- `go vet ./tests/v2/e2e/crud/...` (targeting the package directly, no `-tags e2e`) instead prints
  `no packages to vet` and exits 1 - a soft failure, not a hard typecheck error.
- `golangci-lint run --new-from-rev=HEAD ./tests/v2/e2e/crud/` (what the local Claude Code
  `swarm-stop-verify.sh` Stop-hook runs, per-directory, for any edited `.go` file) instead **hard
  fails** with `typechecking error: build constraints exclude all Go files` (exit 7). This is an
  artifact of golangci-lint's stricter handling of an explicitly-named path with zero buildable
  files, and is unrelated to any code defect - it reproduces identically on a totally clean
  checkout of `main`, for literally any edit inside these three directories, before and after any
  actual code change.

**Do not** add `build-tags: [e2e]` to the root `.golangci.json`/`.golangci.yaml` - that flips lint
scope repo-wide and would surface a large volume of pre-existing, unreviewed findings across
`config`/`crud`/`hdf5` in every future `make lint` / CI run, far outside any single task's scope.
Editing `~/.claude/hooks/swarm-stop-verify.sh` itself is also out of a Maker sub-agent's authority
(hook changes require human approval, per `SWARM.md` §5's swarm-evolve gate).

**Accepted minimal fix, scoped to whichever of these packages you are actively editing:** add a
single `doc_test.go` (matching the `_test.go` naming convention already used by every file in
these packages) with **no build tag**, containing only the package doc comment + `package <name>`
declaration - no test functions, no other symbols. `go test`'s package-vet step considers
`_test.go` files regardless of tag, so this one extra tagless file is enough to give
`golangci-lint run ./tests/v2/e2e/<pkg>/` (and any other per-directory, no-tag tooling) at least
one buildable file to typecheck, without changing what actually executes under `-tags e2e` (it
contributes zero tests) and without touching repo-wide lint config. Verified for
`tests/v2/e2e/crud` (see `doc_test.go`): `golangci-lint run --new-from-rev=HEAD
./tests/v2/e2e/crud/` goes from a hard `typechecking error` to `0 issues`, while
`go test -tags e2e ./tests/v2/e2e/crud/... -args -config <cfg>` and `golangci-lint run
--build-tags=e2e ./tests/v2/e2e/crud/...` remain unaffected/clean. Apply the same pattern to
`config`/`hdf5` only if/when a task actually touches those packages - don't proactively patch
untouched packages.

Also verify these packages directly with `go vet -tags e2e ./tests/v2/e2e/<pkg>/...` and
`golangci-lint run --build-tags=e2e ./tests/v2/e2e/<pkg>/...` (optionally `--new-from-rev=HEAD`)
as a secondary check, since those exercise the real `-tags e2e` build these packages actually run
under.

## 8. 📜 Mission Trajectory Log

2026-07-18 | tests/v2/e2e に Recall-QPS Pareto チャート生成機能を追加(recall記録・AchievedQPS・chartパッケージ・CLI) | 2試行 | done | 試行1でChecker/決定論的検証はPASSしたが、タスク仕様が要求した「スタンドアロンCLIツール」が未実装(ライブラリのみ)というギャップをオーケストレーター側の独立検証で発見。試行2でcmd/chart/main.goを追加し実バイナリで実SVG生成をObserveして解消。学び: Checkerは指定した契約(テスト)の充足は厳密に見るが、タスク仕様全体との突き合わせは別途行う必要がある。
