#
# Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
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

.PHONY: tparse/install
## install tparse
tparse/install: \
        $(GOPATH)bin/tparse

$(GOPATH)bin/tparse:
	$(call go-install, github.com/mfridman/tparse)

.PHONY: gotestfmt/install
## install gotestfmt
gotestfmt/install: \
        $(GOPATH)bin/gotestfmt

$(GOPATH)bin/gotestfmt:
	$(call go-install, github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt)

.PHONY: gotests/install
## install gotests
gotests/install: \
        $(GOPATH)bin/gotests

$(GOPATH)bin/gotests:
	$(call go-install, github.com/cweill/gotests/gotests)

## run tests for cmd, internal, pkg
test:
	go test -shuffle=on -race -mod=readonly -cover -timeout=30m ./cmd/... ./internal/... ./pkg/...

.PHONY: test/tparse
## run tests for cmd, internal, pkg and show table
test/tparse: \
        tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./cmd/... ./internal/... ./pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests

.PHONY: test/cmd/tparse
## run tests for cmd and show table
test/cmd/tparse: \
        tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./cmd/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests

.PHONY: test/internal/tparse
## run tests for internal and show table
test/internal/tparse: \
        tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./internal/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests

.PHONY: test/pkg/tparse
## run tests for pkg and who table
test/pkg/tparse: \
        tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests

.PHONY: test/hack/tparse
## run tests for hack and show table
test/hack/tparse: \
        tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover \
		./hack/gorules/... \
		./hack/helm/... \
		./hack/license/... \
		./hack/tools/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests

.PHONY: test/all/tparse
## run tests for all Go codes and show table
test/all/tparse: \
        tparse/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| tparse -pass -notests

.PHONY: test/gotestfmt
## run tests for cmd, internal, pkg and show table
test/gotestfmt: \
        gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./cmd/... ./internal/... ./pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus

.PHONY: test/cmd/gotestfmt
## run tests for cmd and show table
test/cmd/gotestfmt: \
        gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./cmd/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus

.PHONY: test/internal/gotestfmt
## run tests for internal and show table
test/internal/gotestfmt: \
        gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./internal/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus

.PHONY: test/pkg/gotestfmt
## run tests for pkg and who table
test/pkg/gotestfmt: \
        gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./pkg/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus

.PHONY: test/hack/gotestfmt
## run tests for hack and show table
test/hack/gotestfmt: \
        gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover \
		./hack/gorules/... \
		./hack/helm/... \
		./hack/license/... \
		./hack/tools/... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus

.PHONY: test/all/gotestfmt
## run tests for all Go codes and show table
test/all/gotestfmt: \
        gotestfmt/install
	set -euo pipefail
	rm -rf "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json"
	go test -shuffle=on -race -mod=readonly -json -cover -timeout=30m ./... \
	| tee "$(TEST_RESULT_DIR)/`echo $@ | sed -e 's%/%-%g'`-result.json" \
	| gotestfmt -showteststatus

.PHONY: test/pkg
## run tests for pkg
test/pkg:
	go test -shuffle=on -race -mod=readonly -cover ./pkg/...

.PHONY: test/internal
## run tests for internal
test/internal:
	go test -shuffle=on -race -mod=readonly -cover ./internal/...

.PHONY: test/cmd
## run tests for cmd
test/cmd:
	go test -shuffle=on -race -mod=readonly -cover ./cmd/...

.PHONY: test/hack
## run tests for hack
test/hack:
	go test -shuffle=on -race -mod=readonly -cover \
		./hack/gorules... \
		./hack/helm/... \
		./hack/license/...\
		./hack/tools/...

.PHONY: test/all
## run tests for all Go codes
test/all:
	go test -shuffle=on -race -mod=readonly -cover ./...

.PHONY: coverage
## calculate coverages
coverage:
	go test -shuffle=on -race -mod=readonly -v -race -covermode=atomic -timeout=30m -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: gotests/gen
## generate missing go test files
gotests/gen: \
	$(GO_TEST_SOURCES) \
	$(GO_OPTION_TEST_SOURCES) \
	gotests/patch

.PHONY: gotests/patch
## apply patches to generated go test files
gotests/patch: \
	$(GO_TEST_SOURCES) \
	$(GO_OPTION_TEST_SOURCES)
	@$(call green, "apply patches to go test files...")
	find $(ROOTDIR)/internal/k8s/* -name '*_test.go' | xargs sed -i -E "s%k8s.io/apimachinery/pkg/api/errors%github.com/vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%cockroachdb/errors%vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%golang.org/x/sync/errgroup%github.com/vdaas/vald/internal/errgroup%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%pkg/errors%vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%go-errors/errors%vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%go.uber.org/goleak%github.com/vdaas/vald/internal/test/goleak%g"
	find $(ROOTDIR)/internal/errors -name '*_test.go' | xargs sed -i -E "s%\"github.com/vdaas/vald/internal/errors\"%%g"
	find $(ROOTDIR)/internal/errors -name '*_test.go' | xargs sed -i -E "s/errors\.//g"
	find $(ROOTDIR)/internal/test/goleak -name '*_test.go' | xargs sed -i -E "s%\"github.com/vdaas/vald/internal/test/goleak\"%%g"
	find $(ROOTDIR)/internal/test/goleak -name '*_test.go' | xargs sed -i -E "s/goleak\.//g"

$(GO_TEST_SOURCES): \
	./assets/test/templates/common \
	$(GO_SOURCES)
	@$(call green, $(patsubst %,"generating go test file: %",$@))
	gotests -w -template_dir ./assets/test/templates/common -all $(patsubst %_test.go,%.go,$@)

$(GO_OPTION_TEST_SOURCES): \
	./assets/test/templates/option \
	$(GO_OPTION_SOURCES)
	@$(call green, $(patsubst %,"generating go test file: %",$@))
	gotests -w -template_dir ./assets/test/templates/option -all $(patsubst %_test.go,%.go,$@)
