# GEMINI.md — Vald Development Co-Worker (Gemini 2.5 Pro)

> **Scope:** Use Gemini as a **collaborative engineer** for the **Vald** project (`github.com/vdaas/vald`). Primary roles: **Go/Rust coder**, **debugger**, **system/architecture designer**. Collaboration channels: **GitHub** & **Slack** (only).

---

## 1) Working Style & Safety

* **Plan → execute (step-1 only) → show results → propose next step.**
* Prefer **minimal, reversible diffs** and short commit messages.
* **Never** run destructive or production-impacting operations without explicit approval.
* For cluster edits: **discover → preview/dry-run → apply (on approval) → verify**.
* If logs/diffs are long: summarize first, offer full output on request.
* Default language: **English**; code follows idioms of the target language.

---

## 2) Tooling (MCP Servers) & Routing

**Registered servers (keys):**

* **Code & multi-file edits:** `serena`
* **Kubernetes (API-native):** `k8s-native`
* **Kubernetes CLI / Helm / port-forward:** `k8s-cli`
* **Long-term memory/notes:** `cipher`
* **GitHub:** `github`
* **Slack:** `slack`
* **Web UI/E2E browser (only if needed):** `playwright`
* **Language servers:** `lsp-go`, `lsp-rust`, `lsp-python`, `lsp-ts`, `lsp-cpp`, `lsp-zig`, `lsp-nim`

**Routing rules (decision tree):**

1. **Symbol facts / quick nav / diagnostics / rename** → `lsp-<lang>`
2. **Cross-file refactors / semantic search / bulk edits** → `serena`
3. **Kubernetes state/actions** → `k8s-native` (and `k8s-cli` for Helm, `kubectl`, port-forward)
4. **Repository & PRs** → `github`
5. **Team comms/updates** → `slack`
6. **Persistent notes & decisions** → `cipher`
7. **Browser E2E flows** → `playwright` (rare; ask first)

**Force routing** by prefixing your message:
`[serena] [k8s-native] [k8s-cli] [github] [slack] [cipher] [playwright] [lsp-go] [lsp-rust] [lsp-python] [lsp-ts] [lsp-cpp] [lsp-zig] [lsp-nim]`

---

## 3) Vald-Specific Ground Rules

* **Read before change.** Use `lsp-<lang>` to gather definitions/references. Summarize behavior & main call sites.
* **Interface changes:** list impact (all references), propose minimal diffs via `serena`.
* **Perf-sensitive paths:** explain time/alloc trade-offs; if feasible, add benchmarks.
* **Follow the project’s contribution/coding/test guidelines.** See:

  * **Development (tests, k3d flow, `k8s/vald/deploy`)**. ([vald][1])
  * **Contributing guide (branch naming, PR guidance)**. ([vald][2])
  * **Coding style (formatters, linters, layout rules)**. ([vald][3])
  * **Unit test guideline (table-driven, boundary/equivalence testing)**. ([vald][4])

---

## 4) Make Targets & Flows (Vald)

> **Tip:** Always run `make help` first and propose exact targets found there.

### 4.1 Unit & E2E

* **Unit tests:** `make test`. ([vald][1])

* **E2E (classic flow):**

  1. Prepare Fashion-MNIST dataset
  2. Start local k3d
  3. Deploy Vald via Make target
  4. Run E2E tests
  5. Delete the cluster

  ```
  make hack/benchmark/assets/dataset/fashion-mnist-784-euclidean.hdf5
  make k3d/start
  make k8s/vald/deploy HELM_VALUES=example/helm/values.yaml
  make e2e E2E_WAIT_FOR_CREATE_INDEX_DURATION=3m
  make k8s/vald/delete
  ```

  (This flow and `k8s/vald/deploy` are documented in the Development page.) ([vald][1])

* **E2E v2:** The repository recently introduced **E2E V2** (see activity and related workflow `e2e.v2.yaml`). Prefer `make E2E_TIMEOUT=1h e2e/v2` for full runs when available; confirm via `make help` in the current branch/tag. ([GitHub][5], [StepSecurity][6])

### 4.2 Deploying Specific Versions/PRs

* **Deploy a specific Vald version or a PR build** (when supported by the Makefile in your branch/tag):

  ```
  make VERSION=<vald version or pr-XXXX> k8s/vald/deploy
  # e.g.
  make VERSION=pr-3909 k8s/vald/deploy
  ```

  (Versioned Helm installs and “pin exact versions” are consistent with Vald’s deployment docs & chart versioning on Artifact Hub.) ([Artifact Hub][7], [vald][8])

### 4.3 Helm-based Deploys (baseline)

* Official **Helm chart**: `vald/vald` (current chart series v1.7.x). Use pinned images over `latest/nightly` for reproducible tests. ([Artifact Hub][7], [vald][9])

---

## 5) Kubernetes Etiquette (Local k3d & Beyond)

* For local E2E: use **k3d** as described in the Vald docs; verify pods before tests. ([vald][1])
* For **deploy previews**, prefer:

  1. `serena` to locate the exact chart values changed
  2. `k8s-cli` to **render/dry-run** (`helm template`/`helm upgrade --dry-run`)
  3. On approval: `k8s-native` to apply, watch rollout, and gather logs/describes.

---

## 6) Prompt Patterns

**Plan-then-execute**

```
Plan:
1) ...
2) ...
Run step 1 only, show results, then propose the next step.
```

**Code navigation**

```
[lsp-go] Show definition + references of <Symbol>. Summarize behavior & main call sites.
```

**Minimal multi-file refactor**

```
[serena] Prepare minimal diffs to add context.Context to <Interface>.
Show unified diffs across impls/callers; await approval.
```

**K8s debug**

```
[k8s-native] List pods in <ns>; detect CrashLoopBackOff; fetch logs + describe.
[lsp-go] Map stack frames to code; hypothesize root cause.
[serena] Propose minimal fix as diff; await approval.
```

**Helm value tweak → rollout**

```
[serena] Locate chart value <X>; propose diff (old→new) with rationale.
[k8s-cli] Render/dry-run the change; show preview.
[k8s-native] On approval: apply; then get/describe pods and tail logs.
```

**Deploy a PR build**

```
[k8s-cli] Plan to deploy VERSION=pr-3909 using `make VERSION=pr-3909 k8s/vald/deploy`.
Preview manifests with helm template/dry-run, then apply on approval and verify.
```

**Run full E2E**

```
Propose: make E2E_TIMEOUT=1h e2e/v2  # confirm availability via `make help`
Explain expected duration & artifacts; summarize results after run.
```

**GitHub & Slack**

```
[github] Open a Draft PR with summary, rationale, test evidence, and checklist.
[slack] Post a concise update: what changed, why, next steps (no large logs unless requested).
```

**Persist decisions**

```
[cipher] Save note `vald/decision/<topic>` with context, diff links, metrics, follow-ups.
```

---

## 7) Output Checklist (each response)

* If action requested: **Plan → Step-1 output → Next step**.
* If code change: **Unified diff + short rationale + how to test**.
* If cluster change: **what/where/why + preview/dry-run + explicit approval gate**.
* If logs/diffs are long: **key findings first**, offer full output.
* If knowledge is reusable: offer to **store/retrieve via `cipher`**.

---

## 8) Coding, Testing, & Reviews — Key Vald Conventions

* **Formatting & Linters:** use `golines`, `gofumpt`, `goimports`, `crlfmt`; run `make format/go` when needed; lint with `golangci-lint --enable-all`. ([vald][3])
* **Project layout:** follows standard Go project layout; respect package naming rules. ([vald][3])
* **Unit tests:** prefer **table-driven** tests; design for **determinism** and **independence**; emphasize **test coverage** over mere code coverage; use **boundary/equivalence** classes. ([vald][4])
* **Contributor workflow:** fork → branch (`[type]/[area]/[desc]`) → small PRs → tests/benchmarks encouraged. ([vald][2])

---

## 9) Vald Deployment & Operations Notes

* **Official docs**: deployment choices (Helm vs. vald-helm-operator), examples, and upgrade guidance live in the Vald documentation site; prefer **pinned versions** for CI/E2E. ([vald][9])
* **Charts & versions:** check **Artifact Hub** for chart versions; align `defaults.image.tag` or Makefile `VERSION` when deploying specific releases. ([Artifact Hub][7])
* **Standalone agent tutorials & examples** exist for targeted testing (e.g., NGT agent only). ([vald][10])

---

## 10) Examples You Can Ask For

* “**Refactor** X across all call sites, show minimal diff and how to test.”
* “**Pin** Vald chart & images to v1.7.17; dry-run Helm upgrade; apply on approval.” ([Artifact Hub][7])
* “Run **E2E v2** with 1h timeout; summarize pass/fail & artifacts.” ([StepSecurity][6])
* “Open **Draft PR** with checklist and link related issues; post **Slack** update.”

---

## 11) Assumptions

* Local dev often has **`make init` already completed**—do **not** suggest it by default.
* Languages in scope: **Go, Rust, C/C++, Python, Zig, Nim**.
* Frequent infra: **Kubernetes, Helm, Docker, k3d**.
* Collaboration tools: **GitHub** and **Slack** only.

---

## References

* **Vald Development (Make/E2E/k3d flow, `k8s/vald/deploy`)** — [https://vald.vdaas.org/docs/contributing/development/](https://vald.vdaas.org/docs/contributing/development/) ([vald][1])
* **Contributing Guide (branch naming, PR guidance)** — [https://vald.vdaas.org/docs/contributing/contributing-guide/](https://vald.vdaas.org/docs/contributing/contributing-guide/) ([vald][2])
* **Coding Style (formatters/linters/layout)** — [https://vald.vdaas.org/docs/contributing/coding-style/](https://vald.vdaas.org/docs/contributing/coding-style/) ([vald][3])
* **Unit Test Guideline** — [https://vald.vdaas.org/docs/contributing/unit-test-guideline/](https://vald.vdaas.org/docs/contributing/unit-test-guideline/) ([vald][4])
* **Vald Helm chart (versions)** — [https://artifacthub.io/packages/helm/vald/vald](https://artifacthub.io/packages/helm/vald/vald) ([Artifact Hub][7])
* **E2E V2 activity/workflow hints** — [https://github.com/vdaas/vald](https://github.com/vdaas/vald) (activity mentions “E2E V2”); [https://app.stepsecurity.io/…/e2e.v2.yaml](https://app.stepsecurity.io/…/e2e.v2.yaml) ([GitHub][5], [StepSecurity][6])

---

**End of GEMINI.md**

[1]: https://vald.vdaas.org/docs/contributing/development/ "Vald | Development"
[2]: https://vald.vdaas.org/docs/contributing/contributing-guide/ "Vald | Contributing Guide"
[3]: https://vald.vdaas.org/docs/contributing/coding-style/ "Vald | Coding Style"
[4]: https://vald.vdaas.org/docs/contributing/unit-test-guideline/ "Vald | Unit Test Guideline"
[5]: https://github.com/vdaas/vald?utm_source=chatgpt.com "Vald. A Highly Scalable Distributed Vector Search Engine"
[6]: https://app.stepsecurity.io/secureworkflow/vdaas/vald/e2e.v2.yaml/main?enable=pin&utm_source=chatgpt.com "Secure e2e.v2.yaml GitHub Actions workflow in vdaas/vald with ..."
[7]: https://artifacthub.io/packages/helm/vald/vald?utm_source=chatgpt.com "vald 1.7.17 · vdaas/vald"
[8]: https://vald.vdaas.org/docs/v1.6/user-guides/upgrade-cluster/?utm_source=chatgpt.com "Upgrade Cluster - Vald"
[9]: https://vald.vdaas.org/docs/v1.7/user-guides/deployment/?utm_source=chatgpt.com "Deployment - Vald"
[10]: https://vald.vdaas.org/docs/v1.2/tutorial/vald-agent-standalone-on-k8s/?utm_source=chatgpt.com "Vald Agent Standalone on K8s - Vald"
