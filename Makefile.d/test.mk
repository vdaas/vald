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
test: \
	ngt/install \
	hdf5/install \
	certs/gen
	$(call go-test,$(ROOTDIR)/cmd/... $(ROOTDIR)/internal/... $(ROOTDIR)/pkg/...)
	$(MAKE) certs/clean

.PHONY: test/tparse
## run tests for cmd, internal, pkg and show table
test/tparse: \
	ngt/install \
	hdf5/install \
	certs/gen \
	tparse/install
	$(call go-test-tparse,$(ROOTDIR)/cmd/... $(ROOTDIR)/internal/... $(ROOTDIR)/pkg/...,$@)
	$(MAKE) certs/clean

.PHONY: test/cmd/tparse
## run tests for cmd and show table
test/cmd/tparse: \
	ngt/install \
	hdf5/install \
	certs/gen \
	tparse/install
	$(call go-test-tparse,$(ROOTDIR)/cmd/...,$@)
	$(MAKE) certs/clean

.PHONY: test/internal/tparse
## run tests for internal and show table
test/internal/tparse: \
	ngt/install \
	hdf5/install \
	certs/gen \
	tparse/install
	$(call go-test-tparse,$(ROOTDIR)/internal/...,$@)
	$(MAKE) certs/clean

.PHONY: test/pkg/tparse
## run tests for pkg and who table
test/pkg/tparse: \
	ngt/install \
	hdf5/install \
	certs/gen \
	tparse/install
	$(call go-test-tparse,$(ROOTDIR)/pkg/...,$@)
	$(MAKE) certs/clean

.PHONY: test/hack/tparse
## run tests for hack and show table
test/hack/tparse: \
	certs/gen \
	tparse/install
	GOPRIVATE=$(GOPRIVATE) \
	go mod vendor -o $(ROOTDIR)/vendor
	$(call go-test-tparse,-mod=vendor $(ROOTDIR)/hack/gorules/... $(ROOTDIR)/hack/helm/... $(ROOTDIR)/hack/license/... $(ROOTDIR)/hack/tools/...,$@)
	rm -rf $(ROOTDIR)/vendor
	$(MAKE) certs/clean

.PHONY: test/all/tparse
## run tests for all Go codes and show table
test/all/tparse: \
	ngt/install \
	hdf5/install \
	certs/gen \
	tparse/install
	$(call go-test-tparse,$(ROOTDIR)/...,$@)
	$(MAKE) certs/clean

.PHONY: test/gotestfmt
## run tests for cmd, internal, pkg and show table
test/gotestfmt: \
	ngt/install \
	hdf5/install \
	certs/gen \
	gotestfmt/install
	$(call go-test-gotestfmt,$(ROOTDIR)/cmd/... $(ROOTDIR)/internal/... $(ROOTDIR)/pkg/...,$@,-showteststatus)
	$(MAKE) certs/clean

.PHONY: test/cmd/gotestfmt
## run tests for cmd and show table
test/cmd/gotestfmt: \
	ngt/install \
	hdf5/install \
	certs/gen \
	gotestfmt/install
	$(call go-test-gotestfmt,$(ROOTDIR)/cmd/...,$@,-showteststatus -hide="all")
	$(MAKE) certs/clean

.PHONY: test/internal/gotestfmt
## run tests for internal and show table
test/internal/gotestfmt: \
	ngt/install \
	hdf5/install \
	certs/gen \
	gotestfmt/install
	$(call go-test-gotestfmt,$(ROOTDIR)/internal/...,$@,-showteststatus -hide="all")
	$(MAKE) certs/clean

.PHONY: test/pkg/gotestfmt
## run tests for pkg and who table
test/pkg/gotestfmt: \
	ngt/install \
	hdf5/install \
	certs/gen \
	gotestfmt/install
	$(call go-test-gotestfmt,$(ROOTDIR)/pkg/...,$@,-showteststatus -hide="all")
	$(MAKE) certs/clean

.PHONY: test/hack/gotestfmt
## run tests for hack and show table
test/hack/gotestfmt: \
	certs/gen \
	gotestfmt/install
	GOPRIVATE=$(GOPRIVATE) \
	go mod vendor -o $(ROOTDIR)/vendor
	$(call go-test-gotestfmt,-mod=vendor $(ROOTDIR)/hack/gorules/... $(ROOTDIR)/hack/helm/... $(ROOTDIR)/hack/license/... $(ROOTDIR)/hack/tools/...,$@,-showteststatus)
	rm -rf $(ROOTDIR)/vendor
	$(MAKE) certs/clean

.PHONY: test/all/gotestfmt
## run tests for all Go codes and show table
test/all/gotestfmt: \
	ngt/install \
	hdf5/install \
	certs/gen \
	gotestfmt/install
	$(call go-test-gotestfmt,$(ROOTDIR)/...,$@,-showteststatus)
	$(MAKE) certs/clean

.PHONY: test/create-empty
## create empty test file if not exists
test/create-empty:
	@$(call green, "create empty test file if not exists...")
	@echo "$(GO_ALL_TEST_SOURCES)" | tr ' ' '\n' | xargs -P $(CORES) -I {} bash -c ' \
		if [ ! -f "{}" ]; then \
			echo "Creating empty test file {}"; \
			package="$$(basename $$(dirname {}))"; \
			if [ "$$(basename {})" = "main_test.go" ]; then \
				package="main"; \
			fi; \
			echo "package $$package" > "{}"; \
		fi \
	'

.PHONY: test/remove-empty
## remove empty test files
test/remove-empty:
	@$(call green, "remove empty test files...")
	@echo "$(GO_ALL_TEST_SOURCES)" | tr ' ' '\n' | xargs -P $(CORES) -I {} bash -c ' \
		if ! grep -q "func Test" "{}"; then \
			echo "Removing empty test file {}"; \
			rm "{}"; \
		fi \
	'

.PHONY: test/pkg
## run tests for pkg
test/pkg: \
	ngt/install \
	hdf5/install \
	certs/gen
	$(call go-test,$(ROOTDIR)/pkg/...)
	$(MAKE) certs/clean

.PHONY: test/internal
## run tests for internal
test/internal: \
	ngt/install \
	hdf5/install \
	certs/gen
	$(call go-test,$(ROOTDIR)/internal/...)
	$(MAKE) certs/clean

.PHONY: test/cmd
## run tests for cmd
test/cmd: \
	ngt/install \
	hdf5/install \
	certs/gen
	$(call go-test,$(ROOTDIR)/cmd/...)
	$(MAKE) certs/clean

.PHONY: test/rust
## run tests for rust
test/rust: \
	test/rust/qbg \
	test/rust/agent

.PHONY: test/rust/qbg
## run tests for qbg
test/rust/qbg:
	$(CC_ENV_VARS) \
	cargo test --manifest-path rust/Cargo.toml --package qbg --lib -- tests::test_ffi_qbg --exact --show-output
	$(CC_ENV_VARS) \
	cargo test --manifest-path rust/Cargo.toml --package qbg --lib -- tests::test_ffi_qbg_prebuilt --exact --show-output
	rm -rf rust/libs/algorithms/qbg/index/
	$(CC_ENV_VARS) \
	cargo test --manifest-path rust/Cargo.toml --package qbg --lib -- tests::test_property --exact --show-output
	$(CC_ENV_VARS) \
	cargo test --manifest-path rust/Cargo.toml --package qbg --lib -- tests::test_index --exact --show-output
	rm -rf rust/libs/algorithms/qbg/index/

.PHONY: test/rust/agent
## run tests for agent
test/rust/agent:
	$(CC_ENV_VARS) \
	cargo test --manifest-path rust/Cargo.toml --package agent -- handler::common::tests --show-output

.PHONY: test/hack
## run tests for hack
test/hack: \
	ngt/install \
	hdf5/install \
	certs/gen
	GOPRIVATE=$(GOPRIVATE) \
	go mod vendor -o $(ROOTDIR)/vendor
	$(call go-test,-mod=vendor $(ROOTDIR)/hack/gorules/... $(ROOTDIR)/hack/helm/... $(ROOTDIR)/hack/license/... $(ROOTDIR)/hack/tools/...)
	rm -rf $(ROOTDIR)/vendor
	$(MAKE) certs/clean

.PHONY: test/all
## run tests for all Go codes
test/all: \
	ngt/install \
	hdf5/install \
	certs/gen
	$(call go-test,$(ROOTDIR)/...)
	$(MAKE) certs/clean

.PHONY: coverage
## calculate coverages
coverage: \
	ngt/install \
	hdf5/install \
	certs/gen
	$(GO_TEST_ENV) \
	go test $(GO_TEST_FLAGS) -v -covermode=atomic -coverprofile=coverage.out $(ROOTDIR)/...
	$(GO_ENV_VARS) \
	go tool cover -html=coverage.out -o coverage.html
	$(MAKE) certs/clean

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
	@cat $(ROOTDIR)/.gitfiles | grep -E '^(\./)?internal/k8s/.*\_test.go$$' | xargs -I {} -P$(CORES) bash -c ' \
	echo "Replacing internal/k8s Test File {}" && \
	sed -i -E "s%k8s.io/apimachinery/pkg/api/errors%$(GOPKG)/internal/errors%g" {} && \
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
	@echo "$(GO_ALL_TEST_SOURCES)" | tr ' ' '\n' | xargs -I {} -P$(CORES) bash -c ' \
		[ -f "{}" ] || exit 0 ; \
		sed -i -E "/\/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/,\$$d" "{}"; \
		echo "// $(TEST_NOT_IMPL_PLACEHOLDER)" >> {};'

.PHONY: test/comment-unimplemented
## comment out unimplemented tests (from placeholder to EOF)
test/comment-unimplemented:
	@$(call green, "comment out unimplemented test (from placeholder to EOF)...")
	@echo "$(GO_ALL_TEST_SOURCES)" | tr ' ' '\n' | xargs -I {} -P$(CORES) bash -c ' \
		[ -f "{}" ] || exit 0 ; \
		sed -i -E -e " \
		/\/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/,\$$ { \
		s/^/\/\/ /; \
		s/^\/\/ \/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/\/\/ $(TEST_NOT_IMPL_PLACEHOLDER)/; \
		}" "{}"'