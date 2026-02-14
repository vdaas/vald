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

.PHONY: test

## run tests for cmd, internal, pkg
test: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/cmd/... $(ROOTDIR)/internal/... $(ROOTDIR)/pkg/...
	$(MAKE) certs/clean

.PHONY: test/tparse
## run tests for cmd, internal, pkg and show table
test/tparse: \
	certs/gen \
	tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/cmd/... $(ROOTDIR)/internal/... $(ROOTDIR)/pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests
	$(MAKE) certs/clean

.PHONY: test/cmd/tparse
## run tests for cmd and show table
test/cmd/tparse: \
	certs/gen \
	tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/cmd/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests
	$(MAKE) certs/clean

.PHONY: test/internal/tparse
## run tests for internal and show table
test/internal/tparse: \
	certs/gen \
	tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/internal/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests
	$(MAKE) certs/clean

.PHONY: test/pkg/tparse
## run tests for pkg and who table
test/pkg/tparse: \
	certs/gen \
	tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests
	$(MAKE) certs/clean

.PHONY: test/hack/tparse
## run tests for hack and show table
test/hack/tparse: \
	certs/gen \
	tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	go mod vendor -o $(ROOTDIR)/vendor
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=vendor -json -cover \
	$(ROOTDIR)/hack/gorules/... \
	$(ROOTDIR)/hack/helm/... \
	$(ROOTDIR)/hack/license/... \
	$(ROOTDIR)/hack/tools/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests
	rm -rf $(ROOTDIR)/vendor
	$(MAKE) certs/clean

.PHONY: test/all/tparse
## run tests for all Go codes and show table
test/all/tparse: \
	certs/gen \
	tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests
	$(MAKE) certs/clean

.PHONY: test/gotestfmt
## run tests for cmd, internal, pkg and show table
test/gotestfmt: \
	certs/gen \
	gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	GODEBUG=$(GODEBUG) \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) $(ROOTDIR)/cmd/... $(ROOTDIR)/internal/... $(ROOTDIR)/pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus
	$(MAKE) certs/clean

.PHONY: test/cmd/gotestfmt
## run tests for cmd and show table
test/cmd/gotestfmt: \
	certs/gen \
	gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	GODEBUG=$(GODEBUG) \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) -ldflags="-linkmode=external" $(ROOTDIR)/cmd/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus
	$(MAKE) certs/clean

.PHONY: test/internal/gotestfmt
## run tests for internal and show table
test/internal/gotestfmt: \
	certs/gen \
	gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	GODEBUG=$(GODEBUG) \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) -ldflags="-linkmode=external" $(ROOTDIR)/internal/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus -hide="all"
	$(MAKE) certs/clean

.PHONY: test/pkg/gotestfmt
## run tests for pkg and who table
test/pkg/gotestfmt: \
	certs/gen \
	gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	GODEBUG=$(GODEBUG) \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) -ldflags="-linkmode=external" $(ROOTDIR)/pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus
	$(MAKE) certs/clean

.PHONY: test/hack/gotestfmt
## run tests for hack and show table
test/hack/gotestfmt: \
	certs/gen \
	gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	go mod vendor -o $(ROOTDIR)/vendor
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	GODEBUG=$(GODEBUG) \
	go test -short -shuffle=on -race -mod=vendor -json -cover -ldflags="-linkmode=external" \
	$(ROOTDIR)/hack/gorules/... \
	$(ROOTDIR)/hack/helm/... \
	$(ROOTDIR)/hack/license/... \
	$(ROOTDIR)/hack/tools/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus
	rm -rf $(ROOTDIR)/vendor
	$(MAKE) certs/clean

.PHONY: test/all/gotestfmt
## run tests for all Go codes and show table
test/all/gotestfmt: \
	certs/gen \
	gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	GODEBUG=$(GODEBUG) \
	go test -short -shuffle=on -race -mod=readonly -json -cover -timeout=$(GOTEST_TIMEOUT) -ldflags="-linkmode=external" $(ROOTDIR)/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus
	$(MAKE) certs/clean

.PHONY: test/create-empty
## create empty test file if not exists
test/create-empty:
	@$(call green, "create empty test file if not exists...")
	@for f in $(GO_ALL_TEST_SOURCES) ; do \
		if [ ! -f "$$f" ]; then \
			echo "Creating empty test file $$f"; \
			package="$$(dirname $$f)" ; \
			package="$$(basename $$package)" ; \
			if [ "$$(basename $$f)" = "main.go" ]; then \
				package="main"; \
			fi; \
	echo "package $$package" >> "$$f"; \
	fi; \
	done

.PHONY: test/remove-empty
## remove empty test files
test/remove-empty:
	@$(call green, "remove empty test files...")
	@for f in $(GO_ALL_TEST_SOURCES) ; do \
		if ! grep -q "func Test" "$$f"; then \
			echo "Removing empty test file $$f"; \
			rm "$$f"; \
		fi; \
	done

.PHONY: test/pkg
## run tests for pkg
test/pkg: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -cover $(ROOTDIR)/pkg/...
	$(MAKE) certs/clean

.PHONY: test/internal
## run tests for internal
test/internal: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -cover $(ROOTDIR)/internal/...
	$(MAKE) certs/clean

.PHONY: test/cmd
## run tests for cmd
test/cmd: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -cover $(ROOTDIR)/cmd/...
	$(MAKE) certs/clean

.PHONY: test/rust
## run tests for rust
test/rust: \
	test/rust/qbg \
	test/rust/kvs \
	test/rust/vqueue \
	test/rust/observability \
	test/rust/agent

.PHONY: test/rust/qbg
## run tests for qbg crate
test/rust/qbg:
	cargo test --manifest-path rust/Cargo.toml --package qbg --lib -- --show-output

.PHONY: test/rust/kvs
## run tests for kvs crate
test/rust/kvs:
	cargo test --manifest-path rust/Cargo.toml --package kvs --lib -- --show-output

.PHONY: test/rust/vqueue
## run tests for vqueue crate
test/rust/vqueue:
	cargo test --manifest-path rust/Cargo.toml --package vqueue --lib -- --show-output

.PHONY: test/rust/observability
## run tests for observability crate
test/rust/observability:
	cargo test --manifest-path rust/Cargo.toml --package observability --lib -- --show-output

.PHONY: test/rust/agent
## run tests for agent
test/rust/agent:
	cargo test --manifest-path rust/Cargo.toml --package agent -- --show-output

.PHONY: test/hack
## run tests for hack
test/hack: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	go mod vendor -o $(ROOTDIR)/vendor
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=vendor -cover \
	$(ROOTDIR)/hack/gorules/... \
	$(ROOTDIR)/hack/helm/... \
	$(ROOTDIR)/hack/license/... \
	$(ROOTDIR)/hack/tools/...
	rm -rf $(ROOTDIR)/vendor
	$(MAKE) certs/clean

.PHONY: test/all
## run tests for all Go codes
test/all: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -cover $(ROOTDIR)/...
	$(MAKE) certs/clean

.PHONY: coverage
## calculate coverages
coverage: \
	coverage/go \
	coverage/rust

.PHONY: coverage/go
## calculate go coverages
coverage/go: certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test -short -shuffle=on -race -mod=readonly -v -race -covermode=atomic -timeout=$(GOTEST_TIMEOUT) -coverprofile=coverage.out $(ROOTDIR)/...
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go tool cover -html=coverage.out -o coverage.html
	$(MAKE) certs/clean

.PHONY: coverage/rust
## calculate rust coverages
coverage/rust:
	cargo llvm-cov --manifest-path rust/Cargo.toml --workspace --exclude proto --lcov --output-path rust-coverage.out

.PHONY: gotests/gen
## generate missing go test files
gotests/gen:
	$(MAKE) test/create-empty
	$(MAKE) test/patch-placeholder
	$(MAKE) gotests/gen-test
	$(MAKE) test/remove-empty
	$(MAKE) gotests/patch
	$(MAKE) format/go/test
	$(MAKE) test/comment-unimplemented
	$(MAKE) format/go/test

.PHONY: gotests/gen-test
## generate test implementation
gotests/gen-test:
	@$(call green, "generate go test files...")
	$(call gen-go-test-sources)
	$(call gen-go-option-test-sources)

.PHONY: gotests/patch
## apply patches to generated go test files
gotests/patch: \
	files
	@$(call green, "apply patches to go test files...")
	find $(ROOTDIR)/internal/k8s/* -name '*_test.go' | xargs -P$(CORES) sed -i -E "s%k8s.io/apimachinery/pkg/api/errors%$(GOPKG)/internal/errors%g"
	@cat $(ROOTDIR)/.gitfiles | grep -E '^(\./)?internal/k8s/.*\_test.go$$' | xargs -I {} -P$(CORES) bash -c ' \
	echo "Replacing internal/k8s Test File {}" && \
	sed -i -E "s%cockroachdb/errors%$(REPO)/internal/errors%g" {} && \
	sed -i -E "s%golang.org/x/sync/errgroup%$(GOPKG)/internal/sync/errgroup%g" {} && \
	sed -i -E "s%pkg/errors%$(REPO)/internal/errors%g" {} && \
	sed -i -E "s%go-errors/errors%$(REPO)/internal/errors%g" {} && \
	sed -i -E "s%go.uber.org/goleak%$(GOPKG)/internal/test/goleak%g" {}'
	@cat $(ROOTDIR)/.gitfiles | grep -E '^(\./)?internal/errors/.*\_test.go$$' | xargs -I {} -P$(CORES) bash -c ' \
	echo "Replacing internal/errors Test {}" && \
	sed -i -E "s%\"$(GOPKG)/internal/errors\"%%g" {} && \
	sed -i -E "s/errors\.//g" {}'
	@cat $(ROOTDIR)/.gitfiles | grep -E '^(\./)?internal/test/goleak/.*\_test.go$$' | xargs -I {} -P$(CORES) bash -c ' \
	echo "Replacing goleak Test file {}" && \
	sed -i -E "s%\"$(GOPKG)/internal/test/goleak\"%%g" {} && \
	sed -i -E "s/goleak\.//g" {}'
	@sed -i -E '/"github.com\/vdaas\/vald\/internal\/strings"/d' pkg/gateway/lb/handler/grpc/pairing_heap_test.go
	@sed -i -E '/^\t"github.com\/vdaas\/vald\/internal\/k8s\/vald"$$/d' pkg/agent/core/ngt/service/ngt_test.go

.PHONY: test/patch-placeholder
## delete from placeholder to EOF, then re-append placeholder
test/patch-placeholder:
	@$(call green, "apply placeholder patches (delete to EOF and append)...")
	@for f in $(GO_ALL_TEST_SOURCES); do echo "$$f"; done | \
		xargs -I {} -P$(CORES) bash -c ' \
		[ -f "{}" ] || exit 0 ; \
		sed -i -E "/\/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/,\$$d" "{}"; \
		echo "// $(TEST_NOT_IMPL_PLACEHOLDER)" >> {};'

.PHONY: test/comment-unimplemented
## comment out unimplemented tests (from placeholder to EOF)
test/comment-unimplemented:
	@$(call green, "comment out unimplemented test (from placeholder to EOF)...")
	@for f in $(GO_ALL_TEST_SOURCES); do echo "$$f"; done | \
		xargs -I {} -P$(CORES) bash -c ' \
		[ -f "{}" ] || exit 0 ; \
		sed -i -E -e " \
		/\/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/,\$$ { \
		s/^/\/\/ /; \
		s/^\/\/ \/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/\/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/; \
		}" "{}"'
