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

.PHONY: format
## format go codes
format:
	- @$(MAKE) format/proto
	- @$(MAKE) format/json
	- @$(MAKE) format/md
	- @$(MAKE) remove/empty/file
	- @$(MAKE) replace/busybox
	- @$(MAKE) license
	- @$(MAKE) dockerfile
	- @$(MAKE) format/go
	- @$(MAKE) format/go/test
	- @$(MAKE) format/yaml
	- @$(MAKE) format/make
	- @$(MAKE) format/rust

.PHONY: format/diff
## format diff
format/diff:
	@$(MAKE) format/go/diff
	@$(MAKE) format/yaml/diff

.PHONY: remove/empty/file
## removes empty file such as just includes \r \n space tab
remove/empty/file: \
	files
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -vE '^\s*#' "$(ROOTDIR)/.gitfiles" | grep -v gitkeep \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" -n1 sh -c ' \
		if [ -f "{}" ] && [ -z "$$(tr -d '\''[:space:]'\'' < "{}")" ]; then rm "{}"; fi'; \
		fi

.PHONY: format/go
## run golines, gofumpt, goimports for all go files
format/go: \
	crlfmt/install \
	strictgoimports/install \
	golangci-lint/install \
	files
	@echo "Formatting Go files..."
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -e "\.go$$" "$(ROOTDIR)/.gitfiles" | grep -v "_test\.go$$" \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		echo "Formatting Go file {}" && \
		$(GOBIN)/strictgoimports -w {} && \
		$(GOBIN)/crlfmt -w -diff=false {} && \
		$(BINDIR)/golangci-lint fmt --config $(ROOTDIR)/.golangci.json {}'; \
	fi
	go fix $(ROOTDIR)/...
	@echo "Go formatting complete."

.PHONY: format/go/test
## run golines, gofumpt, goimports for go test files
format/go/test: \
	crlfmt/install \
	strictgoimports/install \
	golangci-lint/install \
	files
	@echo "Formatting Go Test files..."
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -e "_test\.go$$" "$(ROOTDIR)/.gitfiles" \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		echo "Formatting Go Test file {}" && \
		$(GOBIN)/strictgoimports -w {} && \
		$(GOBIN)/crlfmt -w -diff=false {} && \
		$(BINDIR)/golangci-lint fmt --config $(ROOTDIR)/.golangci.json {}'; \
	fi
	@echo "Go test file formatting complete."

.PHONY: format/go/diff
## run golines, gofumpt, goimports for go diff files
format/go/diff: \
	crlfmt/install \
	strictgoimports/install \
	golangci-lint/install \
	files
	@echo "Formatting Go Diff files..."
	@git diff --name-only --diff-filter=ACM HEAD | grep -e ".go$$" | xargs -I {} -P$(CORES) bash -c ' \
	echo "Formatting Go file {}" && \
	$(GOBIN)/strictgoimports -w {} && \
	$(GOBIN)/crlfmt -w -diff=false {} && \
	$(BINDIR)/golangci-lint fmt --config $(ROOTDIR)/.golangci.json {}'
	@echo "Go file formatting complete."

.PHONY: format/rust
## format rust codes
format/rust: \
	rustfmt/install \
	files
	@echo "Formatting Rust files..."
	@cd $(ROOTDIR)/rust && $(CARGO_HOME)/bin/cargo fmt
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -e "\.rs$$" "$(ROOTDIR)/.gitfiles" \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		echo "Formatting Rust file {}" && \
		$(CARGO_HOME)/bin/rustfmt --edition 2024 --style-edition 2024 {}'; \
	fi
	@echo "Rust formatting complete."

.PHONY: format/yaml
## format yaml file
format/yaml:
	- @$(MAKE) prettier/install
	- @$(MAKE) yamlfmt/install
	- @$(MAKE) clean/yaml
	- @$(MAKE) files
	@echo "Formatting YAML files..."
	- @if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -E '\.ya?ml$$' "$(ROOTDIR)/.gitfiles" | grep -Ev '(templates|s3)' \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		if [ ! -L "{}" ]; then \
			echo "Formatting YAML file {}" && \
			yamlfmt {} && \
			bunx prettier --write {}; \
		fi'; \
	fi
	@echo "YAML file formatting complete."

.PHONY: clean/yaml
## cleanup empty yaml file
clean/yaml:
	- @$(MAKE) files
	- @if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -E '\.ya?ml$$' "$(ROOTDIR)/.gitfiles" | grep -Ev '(templates|s3)' \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" -n1 sh -c ' \
		if [ -f "{}" ] && [ ! -L "{}" ] && [ -z "$$(tr -d '\''[:space:]'\'' < "{}")" ]; then rm "{}"; fi'; \
		fi

.PHONY: format/yaml/diff
format/yaml/diff:
	- @$(MAKE) prettier/install
	- @$(MAKE) yamlfmt/install
	- @$(MAKE) clean/yaml
	- @$(MAKE) files
	@echo "Formatting YAML files..."
	- @git diff --name-only --diff-filter=ACM HEAD | grep -E '\.ya?ml$$' | grep -Ev '(templates|s3)' | xargs -I {} -P$(CORES) bash -c ' \
		if [ ! -L "{}" ]; then \
			echo "Formatting YAML file {}" && \
			yamlfmt {} && \
			bunx prettier --write {}; \
		fi'
	@echo "YAML file formatting complete."

.PHONY: format/make
format/make: \
	mbake/install \
	files
	@echo "Formatting Makefile and Makefile.d/*.mk files..."
	- @if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -E '(Makefile|\.mk$$)' "$(ROOTDIR)/.gitfiles" \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		if [ ! -L "{}" ]; then \
			echo "Formatting Make file {}" && \
			mbake format {}; \
		fi'; \
	fi
	@echo "Makefile formatting complete."

.PHONY: format/md
format/md: \
	prettier/install
	@echo "Formatting Markdown files..."
	- @if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -E '\.md$$' "$(ROOTDIR)/.gitfiles" \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		if [ ! -L "{}" ]; then \
			echo "Formatting Markdown file {}" && \
			bunx prettier --write {}; \
		fi'; \
	fi
	@echo "Markdown file formatting complete."

.PHONY: format/json
format/json: \
	prettier/install
	@echo "Formatting JSON files..."
	- @if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -E '\.json$$' "$(ROOTDIR)/.gitfiles" \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c ' \
		if [ ! -L "{}" ]; then \
			echo "Formatting JSON file {}" && \
			bunx prettier --write {}; \
		fi'; \
	fi
	@echo "JSON file formatting complete."
.PHONY: format/proto
format/proto: \
	buf/install
	buf format -w