#
# Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

.PHONY: lint
## run lints
lint: \
	shell/options/check \
	options/check \
	docs/lint \
	files/lint \
	vet \
	go/lint \
	workflow/lint

# --- shell option policy ------------------------------------------------------
# The canonical option set is owned by the root Makefile (`SHELL_STRICT_OPTS`).
# The values below are what hack/docker/gen renders into the generated
# Dockerfiles; the check fails when either side of that contract drifts.
DOCKER_SHELL_LINE := SHELL ["/bin/bash", "$(subst $() ,"$(COMMA) ",$(SHELL_STRICT_OPTS))", "-c"]
DOCKER_RUN_SET = $(SHELL_STRICT) -x

# The scan patterns live in variables so that this file itself never matches them
# (`[[:alpha:]]` / `[o]` spellings keep the literal text non-matching). This is why
# Makefile.d/lint.mk no longer has to be excluded from the scan.
SHELL_ADHOC_RE := (^|[^[:alnum:]_-])(set -[[:alpha:]]|(bash|sh) -c)
SHELL_PIPEFAIL_RE := ^[[:space:]]*set -[o] pipefail

.PHONY: shell/options/check
## verify every shell in the repository uses the canonical option set
shell/options/check:
	@echo "Checking canonical shell options ($(SHELL_STRICT_OPTS))..."
	@fail=0; \
	if git -C $(ROOTDIR) grep -nE '$(SHELL_ADHOC_RE)' -- Makefile Makefile.d; then \
		echo "ERROR: use \$$(SHELL_STRICT) or \$$(SHELL_STRICT_EXEC) instead of ad-hoc shell options (single source of truth: SHELL_STRICT_OPTS in Makefile)" >&2; \
		fail=1; \
	fi; \
	if git -C $(ROOTDIR) grep -nE '$(SHELL_PIPEFAIL_RE)' -- .github/workflows; then \
		echo "ERROR: pipefail is provided by defaults.run.shell: bash; drop this line" >&2; \
		fail=1; \
	fi; \
	missing=$$( \
		git -C $(ROOTDIR) grep -lE '^[[:space:]]+run:' -- .github/workflows \
		| while read -r f; do \
			git -C $(ROOTDIR) grep -qE '^[[:space:]]+shell: bash$$' -- "$$f" || echo "$$f"; \
		  done \
	); \
	if [ -n "$$missing" ]; then \
		echo "ERROR: workflow with run: steps is missing defaults.run.shell: bash (pipefail): $$missing" >&2; \
		fail=1; \
	fi; \
	bad=$$(git -C $(ROOTDIR) grep -nE '^SHELL ' -- dockers | grep -vF '$(DOCKER_SHELL_LINE)' || true); \
	if [ -n "$$bad" ]; then \
		echo "ERROR: generated Dockerfile SHELL instruction must match SHELL_STRICT_OPTS (run make dockerfile): $$bad" >&2; \
		fail=1; \
	fi; \
	bad=$$(git -C $(ROOTDIR) grep -nE '^[[:space:]]*set ' -- dockers | grep -vF '$(DOCKER_RUN_SET)' || true); \
	if [ -n "$$bad" ]; then \
		echo "ERROR: generated Dockerfile RUN must start with the canonical set instruction (run make dockerfile): $$bad" >&2; \
		fail=1; \
	fi; \
	for f in .devcontainer/postAttachCommand.sh hack/git/hooks/pre-commit; do \
		grep -qxF '$(SHELL_STRICT)' "$(ROOTDIR)/$$f" \
		|| { echo "ERROR: $$f must declare '$(SHELL_STRICT)'" >&2; fail=1; }; \
	done; \
	exit $$fail

# --- shared command option policy --------------------------------------------
# The per-command option sets live in the root Makefile (CURL*, XARGS,
# GOLANGCI_LINT_ARGS, KUBECTL_WAIT_*, K8S_WAIT_TIMEOUT*, DOCKER_*_ARGS) and in
# hack/docker/gen for generated images. Call sites must not re-declare them.
UBUNTU_VERSION_FILE := versions/UBUNTU_VERSION
DOCKER_APT_INSTALL_LINE := apt-get install -y --no-install-recommends --fix-missing

# --- canonical option variables ----------------------------------------------
# Values that recipes depend on. Each must expand without error (a definition such
# as `X ?= $(X)` makes make abort with "Recursive variable references itself") and
# must not be empty. This is the regression gate for the class of bug where a
# blanket substitution rewrote a variable's own definition.
CANONICAL_OPTION_VARS := \
	CURL CURL_AUTH CURL_RETRY CURL_RETRY_ALL CURL_RETRY_SHORT CURL_DOWNLOAD \
	XARGS GOLANGCI_LINT_ARGS GIT_DIFF_FILES \
	KUBECTL_WAIT_READY KUBECTL_WAIT_AVAILABLE KUBECTL_WAIT_CONTAINERS_READY KUBECTL_WAIT_COMPLETE KUBECTL_WAIT_DELETE \
	K8S_WAIT_TIMEOUT K8S_WAIT_TIMEOUT_SHORT K8S_WAIT_TIMEOUT_MID K8S_WAIT_TIMEOUT_LONG \
	DOCKER_TAG_ARGS DOCKER_NETWORK_ARGS \
	SHELL_STRICT_OPTS SHELL_STRICT SHELL_STRICT_EXEC

.PHONY: options/vars/print
## print one canonical option variable (helper for options/vars/check)
options/vars/print:
	@printf '%s|%s' '$(origin $(VAR))' '$($(VAR))'

.PHONY: options/vars/check
## verify every canonical option variable expands to a non-empty value
options/vars/check:
	@echo "Checking canonical option variables..."
	@fail=0; \
	for v in $(CANONICAL_OPTION_VARS); do \
		if ! out=$$($(MAKE) --no-print-directory -s VAR="$$v" options/vars/print 2>/dev/null); then \
			echo "ERROR: $$v does not expand (self-reference or syntax error in its definition)" >&2; \
			fail=1; \
			continue; \
		fi; \
		origin=$${out%%|*}; \
		val=$${out#*|}; \
		# environment overrides are intentional; only make them visible instead of silent \
		[ "$$origin" = environment ] && echo "note: $$v overridden by the environment: '$$val'"; \
		if [ -z "$$val" ]; then \
			echo "ERROR: $$v expands to an empty value" >&2; \
			fail=1; \
			continue; \
		fi; \
		shape_fail() { echo "ERROR: $$v=$$1 does not look like $$2 (origin: $$origin)" >&2; fail=1; }; \
		case "$$v" in \
			CURL|CURL_AUTH|CURL_DOWNLOAD) printf '%s' "$$val" | grep -q 'curl' || shape_fail "$$val" "a curl invocation" ;; \
			CURL_RETRY|CURL_RETRY_SHORT) printf '%s' "$$val" | grep -q -- '--retry' || shape_fail "$$val" "curl retry flags" ;; \
			CURL_RETRY_ALL) printf '%s' "$$val" | grep -q -- '--retry-all-errors' || shape_fail "$$val" "curl retry-all-errors flags" ;; \
			XARGS) printf '%s' "$$val" | grep -q 'xargs' || shape_fail "$$val" "an xargs invocation" ;; \
			GOLANGCI_LINT_ARGS) printf '%s' "$$val" | grep -q -- '--config' || shape_fail "$$val" "a golangci-lint --config flag" ;; \
			GIT_DIFF_FILES) printf '%s' "$$val" | grep -q 'git diff' || shape_fail "$$val" "a git diff invocation" ;; \
			KUBECTL_WAIT_*) printf '%s' "$$val" | grep -q -- '--for=' || shape_fail "$$val" "a kubectl --for= condition" ;; \
			K8S_WAIT_TIMEOUT*) printf '%s' "$$val" | grep -qE '^[0-9]' || shape_fail "$$val" "a duration such as 600s" ;; \
			DOCKER_TAG_ARGS) printf '%s' "$$val" | grep -q -- '-t ' || shape_fail "$$val" "at least one -t tag" ;; \
			DOCKER_NETWORK_ARGS) printf '%s' "$$val" | grep -q -- '--network' || shape_fail "$$val" "a --network flag" ;; \
			SHELL_STRICT_OPTS) printf '%s' "$$val" | grep -q 'pipefail' || shape_fail "$$val" "shell options containing pipefail" ;; \
			SHELL_STRICT) printf '%s' "$$val" | grep -q '^set ' || shape_fail "$$val" "a set builtin invocation" ;; \
			SHELL_STRICT_EXEC) printf '%s' "$$val" | grep -q -- ' -c' || shape_fail "$$val" "an interactive shell -c invocation" ;; \
		esac; \
		case "$$val" in \
			*'$$('*|*'  '*) \
			echo "ERROR: $$v has an unresolved or blank substitution: '$$val' (an eagerly evaluated 'VAR :=' before its dependency is defined looks like this)" >&2; \
			fail=1 ;; \
		esac; \
	done; \
	exit $$fail

.PHONY: options/check
## verify shared command options are not re-declared at call sites
options/check: options/vars/check
	@echo "Checking shared command options..."
	@fail=0; \
	for f in $$(git -C $(ROOTDIR) grep -l 'uses: actions/checkout@' -- .github/workflows); do \
		uses=$$(git -C $(ROOTDIR) grep -hcE 'uses: actions/checkout@' -- "$$f"); \
		persist=$$(git -C $(ROOTDIR) grep -hcE 'persist-credentials: false' -- "$$f"); \
		if [ "$$uses" != "$$persist" ]; then \
			echo "ERROR: $$f has $$uses actions/checkout but only $$persist 'persist-credentials: false' (credentials must never be persisted)" >&2; \
			fail=1; \
		fi; \
	done; \
	for f in $$(git -C $(ROOTDIR) ls-files '.github/workflows/*.yaml' '.github/workflows/*.yml'); do \
		runs=$$(git -C $(ROOTDIR) grep -hcE '^[[:space:]]+runs-on:' -- "$$f" || true); \
		tmo=$$(git -C $(ROOTDIR) grep -hcE '^[[:space:]]+timeout-minutes:' -- "$$f" || true); \
		if [ "$$runs" != "$$tmo" ]; then \
			echo "ERROR: $$f has $$runs job(s) but $$tmo timeout-minutes (every job must set one)" >&2; \
			fail=1; \
		fi; \
	done; \
	bad=$$(git -C $(ROOTDIR) grep -n 'apt-get install' -- dockers | grep -vF "$(DOCKER_APT_INSTALL_LINE)" || true); \
	if [ -n "$$bad" ]; then \
		echo "ERROR: generated Dockerfile apt-get flags must be \"$(DOCKER_APT_INSTALL_LINE)\" (run make dockerfile): $$bad" >&2; \
		fail=1; \
	fi; \
	bad=$$(git -C $(ROOTDIR) grep -nE '^[[:space:]]*(&&[[:space:]]+)?make .*(GOPATH|GOROOT|GO_VERSION)=' -- dockers | grep -vF 'GOPATH="$${GOPATH}" GOROOT="$${GOROOT}" GO_VERSION="$${GO_VERSION}"' || true); \
	if [ -n "$$bad" ]; then \
		echo "ERROR: generated Dockerfile must use the shared goToolchainMakeArgs triple (run make dockerfile): $$bad" >&2; \
		fail=1; \
	fi; \
	uv=$$(cat "$(ROOTDIR)/$(UBUNTU_VERSION_FILE)"); \
	if ! grep -qE "^FROM ubuntu:$$uv( |$$)" "$(ROOTDIR)/dockers/dev/Dockerfile"; then \
		echo "ERROR: dockers/dev/Dockerfile builder base must be ubuntu:$$uv from $(UBUNTU_VERSION_FILE) (run make dockerfile)" >&2; \
		fail=1; \
	fi; \
	for f in $$(git -C $(ROOTDIR) ls-files ':(glob)dockers/**/Dockerfile'); do \
		for name in $$(grep -oE '^ARG[[:space:]]+[A-Za-z_][A-Za-z0-9_]*' "$(ROOTDIR)/$$f" | awk '{print $$2}'); do \
			case "$$name" in BUILDKIT_SBOM_SCAN_STAGE|BUILDKIT_SBOM_SCAN_CONTEXT) continue ;; esac; \
				grep -qF "\$${$$name}" "$(ROOTDIR)/$$f" || grep -qE "\\$$$${$$name}([^A-Za-z0-9_]|$$)" "$(ROOTDIR)/$$f" \
				|| { echo "ERROR: $$f declares ARG $$name but never expands it (dead declaration: generated args follow the instructions that use them)" >&2; fail=1; }; \
			done; \
	done; \
	extra_args=$$(git -C $(ROOTDIR) grep -hoE -- '--build-arg [A-Za-z_][A-Za-z0-9_]*' -- Makefile.d/docker.mk \
		| sed 's/--build-arg //' | sort -u); \
	for name in $$extra_args; do \
		git -C $(ROOTDIR) grep -qE "(ARG|ENV)[[:space:]]+$$name([[:space:]=]|$$)" -- ':(glob)dockers/**/Dockerfile' \
		|| { echo "ERROR: $$name is passed by Makefile.d/docker.mk but declared by no Dockerfile" >&2; fail=1; }; \
	done; \
	exit $$fail

.PHONY: go/lint
## run golangci-lint
go/lint:
	$(call go-lint)

.PHONY: vet
## run go vet
vet:
	$(call go-vet)

.PHONY: docs/lint
## run lint for document
docs/lint: \
	docs/cspell \
	docs/textlint

.PHONY: files/lint
## run lint for document
files/lint: \
	files/cspell \
	files/textlint

.PHONY: workflow/lint workflow/fix actionlint/lint ghalint/lint
## run lint for workflow files
workflow/lint:
	@echo "Please run make workflow/fix beforehand"
	@echo "Linting workflow files..."
	@printf '%s\0' \
	"actionlint/lint" \
	"ghalint/lint" \
	| xargs -0 -I{} -P$(CORES) $(MAKE) --no-print-directory {}
	@echo "Workflow linting completed."

## run lint for workflow files
workflow/fix:
	@$(MAKE) --no-print-directory pinact/lint
	@$(MAKE) --no-print-directory ghatm/lint

ACTIONLINT_IGNORES = \
	-ignore 'when a reusable workflow is called with "uses", "timeout-minutes" is not available' \
	-ignore 'property "tag" is not defined in object type' \
	-ignore 'input "file" is not defined in action "codecov/codecov-action@v5"' \
	-ignore 'label "ubuntu-slim" is unknown.' # TODO: remove this line after https://github.com/rhysd/actionlint/issues/587 is merged

.PHONY: actionlint/lint
## run actionlint
actionlint/lint: actionlint/install
	@$(GOBIN)/actionlint -shellcheck= $(ACTIONLINT_IGNORES)

.PHONY: ghalint/lint
## run ghalint
ghalint/lint: \
	ghalint/install
	@$(GOBIN)/ghalint run .github/workflows

.PHONY: pinact/lint
## run pinact
pinact/lint: \
	pinact/install
	@GITHUB_TOKEN=$(shell gh auth token 2>/dev/null || :) $(GOBIN)/pinact run

.PHONY: ghatm/lint
## run ghatm
ghatm/lint: \
	ghatm/install
	@$(GOBIN)/ghatm set

.PHONY: docs/textlint
## run textlint for document
docs/textlint: \
	textlint/install
	bunx textlint $(ROOTDIR)/docs/**/*.md $(TEXTLINT_EXTRA_OPTIONS)

.PHONY: files/textlint
## run textlint for document
files/textlint: \
	files \
	textlint/install
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then bunx textlint "$(ROOTDIR)/.gitfiles" $(TEXTLINT_EXTRA_OPTIONS); fi

.PHONY: docs/cspell
## run cspell for document
docs/cspell: \
	cspell/install
	bunx cspell $(ROOTDIR)/docs/**/*.md --show-suggestions $(CSPELL_EXTRA_OPTIONS)

.PHONY: files/cspell
## run cspell for document
files/cspell: \
	files \
	cspell/install
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then bunx cspell "$(ROOTDIR)/.gitfiles" --show-suggestions $(CSPELL_EXTRA_OPTIONS); fi
