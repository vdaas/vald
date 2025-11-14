#
# Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
golangci-lint/install: \
	$(BINDIR)/golangci-lint

$(BINDIR)/golangci-lint:
	curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
	| sh -s -- -b $(BINDIR) $(GOLANGCILINT_VERSION)

.PHONY: goimports/install
## install goimports
goimports/install: \
	$(GOBIN)/goimports

$(GOBIN)/goimports:
	$(call go-tool-install)

.PHONY: strictgoimports/install
## install strictgoimports
strictgoimports/install: \
	$(GOBIN)/strictgoimports

$(GOBIN)/strictgoimports:
	$(call go-tool-install)

.PHONY: gofumpt/install
## install gofumpt
gofumpt/install: \
	$(GOBIN)/gofumpt

$(GOBIN)/gofumpt:
	$(call go-tool-install)

.PHONY: golines/install
## install golines
golines/install: \
	$(GOBIN)/golines

$(GOBIN)/golines:
	$(call go-tool-install)

.PHONY: crlfmt/install
## install crlfmt
crlfmt/install: \
	$(GOBIN)/crlfmt

$(GOBIN)/crlfmt:
	$(call go-tool-install)

.PHONY: actionlint/install
## install actionlint
actionlint/install: \
	$(GOBIN)/actionlint

$(GOBIN)/actionlint:
	$(call go-tool-install)

.PHONY: ghalint/install
## install ghalint
ghalint/install: \
	$(GOBIN)/ghalint

$(GOBIN)/ghalint:
	$(call go-tool-install)

.PHONY: pinact/install
## install pinact
pinact/install: \
	$(GOBIN)/pinact

$(GOBIN)/pinact:
	$(call go-tool-install)

.PHONY: ghatm/install
## install ghatm
ghatm/install: \
	$(GOBIN)/ghatm

$(GOBIN)/ghatm:
	$(call go-tool-install)

.PHONY: buf/install
## install buf
buf/install: \
	$(GOBIN)/buf

$(GOBIN)/buf:
	$(call go-tool-install)

.PHONY: k9s/install
## install k9s
k9s/install: \
	$(GOBIN)/k9s

$(GOBIN)/k9s:
	$(call go-tool-install)

.PHONY: stern/install
## install stern
stern/install: \
	$(GOBIN)/stern

$(GOBIN)/stern:
	$(call go-tool-install)

.PHONY: yamlfmt/install
## install yamlfmt
yamlfmt/install: \
	$(GOBIN)/yamlfmt

$(GOBIN)/yamlfmt:
	$(call go-tool-install)

.PHONY: gomodifytags/install
## install gomodifytags
gomodifytags/install: \
	$(GOBIN)/gomodifytags

$(GOBIN)/gomodifytags:
	$(call go-tool-install)

.PHONY: impl/install
## install impl
impl/install: \
	$(GOBIN)/impl

$(GOBIN)/impl:
	$(call go-tool-install)

.PHONY: delve/install
## install delve
delve/install: \
	$(GOBIN)/dlv

$(GOBIN)/dlv:
	$(call go-tool-install)

.PHONY: staticcheck/install
## install staticcheck
staticcheck/install: \
	$(GOBIN)/staticcheck

$(GOBIN)/staticcheck:
	$(call go-tool-install)

.PHONY: tparse/install
## install tparse
tparse/install: \
	$(GOBIN)/tparse

$(GOBIN)/tparse:
	$(call go-tool-install)

.PHONY: gotestfmt/install
## install gotestfmt
gotestfmt/install: \
	$(GOBIN)/gotestfmt

$(GOBIN)/gotestfmt:
	$(call go-tool-install)

.PHONY: gotests/install
## install gotests
gotests/install: \
	$(GOBIN)/gotests

$(GOBIN)/gotests:
	$(call go-tool-install)

.PHONY: protoc-gen-doc/install
## install protoc-gen-doc
protoc-gen-doc/install: \
	$(GOBIN)/protoc-gen-doc

$(GOBIN)/protoc-gen-doc:
	$(call go-tool-install)

.PHONY: go/tools/install
## install go tools
go/tools/install:
	$(call go-tool-install)

.PHONY: gopls/install
## install gopls
gopls/install: \
	$(GOBIN)/gopls

$(GOBIN)/gopls:
	GO111MODULE=on go install -mod=readonly golang.org/x/tools/gopls@latest

.PHONY: prettier/install
## Install prettier via Bun (global)
prettier/install: $(BUN_GLOBAL_BIN)/prettier
$(BUN_GLOBAL_BIN)/prettier: bun/install
	command -v prettier >/dev/null 2>&1 || bun add --global prettier

.PHONY: reviewdog/install
## install reviewdog
reviewdog/install: $(BINDIR)/reviewdog

$(BINDIR)/reviewdog:
	curl -fsSL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
	| sh -s -- -b $(BINDIR) $(REVIEWDOG_VERSION)

.PHONY: mbake/install
## install mbake
mbake/install: $(BINDIR)/mbake

$(BINDIR)/mbake:
	pip install mbake --break-system-packages --prefix /usr

.PHONY: kubectl/install
## install kubectl
kubectl/install: $(BINDIR)/kubectl

$(BINDIR)/kubectl:
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	curl -fsSL "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/$(OS)/$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]'))/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl

.PHONY: textlint/install
## Install textlint & rules via Bun (global)
textlint/install: $(BUN_GLOBAL_BIN)/textlint
$(BUN_GLOBAL_BIN)/textlint: bun/install
	bun add --global \
	textlint \
	textlint-rule-en-spell \
	textlint-rule-prh \
	textlint-rule-write-good

.PHONY: textlint/ci/install
## Install textlint & rules for CI via Bun (local devDependencies)
textlint/ci/install: bun/install
	[ -f package.json ] || (bun init -y >/dev/null 2>&1 || echo '{}' > package.json)
	bun add --dev \
	textlint \
	textlint-rule-en-spell \
	textlint-rule-prh \
	textlint-rule-write-good

.PHONY: cspell/install
## Install cspell & dictionaries via Bun (global)
cspell/install: $(BUN_GLOBAL_BIN)/cspell
$(BUN_GLOBAL_BIN)/cspell: bun/install
	bun add --global \
	cspell@latest \
	@cspell/dict-cpp \
	@cspell/dict-docker \
	@cspell/dict-en_us \
	@cspell/dict-fullstack \
	@cspell/dict-git \
	@cspell/dict-golang \
	@cspell/dict-k8s \
	@cspell/dict-makefile \
	@cspell/dict-markdown \
	@cspell/dict-npm \
	@cspell/dict-public-licenses \
	@cspell/dict-rust \
	@cspell/dict-shell
	cspell link add @cspell/dict-cpp
	cspell link add @cspell/dict-docker
	cspell link add @cspell/dict-en_us
	cspell link add @cspell/dict-fullstack
	cspell link add @cspell/dict-git
	cspell link add @cspell/dict-golang
	cspell link add @cspell/dict-k8s
	cspell link add @cspell/dict-makefile
	cspell link add @cspell/dict-markdown
	cspell link add @cspell/dict-npm
	cspell link add @cspell/dict-public-licenses
	cspell link add @cspell/dict-rust
	cspell link add @cspell/dict-shell

.PHONY: bun/install
## Install Bun runtime into $(BUN_INSTALL) if not already installed
bun/install: $(BINDIR)/bun

$(BINDIR)/bun:
	curl -fsSL https://bun.sh/install | BUN_INSTALL=$(BUN_INSTALL) bash

.PHONY: go/install
## install go
go/install: $(GOROOT)/bin/go

$(GOROOT)/bin/go:
	TAR_NAME=go$(GO_VERSION).$(OS)-$(subst x86_64,amd64,$(subst aarch64,arm64,$(ARCH))).tar.gz \
	&& curl -fsSL "https://go.dev/dl/$${TAR_NAME}" -o "$(TEMP_DIR)/$${TAR_NAME}" \
	&& mkdir -p $(TEMP_DIR)/go \
	&& tar -xzvf "$(TEMP_DIR)/$${TAR_NAME}" -C $(TEMP_DIR)/go --strip-components 1 \
	&& rm -rf "$(TEMP_DIR)/$${TAR_NAME}" \
	&& mv $(TEMP_DIR)/go $(GOROOT) \
	&& $(GOROOT)/bin/go version

.PHONY: rust/install
## install rust
rust/install: $(CARGO_HOME)/bin/cargo

$(CARGO_HOME)/bin/cargo:
	curl --proto '=https' --tlsv1.2 -fsSL https://sh.rustup.rs | CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} sh -s -- --default-toolchain $(RUST_VERSION) -y
	rustup toolchain install $(RUST_VERSION)
	rustup default $(RUST_VERSION)
	source "${CARGO_HOME}/env"

.PHONY: zlib/install
## install zlib
zlib/install: $(LIB_PATH)/libz.a

$(LIB_PATH)/libz.a: $(LIB_PATH)
	curl -fsSL https://github.com/madler/zlib/releases/download/v$(ZLIB_VERSION)/zlib-$(ZLIB_VERSION).tar.gz -o $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz \
	&& mkdir -p $(TEMP_DIR)/zlib \
	&& tar -xzvf $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz -C $(TEMP_DIR)/zlib --strip-components 1 \
	&& cd $(TEMP_DIR)/zlib \
	&& mkdir -p build \
	&& cd build \
	&& cmake	-DCMAKE_BUILD_TYPE=Release \
	-DCMAKE_POLICY_VERSION_MINIMUM=$(CMAKE_VERSION) \
	-DBUILD_SHARED_LIBS=OFF \
	-DBUILD_STATIC_EXECS=ON \
	-DBUILD_TESTING=OFF \
	-DZLIB_BUILD_SHARED=OFF \
	-DZLIB_BUILD_STATIC=ON \
	-DZLIB_COMPAT=ON \
	-DZLIB_USE_STATIC_LIBS=ON \
	-DCMAKE_CXX_FLAGS="$(CXXFLAGS)" \
	-DCMAKE_C_FLAGS="$(CFLAGS)" \
	-DCMAKE_INSTALL_LIBDIR=$(LIB_PATH) \
	-DCMAKE_INSTALL_PREFIX=$(USR_LOCAL) \
	-B $(TEMP_DIR)/zlib/build $(TEMP_DIR)/zlib \
	&& make -j$(CORES) \
	&& make install \
	&& cd $(ROOTDIR) \
	&& rm -rf $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz $(TEMP_DIR)/zlib $(LIB_PATH)/libz.s*

.PHONY: hdf5/install
## install hdf5
hdf5/install: $(LIB_PATH)/libhdf5.a

$(LIB_PATH)/libhdf5.a: $(LIB_PATH) \
	zlib/install
	mkdir -p $(TEMP_DIR)/hdf5 \
	&& curl -fsSL https://github.com/HDFGroup/hdf5/archive/refs/tags/$(HDF5_VERSION).tar.gz -o $(TEMP_DIR)/hdf5.tar.gz \
	&& tar -xzvf $(TEMP_DIR)/hdf5.tar.gz -C $(TEMP_DIR)/hdf5 --strip-components 1 \
	&& mkdir -p $(TEMP_DIR)/hdf5/build \
	&& cd $(TEMP_DIR)/hdf5/build \
	&& cmake -DCMAKE_BUILD_TYPE=Release \
	-DCMAKE_POLICY_VERSION_MINIMUM=$(CMAKE_VERSION) \
	-DBUILD_SHARED_LIBS=OFF \
	-DBUILD_STATIC_EXECS=ON \
	-DBUILD_TESTING=OFF \
	-DHDF5_BUILD_CPP_LIB=OFF \
	-DHDF5_BUILD_HL_LIB=ON \
	-DHDF5_BUILD_STATIC_EXECS=ON \
	-DHDF5_BUILD_TOOLS=OFF \
	-DHDF5_ENABLE_Z_LIB_SUPPORT=ON \
	-DH5_ZLIB_INCLUDE_DIR=$(USR_LOCAL)/include \
	-DH5_ZLIB_LIBRARY=$(LIB_PATH)/libz.a \
	-DCMAKE_CXX_FLAGS="$(CXXFLAGS)" \
	-DCMAKE_C_FLAGS="$(CFLAGS)" \
	-DCMAKE_INSTALL_LIBDIR=$(LIB_PATH) \
	-DCMAKE_INSTALL_PREFIX=$(USR_LOCAL) \
	-B $(TEMP_DIR)/hdf5/build $(TEMP_DIR)/hdf5 \
	&& make -j$(CORES) \
	&& make install \
	&& cd $(ROOTDIR) \
	&& rm -rf $(TEMP_DIR)/hdf5.tar.gz $(TEMP_DIR)/hdf5 \
	&& ldconfig

.PHONY: yq/install
## install yq
yq/install: $(BINDIR)/yq

$(BINDIR)/yq:
	mkdir -p $(BINDIR)
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	cd $(TEMP_DIR) \
	&& curl -fsSL https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(OS)_$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]')) -o $(BINDIR)/yq \
	&& chmod a+x $(BINDIR)/yq

.PHONY: docker-cli/install
## install docker-cli
docker-cli/install: $(BINDIR)/docker

$(BINDIR)/docker: $(BINDIR)
	curl -fsSL https://download.docker.com/linux/static/stable/$(shell uname -m)/docker-$(shell echo $(DOCKER_VERSION) | cut -c2-).tgz -o $(TEMP_DIR)/docker.tgz \
	&& tar -xzvf $(TEMP_DIR)/docker.tgz -C $(TEMP_DIR) \
	&& mv $(TEMP_DIR)/docker/docker $(BINDIR) \
	&& rm -rf $(TEMP_DIR)/docker{.tgz,}

.PHONY: replace/busybox
## replace busybox version
replace/busybox:
	find . -type f \( -name "*.yaml" -o -name "*.md" \) -exec sed -i -E 's/busybox:([0-9]+\.[0-9]+\.[0-9]+|latest)/busybox:$(BUSYBOX_VERSION)/g' {} +