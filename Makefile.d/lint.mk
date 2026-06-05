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
	docs/lint \
	files/lint \
	vet \
	go/lint \
	workflow/lint

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
	textlint $(ROOTDIR)/docs/**/*.md $(TEXTLINT_EXTRA_OPTIONS)

.PHONY: files/textlint
## run textlint for document
files/textlint: \
	files \
	textlint/install
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then textlint "$(ROOTDIR)/.gitfiles" $(TEXTLINT_EXTRA_OPTIONS); fi

.PHONY: docs/cspell
## run cspell for document
docs/cspell: \
	cspell/install
	cspell $(ROOTDIR)/docs/**/*.md --show-suggestions $(CSPELL_EXTRA_OPTIONS)

.PHONY: files/cspell
## run cspell for document
files/cspell: \
	files \
	cspell/install
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then cspell "$(ROOTDIR)/.gitfiles" --show-suggestions $(CSPELL_EXTRA_OPTIONS); fi