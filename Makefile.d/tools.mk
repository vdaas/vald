#
# Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

.PHONY: golangci-lint/install
## install golangci-lint
golangci-lint/install: $(BINDIR)/golangci-lint

$(BINDIR)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(BINDIR) $(GOLANGCILINT_VERSION)

.PHONY: goimports/install
goimports/install: $(GOBIN)/goimports

$(GOBIN)/goimports:
	$(call go-install, golang.org/x/tools/cmd/goimports)

.PHONY: strictgoimports/install
strictgoimports/install: $(GOBIN)/strictgoimports

$(GOBIN)/strictgoimports:
	$(call go-install, github.com/momotaro98/strictgoimports/cmd/strictgoimports)

.PHONY: gofumpt/install
gofumpt/install: $(GOBIN)/gofumpt

$(GOBIN)/gofumpt:
	$(call go-install, mvdan.cc/gofumpt)

.PHONY: golines/install
golines/install: $(GOBIN)/golines

$(GOBIN)/golines:
	$(call go-install, github.com/segmentio/golines)

.PHONY: prettier/install
prettier/install: $(BINDIR)/prettier
$(BINDIR)/prettier:
	type prettier || npm install -g prettier

.PHONY: reviewdog/install
## install reviewdog
reviewdog/install: $(BINDIR)/reviewdog

$(BINDIR)/reviewdog:
	curl -sSfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
		| sh -s -- -b $(BINDIR) $(REVIEWDOG_VERSION)

.PHONY: kubectl/install
kubectl/install: $(BINDIR)/kubectl

$(BINDIR)/kubectl:
	curl -L "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/$(OS)/$(subst x86_64,amd64,$(shell echo $(ARCH) | tr '[:upper:]' '[:lower:]'))/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl

.PHONY: textlint/install
textlint/install: $(NPM_GLOBAL_PREFIX)/bin/textlint

$(NPM_GLOBAL_PREFIX)/bin/textlint:
	npm install -g textlint textlint-rule-en-spell textlint-rule-prh textlint-rule-write-good

.PHONY: textlint/ci/install
textlint/ci/install:
	npm init -y
	npm install --save-dev textlint textlint-rule-en-spell textlint-rule-prh textlint-rule-write-good

.PHONY: cspell/install
cspell/install: $(NPM_GLOBAL_PREFIX)/bin/cspell

$(NPM_GLOBAL_PREFIX)/bin/cspell:
	npm install -g git+https://github.com/streetsidesoftware/cspell-cli

.PHONY: buf/install
buf/install: $(BINDIR)/buf

$(BINDIR)/buf:
	curl -sSL \
	"https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m)" \
	-o "${BINDIR}/buf" && \
	chmod +x "${BINDIR}/buf"

.PHONY: k9s/install
k9s/install: $(GOPATH)/bin/k9s

$(GOPATH)/bin/k9s:
	$(call go-install, github.com/derailed/k9s)

.PHONY: stern/install
stern/install: $(GOPATH)/bin/stern

$(GOPATH)/bin/stern:
	$(call go-install, github.com/stern/stern)

.PHONY: yamlfmt/install
yamlfmt/install: $(GOPATH)/bin/yamlfmt

$(GOPATH)/bin/yamlfmt:
	$(call go-install, github.com/google/yamlfmt/cmd/yamlfmt)

.PHONY: gopls/install
gopls/install: $(GOPATH)/bin/gopls

$(GOPATH)/bin/gopls:
	$(call go-install, golang.org/x/tools/gopls)

.PHONY: gomodifytags/install
gomodifytags/install: $(GOPATH)/bin/gomodifytags

$(GOPATH)/bin/gomodifytags:
	$(call go-install, github.com/fatih/gomodifytags)

.PHONY: impl/install
impl/install: $(GOPATH)/bin/impl

$(GOPATH)/bin/impl:
	$(call go-install, github.com/josharian/impl)

.PHONY: goplay/install
goplay/install: $(GOPATH)/bin/goplay

$(GOPATH)/bin/goplay:
	$(call go-install, github.com/haya14busa/goplay/cmd/goplay)

.PHONY: delve/install
delve/install: $(GOPATH)/bin/dlv

$(GOPATH)/bin/dlv:
	$(call go-install, github.com/go-delve/delve/cmd/dlv)

.PHONY: staticcheck/install
staticcheck/install: $(GOPATH)/bin/staticcheck

$(GOPATH)/bin/staticcheck:
	$(call go-install, honnef.co/go/tools/cmd/staticcheck)

.PHONY: go/install
go/install: $(GOROOT)/bin/go

$(GOROOT)/bin/go:
	TAR_NAME=go$(GO_VERSION).$(OS)-$(subst x86_64,amd64,$(subst aarch64,arm64,$(ARCH))).tar.gz \
	&& curl -fsSLO "https://go.dev/dl/$${TAR_NAME}" \
	&& tar zxf "$${TAR_NAME}" \
	&& rm -rf "$${TAR_NAME}" \
	&& mv go $(GOROOT) \
	&& $(GOROOT)/bin/go version \
	&& mkdir -p "$(GOPATH)/src" "$(GOPATH)/bin" "$(GOPATH)/pkg"

.PHONY: rust/install
rust/install: $(CARGO_HOME)/bin/cargo

$(CARGO_HOME)/bin/cargo:
	curl --proto '=https' --tlsv1.2 -fsSL https://sh.rustup.rs | CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} sh -s -- --default-toolchain nightly -y
	source "${CARGO_HOME}/env" \
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} ${CARGO_HOME}/bin/rustup install stable \
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} ${CARGO_HOME}/bin/rustup install beta \
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} ${CARGO_HOME}/bin/rustup install nightly \
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} ${CARGO_HOME}/bin/rustup toolchain install nightly \
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} ${CARGO_HOME}/bin/rustup default nightly \
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} ${CARGO_HOME}/bin/rustup update
